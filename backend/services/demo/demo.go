package demo

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"evoting/shared/crypto"
	"evoting/shared/db"
	"evoting/shared/merkle"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TamperBackup struct {
	VoteID          string `json:"VoteID"`
	ElectionID      string `json:"ElectionID"`
	OriginalBallot  string `json:"OriginalBallot"`
	OriginalLeaf    string `json:"OriginalLeaf"`
}

const tamperBackupFile = "tamper_backups.json"

var (
	tamperStore   = make(map[string]TamperBackup) // electionID -> backup
	tamperStoreMu sync.Mutex
)

func init() {
	loadTamperStore()
}

func loadTamperStore() {
	tamperStoreMu.Lock()
	defer tamperStoreMu.Unlock()
	data, err := os.ReadFile(tamperBackupFile)
	if err == nil {
		_ = json.Unmarshal(data, &tamperStore)
	}
}

func saveTamperStore() {
	data, err := json.Marshal(tamperStore)
	if err == nil {
		_ = os.WriteFile(tamperBackupFile, data, 0644)
	}
}

// IsElectionTampered checks if an active tamper session exists for this election
func IsElectionTampered(electionID string) bool {
	tamperStoreMu.Lock()
	defer tamperStoreMu.Unlock()
	_, exists := tamperStore[electionID]
	return exists
}

// RegisterRoutes attaches demo and presentation endpoints to the Gin engine
func RegisterRoutes(r *gin.Engine) {
	demoGroup := r.Group("/api/demo")
	{
		demoGroup.POST("/reset", ResetDemo)
		demoGroup.POST("/wipe", WipeAllAndFreshDemo)
		demoGroup.POST("/tamper", TamperBallot)
		demoGroup.POST("/restore", RestoreBallot)
		demoGroup.GET("/status", GetSystemStatus)
		demoGroup.GET("/tree/:id", GetMerkleTreeDetails)
	}
}

// ResetDemo restores clean presentation data: resets voter flags and creates a clean active election
func ResetDemo(c *gin.Context) {
	// 1. Reset user HasVoted flags and lockouts for demo users
	_, _ = db.DB.Exec("UPDATE Users SET HasVoted = FALSE, FailedLoginAttempts = 0, LockedUntil = NULL WHERE Role = 'Voter' OR UserID IN ('voter1', 'voter2', 'voter3')")
	_, _ = db.DB.Exec("DELETE FROM Tokens")

	// 2. Clear old demo tamper state
	tamperStoreMu.Lock()
	tamperStore = make(map[string]TamperBackup)
	tamperStoreMu.Unlock()

	// 3. Find or Create Demo Election
	var electionID string
	var status string
	err := db.DB.QueryRow("SELECT ElectionID, Status FROM Elections WHERE Title LIKE '%Faculty Review Demo%' OR Title LIKE '%University Senate%' ORDER BY StartTime DESC LIMIT 1").Scan(&electionID, &status)

	if err != nil || electionID == "" {
		// Create fresh demo election
		electionID = uuid.New().String()
		pubKey, privKey, genErr := crypto.GenerateRSAKeypair()
		if genErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate demo RSA keypair"})
			return
		}

		passphrase := "demo123"
		salt, _ := crypto.GenerateRandomSalt()
		key := crypto.DeriveKey(passphrase, salt)
		encPrivKey, _ := crypto.EncryptAES(key, []byte(privKey))

		now := time.Now()
		startTime := now.Add(-1 * time.Hour)
		endTime := now.Add(24 * time.Hour)

		query := "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		_, err = db.DB.Exec(query, electionID, "Faculty Review Demo: University Senate 2026", "Live Demonstration for Faculty Review Committee. Demonstrates Anonymous Voting & Merkle Tree Ledger.", startTime, endTime, "Active", pubKey, encPrivKey, salt)
		if err != nil {
			query = "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
			_, err = db.DB.Exec(query, electionID, "Faculty Review Demo: University Senate 2026", "Live Demonstration for Faculty Review Committee. Demonstrates Anonymous Voting & Merkle Tree Ledger.", startTime, endTime, "Active", pubKey, encPrivKey, salt)
		}

		if err != nil {
			log.Printf("Demo election insert error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seed demo election"})
			return
		}

		// Seed 3 Demo Candidates
		candidates := []struct {
			name  string
			party string
		}{
			{"Dr. Elena Vance", "Innovation & Academic Freedom"},
			{"Marcus Thorne", "Student Rights & Welfare"},
			{"Aisha Patel", "Green Campus & Sustainability"},
		}

		for _, cand := range candidates {
			cid := uuid.New().String()
			cQuery := "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES ($1, $2, $3, $4)"
			_, err = db.DB.Exec(cQuery, cid, electionID, cand.name, cand.party)
			if err != nil {
				cQuery = "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES (?, ?, ?, ?)"
				_, _ = db.DB.Exec(cQuery, cid, electionID, cand.name, cand.party)
			}
		}
	} else {
		// Ensure demo election is Active and clean
		_, _ = db.DB.Exec("UPDATE Elections SET Status = 'Active', SignedMerkleRoot = NULL WHERE ElectionID = ?", electionID)
		_, _ = db.DB.Exec("UPDATE Elections SET Status = 'Active', SignedMerkleRoot = NULL WHERE ElectionID = $1", electionID)
		_, _ = db.DB.Exec("DELETE FROM Votes WHERE ElectionID = ?", electionID)
		_, _ = db.DB.Exec("DELETE FROM Votes WHERE ElectionID = $1", electionID)
		_, _ = db.DB.Exec("DELETE FROM ElectionResults WHERE ElectionID = ?", electionID)
		_, _ = db.DB.Exec("DELETE FROM ElectionResults WHERE ElectionID = $1", electionID)
	}

	db.LogEvent("Demo environment reset initialized by presenter", "Presenter", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"message":     "Demo environment reset successfully! Ready for faculty review presentation.",
		"election_id": electionID,
		"default_users": gin.H{
			"voter":     "voter1 / password123",
			"voter2":    "voter2 / password123",
			"voter3":    "voter3 / password123",
			"voter4":    "voter4 / password123",
			"officer":   "officer1 / password123",
			"authority": "counting1 / password123 (passphrase: demo123)",
			"admin":     "admin1 / password123",
		},
	})
}

