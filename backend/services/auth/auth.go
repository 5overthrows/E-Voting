package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"evoting/shared/crypto"
	"evoting/shared/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey            = []byte("secret_key_voting_system")
	pendingTokens     = make(map[string]string) // TokenHash -> UserID
	pendingTokensMu   sync.RWMutex
	LockoutThreshold  = 3
	LockoutDuration   = 5 * time.Minute
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenRequest struct {
	ElectionID string `json:"election_id" binding:"required"`
}

type ConfirmRequest struct {
	TokenHash string `json:"token_hash" binding:"required"`
}

// GenerateJWT creates a JWT for a user
func GenerateJWT(userID, role, name string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT validates a token string and returns the claims
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

// AuthMiddleware checks JWT authentication
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr := authHeader[7:]
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if len(allowedRoles) > 0 {
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("name", claims.Name)
		c.Next()
	})
}

// AddPendingToken registers a token hash with a user ID in memory
func AddPendingToken(tokenHash, userID string) {
	pendingTokensMu.Lock()
	defer pendingTokensMu.Unlock()
	pendingTokens[tokenHash] = userID
}

// RegisterRoutes registers the authentication service routes on the given router
func RegisterRoutes(r *gin.Engine) {
	r.POST("/api/auth/login", login)
	r.POST("/api/auth/confirm-vote", confirmVote) // Internal API

	// Protected routes
	protected := r.Group("/api/auth")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/me", me)
		protected.POST("/token/request", requestToken)
	}
}

// StartAuthServer starts the authentication service on the specified address
func StartAuthServer(addr string) {
	r := gin.Default()

	// Configure CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	RegisterRoutes(r)

	log.Printf("Authentication Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Auth server failed: %v", err)
	}
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var name, passwordHash, role string
	var failedAttempts int
	var lockedUntil sql.NullTime

	query := "SELECT Name, PasswordHash, Role, FailedLoginAttempts, LockedUntil FROM Users WHERE UserID = $1"
	err := db.DB.QueryRow(query, req.UserID).Scan(&name, &passwordHash, &role, &failedAttempts, &lockedUntil)
	if err == sql.ErrNoRows {
		// Try sqlite syntax if query above failed, though placeholder $1 is supported.
		query = "SELECT Name, PasswordHash, Role, FailedLoginAttempts, LockedUntil FROM Users WHERE UserID = ?"
		err = db.DB.QueryRow(query, req.UserID).Scan(&name, &passwordHash, &role, &failedAttempts, &lockedUntil)
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID or Password"})
		db.LogEvent("Failed login attempt: non-existent User ID: "+req.UserID, "None", c.ClientIP())
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if account is locked
	if lockedUntil.Valid && lockedUntil.Time.After(time.Now()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is temporarily locked. Try again later."})
		db.LogEvent("Failed login attempt: locked account for User ID: "+req.UserID, role, c.ClientIP())
		return
	}

	// Verify password
	if !crypto.CheckPasswordHash(req.Password, passwordHash) {
		failedAttempts++
		var lockQuery string
		var lockErr error
		if failedAttempts >= LockoutThreshold {
			unlockTime := time.Now().Add(LockoutDuration)
			lockQuery = "UPDATE Users SET FailedLoginAttempts = $1, LockedUntil = $2 WHERE UserID = $3"
			_, lockErr = db.DB.Exec(lockQuery, failedAttempts, unlockTime, req.UserID)
			if lockErr != nil {
				lockQuery = "UPDATE Users SET FailedLoginAttempts = ?, LockedUntil = ? WHERE UserID = ?"
				db.DB.Exec(lockQuery, failedAttempts, unlockTime, req.UserID)
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "Account locked due to too many failed attempts."})
			db.LogEvent("Account locked for User ID: "+req.UserID, role, c.ClientIP())
		} else {
			lockQuery = "UPDATE Users SET FailedLoginAttempts = $1 WHERE UserID = $2"
			_, lockErr = db.DB.Exec(lockQuery, failedAttempts, req.UserID)
			if lockErr != nil {
				lockQuery = "UPDATE Users SET FailedLoginAttempts = ? WHERE UserID = ?"
				db.DB.Exec(lockQuery, failedAttempts, req.UserID)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID or Password"})
			db.LogEvent("Failed login attempt: wrong password for User ID: "+req.UserID, role, c.ClientIP())
		}
		return
	}

	// Reset failed attempts on success
	resetQuery := "UPDATE Users SET FailedLoginAttempts = 0, LockedUntil = NULL WHERE UserID = $1"
	_, resetErr := db.DB.Exec(resetQuery, req.UserID)
	if resetErr != nil {
		resetQuery = "UPDATE Users SET FailedLoginAttempts = 0, LockedUntil = NULL WHERE UserID = ?"
		db.DB.Exec(resetQuery, req.UserID)
	}

	// Create JWT token
	token, err := GenerateJWT(req.UserID, role, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	db.LogEvent("User login successful: "+req.UserID, role, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"user_id": req.UserID,
			"name":    name,
			"role":    role,
		},
	})
}

