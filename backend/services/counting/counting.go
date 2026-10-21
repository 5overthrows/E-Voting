package counting

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"evoting/services/auth"
	"evoting/shared/crypto"
	"evoting/shared/db"
	"evoting/shared/merkle"
	"evoting/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// System-wide RSA keypair for signing Merkle roots (for prototype ease-of-use)
var (
	systemPrivateKeyPEM string
	systemPublicKeyPEM  string
)

func init() {
	pub, priv, err := crypto.GenerateRSAKeypair()
	if err != nil {
		log.Fatalf("Failed to initialize system keys: %v", err)
	}
	systemPublicKeyPEM = pub
	systemPrivateKeyPEM = priv
}

type CreateElectionRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Passphrase  string    `json:"passphrase" binding:"required"` // For encrypting the private key
}

type CandidateRequest struct {
	Name  string `json:"name" binding:"required"`
	Party string `json:"party" binding:"required"`
}

type TallyRequest struct {
	Passphrase string `json:"passphrase" binding:"required"`
}

// RegisterRoutes registers the counting & election service routes on the given router
func RegisterRoutes(r *gin.Engine) {
	// Public routes
	r.GET("/api/elections", listElections)
	r.GET("/api/elections/:id", getElection)
	r.GET("/api/elections/:id/results", getResults)

	// Election Officer / Counting Authority protected routes
	protected := r.Group("/api/elections")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/create", auth.AuthMiddleware("ElectionOfficer"), createElection)
		protected.POST("/:id/candidates", auth.AuthMiddleware("ElectionOfficer"), addCandidate)
		protected.POST("/:id/open", auth.AuthMiddleware("ElectionOfficer"), openElection)
		protected.POST("/:id/close", auth.AuthMiddleware("ElectionOfficer"), closeElection)
		
		// Counting Authority only
		protected.POST("/:id/tally", auth.AuthMiddleware("CountingAuthority"), tallyVotes)
		protected.GET("/:id/verify-merkle", auth.AuthMiddleware("CountingAuthority", "ElectionOfficer"), verifyMerkle)
	}
}

// StartCountingServer starts the counting & election service on the specified address
func StartCountingServer(addr string) {
	r := gin.Default()

	// Configure CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	RegisterRoutes(r)

	log.Printf("Counting & Election Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Counting server failed: %v", err)
	}
}

func listElections(c *gin.Context) {
	var elections []models.Election
	query := "SELECT ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, SignedMerkleRoot FROM Elections"
	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Election
		var signedRoot sql.NullString
		err := rows.Scan(&e.ElectionID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.PublicKey, &signedRoot)
		if err == nil {
			if signedRoot.Valid {
				e.SignedMerkleRoot = signedRoot.String
			}
			elections = append(elections, e)
		}
	}

	c.JSON(http.StatusOK, elections)
}

func getElection(c *gin.Context) {
	id := c.Param("id")
	var e models.Election
	var signedRoot sql.NullString

	query := "SELECT ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, SignedMerkleRoot FROM Elections WHERE ElectionID = $1"
	err := db.DB.QueryRow(query, id).Scan(&e.ElectionID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.PublicKey, &signedRoot)
	if err != nil {
		query = "SELECT ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, SignedMerkleRoot FROM Elections WHERE ElectionID = ?"
		err = db.DB.QueryRow(query, id).Scan(&e.ElectionID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.PublicKey, &signedRoot)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if signedRoot.Valid {
		e.SignedMerkleRoot = signedRoot.String
	}

	// Fetch candidates
	var candidates []models.Candidate
	cRows, err := db.DB.Query("SELECT CandidateID, CandidateName, Party FROM Candidates WHERE ElectionID = $1", id)
	if err != nil {
		cRows, err = db.DB.Query("SELECT CandidateID, CandidateName, Party FROM Candidates WHERE ElectionID = ?", id)
	}

	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var cand models.Candidate
			cand.ElectionID = id
			if err := cRows.Scan(&cand.CandidateID, &cand.CandidateName, &cand.Party); err == nil {
				candidates = append(candidates, cand)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"election":   e,
		"candidates": candidates,
	})
}

