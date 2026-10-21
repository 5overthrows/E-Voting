package voting

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"evoting/services/auth"
	"evoting/shared/crypto"
	"evoting/shared/db"
	"evoting/shared/merkle"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VoteRequest struct {
	Token       string `json:"token" binding:"required"`
	CandidateID string `json:"candidate_id" binding:"required"`
	ElectionID  string `json:"election_id" binding:"required"`
}

// RegisterRoutes registers the voting service routes on the given router
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/vote/cast", castVote)
}

// StartVotingServer starts the voting service on the specified address
func StartVotingServer(addr string) {
	r := gin.Default()

	// Configure CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	RegisterRoutes(r)

	log.Printf("Voting Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Voting server failed: %v", err)
	}
}

func castVote(c *gin.Context) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Verify candidate exists for the active election
	var candName string
	queryCand := "SELECT CandidateName FROM Candidates WHERE CandidateID = $1 AND ElectionID = $2"
	err := db.DB.QueryRow(queryCand, req.CandidateID, req.ElectionID).Scan(&candName)
	if err != nil {
		queryCand = "SELECT CandidateName FROM Candidates WHERE CandidateID = ? AND ElectionID = ?"
		err = db.DB.QueryRow(queryCand, req.CandidateID, req.ElectionID).Scan(&candName)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Candidate not registered in this election"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 2. Verify election is active and retrieve public key
	var status, publicKey string
	queryElect := "SELECT Status, PublicKey FROM Elections WHERE ElectionID = $1"
	err = db.DB.QueryRow(queryElect, req.ElectionID).Scan(&status, &publicKey)
	if err != nil {
		queryElect = "SELECT Status, PublicKey FROM Elections WHERE ElectionID = ?"
		err = db.DB.QueryRow(queryElect, req.ElectionID).Scan(&status, &publicKey)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if status != "Active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Election is not active"})
		return
	}

	// 3. Validate Token via Token Service
	tokenServiceURL := os.Getenv("TOKEN_VALIDATE_URL")
	if tokenServiceURL == "" {
		tokenServiceURL = "http://127.0.0.1:8082/api/token/validate"
	}
	tokenReqBody, _ := json.Marshal(map[string]string{
		"token": req.Token,
	})

	resp, err := http.Post(tokenServiceURL, "application/json", bytes.NewBuffer(tokenReqBody))
	if err != nil {
		// Fallback to unified port 8080
		fallbackURL := "http://127.0.0.1:8080/api/token/validate"
		if resp2, err2 := http.Post(fallbackURL, "application/json", bytes.NewBuffer(tokenReqBody)); err2 == nil {
			resp = resp2
			err = nil
		} else {
			log.Printf("Token service connection error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token Service unavailable"})
			return
		}
	}
	defer resp.Body.Close()

	var tokenResp struct {
		Valid     bool   `json:"valid"`
		Reason    string `json:"reason"`
		TokenHash string `json:"token_hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token service response"})
		return
	}

	// REQ-VOTE-2: Verify token validity (existence, non-expiry, unused status)
	if !tokenResp.Valid {
		// Reject without indicating specific check detail to prevent leaks (REQ-VOTE-2 / REQ-VOTE-6)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid, expired, or reused voting token"})
		db.LogEvent("Failed vote casting attempt: invalid token: "+tokenResp.Reason, "Voter", c.ClientIP())
		return
	}

	// 4. Encrypt Ballot (candidate_id) using public key (RSA-OAEP)
	encryptedBallot, err := crypto.EncryptRSA(publicKey, []byte(req.CandidateID))
	if err != nil {
		log.Printf("Ballot encryption failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt ballot"})
		return
	}

	// 5. Compute SHA-256 leaf hash for Merkle Tree
	leafHashBytes := sha256.Sum256([]byte(encryptedBallot))
	leafHash := fmt.Sprintf("%x", leafHashBytes)

	// 6. Persist encrypted ballot in database (append-only)
	voteID := uuid.New().String()
	timestamp := time.Now()
	queryVote := "INSERT INTO Votes (VoteID, ElectionID, EncryptedBallot, Timestamp, MerkleLeafHash) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.DB.Exec(queryVote, voteID, req.ElectionID, encryptedBallot, timestamp, leafHash)
	if err != nil {
		queryVote = "INSERT INTO Votes (VoteID, ElectionID, EncryptedBallot, Timestamp, MerkleLeafHash) VALUES (?, ?, ?, ?, ?)"
		_, err = db.DB.Exec(queryVote, voteID, req.ElectionID, encryptedBallot, timestamp, leafHash)
	}

	if err != nil {
		log.Printf("Failed to store vote: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error storing vote"})
		return
	}

	// 7. Incrementally compute current Merkle Root for monitoring/logging
	var leaves []string
	rows, err := db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = $1 ORDER BY Timestamp ASC", req.ElectionID)
	if err != nil {
		rows, err = db.DB.Query("SELECT MerkleLeafHash FROM Votes WHERE ElectionID = ? ORDER BY Timestamp ASC", req.ElectionID)
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
	currentRoot := merkle.ComputeMerkleRoot(leaves)
	log.Printf("Vote cast! Leaf count: %d. Current Merkle Root: %s", len(leaves), currentRoot)

	// 8. Invalidate Token via Token Service
	invalidateURL := os.Getenv("TOKEN_INVALIDATE_URL")
	if invalidateURL == "" {
		invalidateURL = "http://127.0.0.1:8082/api/token/invalidate"
	}
	invalidateReqBody, _ := json.Marshal(map[string]string{
		"token_hash": tokenResp.TokenHash,
	})
	
	// Issue internal call, log error but don't block user if this is just async state update
	if invResp, err := http.Post(invalidateURL, "application/json", bytes.NewBuffer(invalidateReqBody)); err == nil {
		invResp.Body.Close()
	} else {
		// Try fallback to unified port 8080
		fallbackURL := "http://127.0.0.1:8080/api/token/invalidate"
		if invResp2, err2 := http.Post(fallbackURL, "application/json", bytes.NewBuffer(invalidateReqBody)); err2 == nil {
			invResp2.Body.Close()
		} else {
			log.Printf("Warning: Failed to invalidate token: %v", err)
		}
	}

	// 9. Confirm Vote to Auth Service to mark voter as HasVoted = true
	if err := auth.CallConfirmVote(tokenResp.TokenHash); err != nil {
		log.Printf("Warning: Failed to confirm vote with auth service: %v", err)
	}

	db.LogEvent("Encrypted vote added to ledger. Merkle Leaf: "+leafHash, "Voter", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Vote cast successfully",
		"timestamp":  timestamp,
		"receipt_h":  leafHash,
		"root_hash":  currentRoot,
	})
}
