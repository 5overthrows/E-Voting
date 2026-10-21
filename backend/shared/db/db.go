package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"evoting/shared/crypto"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var DBType string // "postgres" or "sqlite"

// Connect initializes the database connection and runs migrations
func Connect() {
	var err error
	dbType := os.Getenv("DB_TYPE")
	dbURL := os.Getenv("DATABASE_URL")

	if dbType == "postgres" || dbURL != "" {
		if dbURL == "" {
			dbURL = "postgres://postgres:password@localhost:5432/evoting?sslmode=disable"
		}
		log.Printf("Connecting to PostgreSQL at: %s", dbURL)
		DB, err = sql.Open("postgres", dbURL)
		DBType = "postgres"
	} else {
		sqlitePath := os.Getenv("SQLITE_PATH")
		if sqlitePath == "" {
			sqlitePath = "evoting.db"
		}
		log.Printf("Connecting to SQLite at: %s", sqlitePath)
		DB, err = sql.Open("sqlite3", sqlitePath)
		DBType = "sqlite"
	}

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Wait for db connection to be ready (useful in docker-compose)
	for i := 0; i < 10; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Database not ready yet, retrying... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection established.")
	runMigrations()
	seedData()
}

func runMigrations() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Users (
			UserID VARCHAR(255) PRIMARY KEY,
			Name VARCHAR(255) NOT NULL,
			PasswordHash VARCHAR(255) NOT NULL,
			HasVoted BOOLEAN NOT NULL DEFAULT FALSE,
			Role VARCHAR(50) NOT NULL,
			FailedLoginAttempts INTEGER DEFAULT 0,
			LockedUntil TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS Elections (
			ElectionID VARCHAR(255) PRIMARY KEY,
			Title VARCHAR(255) NOT NULL,
			Description TEXT,
			StartTime TIMESTAMP NOT NULL,
			EndTime TIMESTAMP NOT NULL,
			Status VARCHAR(50) NOT NULL,
			PublicKey TEXT,
			EncryptedPrivateKey TEXT,
			PrivateKeySalt TEXT,
			SignedMerkleRoot TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS Candidates (
			CandidateID VARCHAR(255) PRIMARY KEY,
			ElectionID VARCHAR(255) NOT NULL,
			CandidateName VARCHAR(255) NOT NULL,
			Party VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS Tokens (
			TokenID VARCHAR(255) PRIMARY KEY,
			TokenHash VARCHAR(255) UNIQUE NOT NULL,
			IssuedTime TIMESTAMP NOT NULL,
			ExpiryTime TIMESTAMP NOT NULL,
			Used BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS Votes (
			VoteID VARCHAR(255) PRIMARY KEY,
			ElectionID VARCHAR(255) NOT NULL,
			EncryptedBallot TEXT NOT NULL,
			Timestamp TIMESTAMP NOT NULL,
			MerkleLeafHash VARCHAR(255) UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS AuditLogs (
			LogID VARCHAR(255) PRIMARY KEY,
			Action TEXT NOT NULL,
			Timestamp TIMESTAMP NOT NULL,
			UserRole VARCHAR(50) NOT NULL,
			IPAddress VARCHAR(50) NOT NULL
		)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Migration failed: %v\nQuery: %s", err, query)
		}
	}
	log.Println("Migrations ran successfully.")
}

func seedData() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check Users count: %v", err)
	}

	if count > 0 {
		return
	}

	log.Println("Seeding default users...")
	users := []struct {
		ID   string
		Name string
		Pass string
		Role string
	}{
		{"voter1", "Alice Voter", "password123", "Voter"},
		{"voter2", "Bob Voter", "password123", "Voter"},
		{"officer1", "Charlie Election Officer", "password123", "ElectionOfficer"},
		{"admin1", "David Administrator", "password123", "Administrator"},
		{"counting1", "Emily Counting Authority", "password123", "CountingAuthority"},
	}

	for _, u := range users {
		hash, err := crypto.HashPassword(u.Pass)
		if err != nil {
			log.Fatalf("Failed to hash seed password: %v", err)
		}
		_, err = DB.Exec(
			"INSERT INTO Users (UserID, Name, PasswordHash, HasVoted, Role) VALUES ($1, $2, $3, $4, $5)",
			u.ID, u.Name, hash, false, u.Role,
		)
		// Fallback for sqlite which supports ? placeholder or $1. Let's make sure it handles both.
		// Wait, pg driver requires $1, $2. sqlite3 driver also supports $1, $2!
		// Yes, sqlite3 supports $1, $2 placeholders! This is perfect.
		if err != nil {
			// Try with ? if $1 fails, just in case. But $1 is standard in sqlite3 as well.
			_, err = DB.Exec(
				"INSERT INTO Users (UserID, Name, PasswordHash, HasVoted, Role) VALUES (?, ?, ?, ?, ?)",
				u.ID, u.Name, hash, false, u.Role,
			)
			if err != nil {
				log.Fatalf("Failed to insert seed user %s: %v", u.ID, err)
			}
		}
	}
	log.Println("Database seeded successfully.")
}

// LogEvent helper to write to AuditLogs table directly
func LogEvent(action, role, ip string) {
	id := fmt.Sprintf("log_%d", time.Now().UnixNano())
	query := "INSERT INTO AuditLogs (LogID, Action, Timestamp, UserRole, IPAddress) VALUES ($1, $2, $3, $4, $5)"
	_, err := DB.Exec(query, id, action, time.Now(), role, ip)
	if err != nil {
		// Try sqlite backup placeholder
		query = "INSERT INTO AuditLogs (LogID, Action, Timestamp, UserRole, IPAddress) VALUES (?, ?, ?, ?, ?)"
		_, err = DB.Exec(query, id, action, time.Now(), role, ip)
		if err != nil {
			log.Printf("Failed to log audit event: %v", err)
		}
	}
}