func createElection(c *gin.Context) {
	var req CreateElectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate RSA keypair for the election
	pubPEM, privPEM, err := crypto.GenerateRSAKeypair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate keys"})
		return
	}

	// Derive AES key from passphrase & salt to encrypt the private key
	salt, err := crypto.GenerateRandomSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}

	aesKey := crypto.DeriveKey(req.Passphrase, salt)
	encPrivKey, err := crypto.EncryptAES(aesKey, []byte(privPEM))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt private key"})
		return
	}

	electionID := uuid.New().String()
	status := "Created"

	query := "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err = db.DB.Exec(query, electionID, req.Title, req.Description, req.StartTime, req.EndTime, status, pubPEM, encPrivKey, salt)
	if err != nil {
		query = "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = db.DB.Exec(query, electionID, req.Title, req.Description, req.StartTime, req.EndTime, status, pubPEM, encPrivKey, salt)
	}

	if err != nil {
		log.Printf("Failed to insert election in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	db.LogEvent("Election created: "+req.Title, "ElectionOfficer", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"election_id": electionID,
		"message":     "Election created successfully",
	})
}

func addCandidate(c *gin.Context) {
	id := c.Param("id")
	var req CandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check election status
	var status string
	err := db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = $1", id).Scan(&status)
	if err != nil {
		err = db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = ?", id).Scan(&status)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if status != "Created" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add candidates to an election that is already active/closed"})
		return
	}

	candidateID := uuid.New().String()
	query := "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, candidateID, id, req.Name, req.Party)
	if err != nil {
		query = "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES (?, ?, ?, ?)"
		_, err = db.DB.Exec(query, candidateID, id, req.Name, req.Party)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	db.LogEvent("Candidate "+req.Name+" added to election: "+id, "ElectionOfficer", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"candidate_id": candidateID})
}

func openElection(c *gin.Context) {
	id := c.Param("id")

	query := "UPDATE Elections SET Status = 'Active' WHERE ElectionID = $1 AND Status = 'Created'"
	res, err := db.DB.Exec(query, id)
	if err != nil {
		query = "UPDATE Elections SET Status = 'Active' WHERE ElectionID = ? AND Status = 'Created'"
		res, err = db.DB.Exec(query, id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Election must be in 'Created' status to open"})
		return
	}

	db.LogEvent("Election opened: "+id, "ElectionOfficer", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"message": "Election status set to Active"})
}

func closeElection(c *gin.Context) {
	id := c.Param("id")

	// 1. Update status to Closed
	query := "UPDATE Elections SET Status = 'Closed' WHERE ElectionID = $1 AND Status = 'Active'"
	res, err := db.DB.Exec(query, id)
	if err != nil {
		query = "UPDATE Elections SET Status = 'Closed' WHERE ElectionID = ? AND Status = 'Active'"
		res, err = db.DB.Exec(query, id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Election must be in 'Active' status to close"})
		return
	}

	// 2. Compute final Merkle Root and digitally sign it (REQ-MRK-3)
	var leaves []string
	vRows, err := db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp ASC", id)
	if err != nil {
		vRows, err = db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp ASC", id)
	}

	if err == nil {
		defer vRows.Close()
		for vRows.Next() {
			var lh string
			if err := vRows.Scan(&lh); err == nil {
				leaves = append(leaves, lh)
			}
		}
	}

	finalRoot := merkle.ComputeMerkleRoot(leaves)

	// Digitally sign the root hash using the system-wide key
	signature, err := crypto.SignRSA(systemPrivateKeyPEM, []byte(finalRoot))
	if err != nil {
		log.Printf("Failed to sign Merkle Root: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign Merkle root"})
		return
	}

	// Save signed root hash in Elections table
	updateQuery := "UPDATE Elections SET SignedMerkleRoot = $1 WHERE ElectionID = $2"
	_, err = db.DB.Exec(updateQuery, signature, id)
	if err != nil {
		updateQuery = "UPDATE Elections SET SignedMerkleRoot = ? WHERE ElectionID = ?"
		_, err = db.DB.Exec(updateQuery, signature, id)
	}

	db.LogEvent("Election closed and Merkle Root signed: "+id, "ElectionOfficer", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"message":            "Election status set to Closed",
		"merkle_root":        finalRoot,
		"signed_merkle_root": signature,
	})
}