func me(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	name := c.GetString("name")

	var hasVoted bool
	query := "SELECT HasVoted FROM Users WHERE UserID = $1"
	err := db.DB.QueryRow(query, userID).Scan(&hasVoted)
	if err != nil {
		query = "SELECT HasVoted FROM Users WHERE UserID = ?"
		db.DB.QueryRow(query, userID).Scan(&hasVoted)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"name":      name,
		"role":      role,
		"has_voted": hasVoted,
	})
}

func requestToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	role := c.GetString("role")

	// 1. Verify election status
	var status string
	err := db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = $1", req.ElectionID).Scan(&status)
	if err != nil {
		err = db.DB.QueryRow("SELECT Status FROM Elections WHERE ElectionID = ?", req.ElectionID).Scan(&status)
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

	// 2. Verify eligibility (HasVoted flag)
	var hasVoted bool
	err = db.DB.QueryRow("SELECT HasVoted FROM Users WHERE UserID = $1", userID).Scan(&hasVoted)
	if err != nil {
		err = db.DB.QueryRow("SELECT HasVoted FROM Users WHERE UserID = ?", userID).Scan(&hasVoted)
	}

	if hasVoted {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have already cast a ballot in this election"})
		db.LogEvent("Failed token request: already voted for User: "+userID, role, c.ClientIP())
		return
	}

	// 3. Request Token Service to generate a token
	tokenServiceURL := os.Getenv("TOKEN_SERVICE_URL")
	if tokenServiceURL == "" {
		tokenServiceURL = "http://127.0.0.1:8082/api/token/generate"
	}
	tokenReqBody, _ := json.Marshal(map[string]string{
		"election_id": req.ElectionID,
	})

	resp, err := http.Post(tokenServiceURL, "application/json", bytes.NewBuffer(tokenReqBody))
	if err != nil {
		// Fallback to unified port 8080 if standalone token service is unreachable
		fallbackURL := "http://127.0.0.1:8080/api/token/generate"
		if resp2, err2 := http.Post(fallbackURL, "application/json", bytes.NewBuffer(tokenReqBody)); err2 == nil {
			resp = resp2
			err = nil
		} else {
			log.Printf("Token Service connection failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token Service unavailable"})
			return
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		c.JSON(resp.StatusCode, gin.H{"error": errResp["error"]})
		return
	}

	var tokenResp struct {
		Token     string `json:"token"`
		TokenHash string `json:"token_hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Token Service"})
		return
	}

	// 4. Map TokenHash -> UserID in memory
	AddPendingToken(tokenResp.TokenHash, userID)

	db.LogEvent("Anonymous token issued to client for user: "+userID, role, c.ClientIP())

	// Return ONLY the raw token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": tokenResp.Token,
	})
}

// confirmVote is called internally by the Voting Server once a vote is written
func confirmVote(c *gin.Context) {
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Internal lookup
	pendingTokensMu.Lock()
	userID, ok := pendingTokens[req.TokenHash]
	if ok {
		delete(pendingTokens, req.TokenHash)
	}
	pendingTokensMu.Unlock()

	if !ok {
		// Fail silently or log internally to prevent leaks
		log.Printf("confirmVote failed: token hash not found in pending list: %s", req.TokenHash)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pending token not found"})
		return
	}

	// Mark user as HasVoted = true in DB
	query := "UPDATE Users SET HasVoted = TRUE WHERE UserID = $1"
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		query = "UPDATE Users SET HasVoted = TRUE WHERE UserID = ?"
		_, err = db.DB.Exec(query, userID)
	}

	if err != nil {
		log.Printf("Failed to update HasVoted for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voter status"})
		return
	}

	db.LogEvent("Vote cast successfully, voter marked: "+userID, "Voter", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"status": "confirmed"})
}

// CallConfirmVote issues a callback to Auth Service to confirm a vote
func CallConfirmVote(tokenHash string) error {
	confirmURL := os.Getenv("AUTH_CONFIRM_URL")
	if confirmURL == "" {
		confirmURL = "http://127.0.0.1:8081/api/auth/confirm-vote"
	}
	reqBody, _ := json.Marshal(ConfirmRequest{TokenHash: tokenHash})
	
	resp, err := http.Post(confirmURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		// Fallback to unified port 8080
		fallbackURL := "http://127.0.0.1:8080/api/auth/confirm-vote"
		if resp2, err2 := http.Post(fallbackURL, "application/json", bytes.NewBuffer(reqBody)); err2 == nil {
			resp = resp2
			err = nil
		} else {
			return err
		}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
