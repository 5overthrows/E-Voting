package token

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"evoting/shared/crypto"
	"evoting/shared/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateRequest struct {
	ElectionID string `json:"election_id" binding:"required"`
}

type ValidateRequest struct {
	Token string `json:"token" binding:"required"`
}

type InvalidateRequest struct {
	TokenHash string `json:"token_hash" binding:"required"`
}

// RegisterRoutes registers the token service routes on the given router
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/token/generate", generateToken)
	r.POST("/api/token/validate", validateToken)
	r.POST("/api/token/invalidate", invalidateToken)
}

// StartTokenServer starts the token service on the specified address
func StartTokenServer(addr string) {
	r := gin.Default()

	// CORS config (allows requests from Auth/Voting service and client if needed)
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

	log.Printf("Anonymous Token Service running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Token server failed: %v", err)
	}
}

func generateToken(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create random token
	rawToken, err := crypto.GenerateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cryptographically secure token"})
		return
	}

	tokenHash := crypto.HashToken(rawToken)
	tokenID := uuid.New().String()
	issuedTime := time.Now()
	expiryTime := issuedTime.Add(15 * time.Minute) // valid for 15 minutes

	query := "INSERT INTO Tokens (TokenID, TokenHash, IssuedTime, ExpiryTime, Used) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.DB.Exec(query, tokenID, tokenHash, issuedTime, expiryTime, false)
	if err != nil {
		query = "INSERT INTO Tokens (TokenID, TokenHash, IssuedTime, ExpiryTime, Used) VALUES (?, ?, ?, ?, ?)"
		_, err = db.DB.Exec(query, tokenID, tokenHash, issuedTime, expiryTime, false)
	}

	if err != nil {
		log.Printf("Failed to insert token in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	db.LogEvent("Anonymous voting token hash created in database", "TokenService", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"token":      rawToken,
		"token_hash": tokenHash,
	})
}

func validateToken(c *gin.Context) {
	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenHash := crypto.HashToken(req.Token)

	var expiryTime time.Time
	var used bool
	query := "SELECT ExpiryTime, Used FROM Tokens WHERE TokenHash = $1"
	err := db.DB.QueryRow(query, tokenHash).Scan(&expiryTime, &used)
	if err != nil {
		query = "SELECT ExpiryTime, Used FROM Tokens WHERE TokenHash = ?"
		err = db.DB.QueryRow(query, tokenHash).Scan(&expiryTime, &used)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"valid": false, "reason": "Token does not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if used {
		c.JSON(http.StatusOK, gin.H{"valid": false, "reason": "Token has already been used"})
		return
	}

	if time.Now().After(expiryTime) {
		c.JSON(http.StatusOK, gin.H{"valid": false, "reason": "Token has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"token_hash": tokenHash,
	})
}

func invalidateToken(c *gin.Context) {
	var req InvalidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE Tokens SET Used = TRUE WHERE TokenHash = $1"
	_, err := db.DB.Exec(query, req.TokenHash)
	if err != nil {
		query = "UPDATE Tokens SET Used = TRUE WHERE TokenHash = ?"
		_, err = db.DB.Exec(query, req.TokenHash)
	}

	if err != nil {
		log.Printf("Failed to invalidate token hash: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	db.LogEvent("Anonymous voting token marked as used", "TokenService", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"status": "invalidated"})
}