func verifyMerkle(c *gin.Context) {
	id := c.Param("id")

	var signedRoot sql.NullString
	var status string
	err := db.DB.QueryRow("SELECT Status, SignedMerkleRoot FROM Elections WHERE ElectionID = $1", id).Scan(&status, &signedRoot)
	if err != nil {
		err = db.DB.QueryRow("SELECT Status, SignedMerkleRoot FROM Elections WHERE ElectionID = ?", id).Scan(&status, &signedRoot)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Retrieve leaves
	var leaves []string
	vRows, err := db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp ASC", id)
	if err != nil {
		vRows, err = db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp ASC", id)
	}

	if err == nil {
		defer vRows.Close()
		for vRows.Next() {
			var lh string
			if err := vRows.Scan(&lh); err == nil {
				leaves = append(leaves, lh)
			}
		}
	}

	computedRoot := merkle.ComputeMerkleRoot(leaves)

	if status != "Closed" && status != "Archived" {
		c.JSON(http.StatusOK, gin.H{
			"status":            "in_progress",
			"ballots_count":     len(leaves),
			"leaves":            leaves,
			"computed_root":     computedRoot,
			"signed_merkle_root": signedRoot.String,
			"verification_pass": true,
			"message":           "Election is active; Merkle root represents active state.",
		})
		return
	}

	if !signedRoot.Valid || signedRoot.String == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":            "error",
			"ballots_count":     len(leaves),
			"leaves":            leaves,
			"computed_root":     computedRoot,
			"verification_pass": false,
			"message":           "No signed Merkle Root archived for this closed election.",
		})
		return
	}

	// Verify the computed root against the signed root using the system's public key
	err = crypto.VerifyRSASignature(systemPublicKeyPEM, []byte(computedRoot), signedRoot.String)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":            "tampered",
			"ballots_count":     len(leaves),
			"leaves":            leaves,
			"computed_root":     computedRoot,
			"signed_merkle_root": signedRoot.String,
			"verification_pass": false,
			"message":           "WARNING: Merkle root signature verification failed. Ledger has been modified!",
		})
		db.LogEvent("WARNING: Merkle integrity check failed for election: "+id, "System", c.ClientIP())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            "verified",
		"ballots_count":     len(leaves),
		"leaves":            leaves,
		"computed_root":     computedRoot,
		"signed_merkle_root": signedRoot.String,
		"verification_pass": true,
		"message":           "Merkle Tree ledger integrity verified successfully.",
	})
}