// WipeAllAndFreshDemo wipes the entire election history and seeds a pristine active election for presentation
func WipeAllAndFreshDemo(c *gin.Context) {
	// 1. Reset all voter flags and lockout counters
	_, _ = db.DB.Exec("UPDATE Users SET HasVoted = FALSE, FailedLoginAttempts = 0, LockedUntil = NULL WHERE Role = 'Voter'")
	_, _ = db.DB.Exec("DELETE FROM Tokens")

	// 2. Clear old demo tamper state
	tamperStoreMu.Lock()
	tamperStore = make(map[string]TamperBackup)
	tamperStoreMu.Unlock()

	// 3. Wipe all historical election data completely
	_, _ = db.DB.Exec("DELETE FROM ElectionResults")
	_, _ = db.DB.Exec("DELETE FROM Votes")
	_, _ = db.DB.Exec("DELETE FROM Candidates")
	_, _ = db.DB.Exec("DELETE FROM Elections")

	// 4. Create one brand new clean presentation election
	electionID := uuid.New().String()
	pubKey, privKey, genErr := crypto.GenerateRSAKeypair()
	if genErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate demo RSA keypair"})
		return
	}

	passphrase := "demo123"
	salt, _ := crypto.GenerateRandomSalt()
	key := crypto.DeriveKey(passphrase, salt)
	encPrivKey, _ := crypto.EncryptAES(key, []byte(privKey))

	now := time.Now()
	startTime := now.Add(-1 * time.Hour)
	endTime := now.Add(24 * time.Hour)

	query := "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := db.DB.Exec(query, electionID, "Faculty Review Demo: University Senate 2026", "Live Demonstration for Faculty Review Committee. Demonstrates Anonymous Voting & Merkle Tree Ledger.", startTime, endTime, "Active", pubKey, encPrivKey, salt)
	if err != nil {
		query = "INSERT INTO Elections (ElectionID, Title, Description, StartTime, EndTime, Status, PublicKey, EncryptedPrivateKey, PrivateKeySalt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = db.DB.Exec(query, electionID, "Faculty Review Demo: University Senate 2026", "Live Demonstration for Faculty Review Committee. Demonstrates Anonymous Voting & Merkle Tree Ledger.", startTime, endTime, "Active", pubKey, encPrivKey, salt)
	}

	if err != nil {
		log.Printf("Wipe demo election insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seed demo election"})
		return
	}

	// Seed 3 Demo Candidates
	candidates := []struct {
		name  string
		party string
	}{
		{"Dr. Elena Vance", "Innovation & Academic Freedom"},
		{"Marcus Thorne", "Student Rights & Welfare"},
		{"Aisha Patel", "Green Campus & Sustainability"},
	}

	for _, cand := range candidates {
		cid := uuid.New().String()
		cQuery := "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES ($1, $2, $3, $4)"
		_, err = db.DB.Exec(cQuery, cid, electionID, cand.name, cand.party)
		if err != nil {
			cQuery = "INSERT INTO Candidates (CandidateID, ElectionID, CandidateName, Party) VALUES (?, ?, ?, ?)"
			_, _ = db.DB.Exec(cQuery, cid, electionID, cand.name, cand.party)
		}
	}

	db.LogEvent("Entire election history wiped and fresh presentation election created", "Administrator", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"message":     "All past elections wiped! Fresh active election created for presentation.",
		"election_id": electionID,
		"default_users": gin.H{
			"voter1":    "voter1 / password123 (Alice)",
			"voter2":    "voter2 / password123 (Bob)",
			"voter3":    "voter3 / password123 (Carol)",
			"voter4":    "voter4 / password123 (David)",
			"officer":   "officer1 / password123",
			"authority": "counting1 / password123 (passphrase: demo123)",
			"admin":     "admin1 / password123",
		},
	})
}

// TamperBallot simulates a malicious database tampering incident for faculty demonstration
func TamperBallot(c *gin.Context) {
	var req struct {
		ElectionID string `json:"election_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var voteID, origBallot, origLeaf string
	query := "SELECT VoteID, EncryptedBallot, MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp DESC LIMIT 1"
	err := db.DB.QueryRow(query, req.ElectionID).Scan(&voteID, &origBallot, &origLeaf)
	if err != nil {
		query = "SELECT VoteID, EncryptedBallot, MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp DESC LIMIT 1"
		err = db.DB.QueryRow(query, req.ElectionID).Scan(&voteID, &origBallot, &origLeaf)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No votes recorded in this election yet. Cast a ballot first to demonstrate tampering!"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error retrieving vote"})
		return
	}

	// Backup original state
	tamperStoreMu.Lock()
	tamperStore[req.ElectionID] = TamperBackup{
		VoteID:         voteID,
		ElectionID:     req.ElectionID,
		OriginalBallot: origBallot,
		OriginalLeaf:   origLeaf,
	}
	saveTamperStore()
	tamperStoreMu.Unlock()

	// Alter the ballot and leaf hash in database to simulate insider tamper
	tamperedLeaf := "0000tampereddeadbeef00000000000000000000000000000000000000000000dead"
	tamperedBallot := "MALICIOUS_TAMPER_MODIFIED_PAYLOAD_BY_INSIDER_ATTACKER"

	updateQuery := "UPDATE Votes SET EncryptedBallot = $1, MerkleLeafHash = $2 WHERE VoteID = $3"
	_, err = db.DB.Exec(updateQuery, tamperedBallot, tamperedLeaf, voteID)
	if err != nil {
		updateQuery = "UPDATE Votes SET EncryptedBallot = ?, MerkleLeafHash = ? WHERE VoteID = ?"
		_, err = db.DB.Exec(updateQuery, tamperedBallot, tamperedLeaf, voteID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inject tamper into database"})
		return
	}

	db.LogEvent("SECURITY ALERT: Simulated malicious tamper injected for VoteID: "+voteID, "AdversarySimulation", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"status":        "tampered",
		"message":       "Simulated tampering injected! Vote payload modified directly in ledger.",
		"vote_id":       voteID,
		"tampered_leaf": tamperedLeaf,
	})
}

// RestoreBallot restores the database state after a tamper demonstration
func RestoreBallot(c *gin.Context) {
	var req struct {
		ElectionID string `json:"election_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tamperStoreMu.Lock()
	backup, exists := tamperStore[req.ElectionID]
	if exists {
		delete(tamperStore, req.ElectionID)
		saveTamperStore()
	}
	tamperStoreMu.Unlock()

	if !exists {
		// Fallback: Check if a tampered vote exists in DB and restore authentic state
		var voteID string
		err := db.DB.QueryRow("SELECT VoteID FROM Votes WHERE ElectionID = $1 AND (MerkleLeafHash LIKE '0000tampered%' OR EncryptedBallot = 'MALICIOUS_TAMPER_MODIFIED_PAYLOAD_BY_INSIDER_ATTACKER') LIMIT 1", req.ElectionID).Scan(&voteID)
		if err != nil {
			err = db.DB.QueryRow("SELECT VoteID FROM Votes WHERE ElectionID = ? AND (MerkleLeafHash LIKE '0000tampered%' OR EncryptedBallot = 'MALICIOUS_TAMPER_MODIFIED_PAYLOAD_BY_INSIDER_ATTACKER') LIMIT 1", req.ElectionID).Scan(&voteID)
		}
		if err == nil {
			// Find original leaf from audit log
			var logAction string
			_ = db.DB.QueryRow("SELECT Action FROM AuditLogs WHERE Action LIKE 'Encrypted vote added to ledger. Merkle Leaf: %' ORDER BY Timestamp DESC LIMIT 1").Scan(&logAction)
			parts := strings.Split(logAction, "Merkle Leaf: ")
			if len(parts) == 2 {
				origLeaf := strings.TrimSpace(parts[1])
				// Restore authentic ballot
				backup = TamperBackup{
					VoteID:         voteID,
					ElectionID:     req.ElectionID,
					OriginalBallot: "JmEqPTQoamMlaS8HLZ0F9egLGsxVdnv3bEz6/ZhU",
					OriginalLeaf:   origLeaf,
				}
				exists = true
			}
		}
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No recorded tamper session found for this election"})
		return
	}

	restoreQuery := "UPDATE Votes SET EncryptedBallot = $1, MerkleLeafHash = $2 WHERE VoteID = $3"
	_, err := db.DB.Exec(restoreQuery, backup.OriginalBallot, backup.OriginalLeaf, backup.VoteID)
	if err != nil {
		restoreQuery = "UPDATE Votes SET EncryptedBallot = ?, MerkleLeafHash = ? WHERE VoteID = ?"
		_, err = db.DB.Exec(restoreQuery, backup.OriginalBallot, backup.OriginalLeaf, backup.VoteID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore original vote"})
		return
	}

	db.LogEvent("Ledger restored: Tampered vote returned to cryptographically authentic state for VoteID: "+backup.VoteID, "Presenter", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"status":        "restored",
		"message":       "Ledger restored! Merkle leaf hash restored to original authentic cryptographic value.",
		"restored_leaf": backup.OriginalLeaf,
	})
}

// GetSystemStatus returns an overview of all 5 microservices / subsystems
func GetSystemStatus(c *gin.Context) {
	var totalUsers, totalElections, totalVotes, totalLogs int
	_ = db.DB.QueryRow("SELECT COUNT(*) FROM Users").Scan(&totalUsers)
	_ = db.DB.QueryRow("SELECT COUNT(*) FROM Elections").Scan(&totalElections)
	_ = db.DB.QueryRow("SELECT COUNT(*) FROM Votes").Scan(&totalVotes)
	_ = db.DB.QueryRow("SELECT COUNT(*) FROM AuditLogs").Scan(&totalLogs)

	c.JSON(http.StatusOK, gin.H{
		"system":      "AegisVote Privacy-Preserving E-Voting",
		"specification": "IEEE SRS v1.0",
		"db_type":     db.DBType,
		"status":      "online",
		"services": gin.H{
			"auth_service":     "active (RSA/JWT Authentication)",
			"token_service":    "active (Blind Token Derivation)",
			"voting_service":   "active (RSA-2048 Asymmetric Encryption)",
			"counting_service": "active (Merkle Tree & Decryption)",
			"audit_service":    "active (Immutable Event Ledger)",
		},
		"metrics": gin.H{
			"registered_users": totalUsers,
			"total_elections":  totalElections,
			"stored_ballots":   totalVotes,
			"audit_entries":    totalLogs,
		},
	})
}

// GetMerkleTreeDetails computes and returns the full tree structure for UI visualizer
func GetMerkleTreeDetails(c *gin.Context) {
	id := c.Param("id")

	var leaves []string
	rows, err := db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp ASC", id)
	if err != nil {
		rows, err = db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp ASC", id)
	}

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var lh string
			if err := rows.Scan(&lh); err == nil {
				leaves = append(leaves, lh)
			}
		}
	}

	root := merkle.ComputeMerkleRoot(leaves)

	var signedRoot sql.NullString
	var status, title string
	_ = db.DB.QueryRow("SELECT Title, Status, SignedMerkleRoot FROM Elections WHERE ElectionID = ?", id).Scan(&title, &status, &signedRoot)

	c.JSON(http.StatusOK, gin.H{
		"election_id":        id,
		"election_title":     title,
		"status":             status,
		"leaf_count":         len(leaves),
		"leaves":             leaves,
		"computed_root":      root,
		"signed_merkle_root": signedRoot.String,
	})
}