func tallyVotes(c *gin.Context) {
	id := c.Param("id")
	var req TallyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Verify election is Closed
	var status, encPrivKey, salt, signedRoot string
	var sRoot sql.NullString
	query := "SELECT Status, EncryptedPrivateKey, PrivateKeySalt, SignedMerkleRoot FROM Elections WHERE ElectionID = $1"
	err := db.DB.QueryRow(query, id).Scan(&status, &encPrivKey, &salt, &sRoot)
	if err != nil {
		query = "SELECT Status, EncryptedPrivateKey, PrivateKeySalt, SignedMerkleRoot FROM Elections WHERE ElectionID = ?"
		err = db.DB.QueryRow(query, id).Scan(&status, &encPrivKey, &salt, &sRoot)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if status != "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Election must be closed to decrypt and tally"})
		return
	}

	if sRoot.Valid {
		signedRoot = sRoot.String
	}

	// 2. Verify Merkle Tree integrity (REQ-CNT-1)
	var leaves []string
	vRows, err := db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp ASC", id)
	if err != nil {
		vRows, err = db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp ASC", id)
	}

	if err == nil {
		defer vRows.Close()
		for vRows.Next() {
			var lh string
			if err := vRows.Scan(&lh); err == nil {
				leaves = append(leaves, lh)
			}
		}
	}

	computedRoot := merkle.ComputeMerkleRoot(leaves)
	if signedRoot != "" {
		err = crypto.VerifyRSASignature(systemPublicKeyPEM, []byte(computedRoot), signedRoot)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tally aborted: Merkle Tree integrity check failed (tampering detected)"})
			db.LogEvent("Tally attempt aborted due to Merkle mismatch for election: "+id, "CountingAuthority", c.ClientIP())
			return
		}
	}

	// 3. Decrypt Election Private Key using Counting Authority's passphrase (TBD-1 / REQ-ENC-3)
	aesKey := crypto.DeriveKey(req.Passphrase, salt)
	privKeyBytes, err := crypto.DecryptAES(aesKey, encPrivKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Counting Authority passphrase"})
		db.LogEvent("Failed tally attempt: incorrect passphrase for election: "+id, "CountingAuthority", c.ClientIP())
		return
	}

	privateKeyPEM := string(privKeyBytes)

	// 4. Fetch encrypted ballots, decrypt them and tally (REQ-CNT-2 / REQ-CNT-3)
	ballotsQuery := "SELECT EncryptedBallot FROM Votes WHERE ElectionID = $1"
	rows, err := db.DB.Query(ballotsQuery, id)
	if err != nil {
		ballotsQuery = "SELECT EncryptedBallot FROM Votes WHERE ElectionID = ?"
		rows, err = db.DB.Query(ballotsQuery, id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error reading ballots"})
		return
	}
	defer rows.Close()

	tallyMap := make(map[string]int)
	totalVotes := 0

	for rows.Next() {
		var encBallot string
		if err := rows.Scan(&encBallot); err == nil {
			decBytes, err := crypto.DecryptRSA(privateKeyPEM, encBallot)
			if err == nil {
				candidateID := string(decBytes)
				tallyMap[candidateID]++
				totalVotes++
			} else {
				log.Printf("Failed to decrypt ballot: %v", err)
			}
		}
	}

	// Update status to Archived (Tally complete)
	updateQuery := "UPDATE Elections SET Status = 'Archived' WHERE ElectionID = $1"
	_, _ = db.DB.Exec(updateQuery, id)

	// Fetch candidates details to return user-friendly results
	var results []gin.H
	cRows, err := db.DB.Query("SELECT CandidateID, CandidateName, Party FROM Candidates WHERE ElectionID = $1", id)
	if err != nil {
		cRows, err = db.DB.Query("SELECT CandidateID, CandidateName, Party FROM Candidates WHERE ElectionID = ?", id)
	}

	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var cid, name, party string
			if err := cRows.Scan(&cid, &name, &party); err == nil {
				votes := tallyMap[cid]
				results = append(results, gin.H{
					"candidate_id":   cid,
					"candidate_name": name,
					"party":          party,
					"votes":          votes,
				})
			}
		}
	}

	// Store results in an in-memory or static file or just return it. 
	// For simplicity and since we transition the status to 'Archived', we can return it.
	// We can save the tally results in a JSON file or keep it in database.
	// Wait, we can save the results back to a new table or field, but keeping it in the database via status update is standard,
	// and we can dynamically count or store a serialized version in Elections table if we want. But dynamic count post-decryption is not possible once memory is cleared, unless we store the decrypted counts. Let's serialize the results and store them in the Elections table, or just write them to a results file/tally table!
	// Oh, we can save the results directly in `Elections.Description` or store them in a new table, or write a JSON results file.
	// Writing a JSON results file `results_<electionID>.json` is extremely clean! Or we can serialize results and save them to a new column or table.
	// Actually, let's create a table `ElectionResults` or just add a `ResultsJSON` column to the `Elections` table!
	// Adding a `ResultsJSON` column to `Elections` table: we can alter schema or just write to a new table. Let's add a table `ElectionResults` in our schema, or let's create it on the fly, or store the results in the `Elections` table. Wait! Let's check: can we just write it to a database table `ElectionResults`?
	// Yes! In `db.go`, let's see if we need it. We can just create `ElectionResults` table dynamically or save it to a JSON results file.
	// Let's create `CREATE TABLE IF NOT EXISTS ElectionResults (ElectionID VARCHAR(255) PRIMARY KEY, ResultsJSON TEXT)` in `db.go` migrations or just run it.
	// Let's check, does SQLite/Postgres support it? Yes!
	// Let's write the results to `ElectionResults` table.
	resultsJSON, _ := json.Marshal(results)
	// Wait, since we already have the Election, we can just save it into a new table `ElectionResults` or overwrite the Description field with results, or just create the table.
	// Let's execute:
	db.DB.Exec("CREATE TABLE IF NOT EXISTS ElectionResults (ElectionID VARCHAR(255) PRIMARY KEY, ResultsJSON TEXT)")
	insertRes := "INSERT INTO ElectionResults (ElectionID, ResultsJSON) VALUES ($1, $2)"
	_, err = db.DB.Exec(insertRes, id, string(resultsJSON))
	if err != nil {
		insertRes = "INSERT INTO ElectionResults (ElectionID, ResultsJSON) VALUES (?, ?)"
		_, err = db.DB.Exec(insertRes, id, string(resultsJSON))
	}

	db.LogEvent("Tally completed and archived for election: "+id, "CountingAuthority", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"election_id": id,
		"total_votes": totalVotes,
		"results":     results,
	})
}

func getResults(c *gin.Context) {
	id := c.Param("id")

	var status string
	err := db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = $1", id).Scan(&status)
	if err != nil {
		err = db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = ?", id).Scan(&status)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if status != "Archived" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Results are not yet tallied or published for this election"})
		return
	}

	var resultsJSON string
	err = db.DB.QueryRow("SELECT ResultsJSON FROM ElectionResults WHERE ElectionID = $1", id).Scan(&resultsJSON)
	if err != nil {
		err = db.DB.QueryRow("SELECT ResultsJSON FROM ElectionResults WHERE ElectionID = ?", id).Scan(&resultsJSON)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Results not found"})
		return
	}

	var results []interface{}
	json.Unmarshal([]byte(resultsJSON), &results)

	c.JSON(http.StatusOK, gin.H{
		"election_id": id,
		"results":     results,
	})
}
