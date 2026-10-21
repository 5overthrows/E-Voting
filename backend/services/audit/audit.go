package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"evoting/services/auth"
	"evoting/shared/db"
	"evoting/shared/models"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the audit service routes on the given router
func RegisterRoutes(r *gin.Engine) {
	protected := r.Group("/api/audit")
	protected.Use(auth.AuthMiddleware("Administrator"))
	{
		protected.GET("/logs", getAuditLogs)
		protected.POST("/backup", backupDatabase)
	}
}

// StartAuditServer starts the audit service on the specified address
func StartAuditServer(addr string) {
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

	log.Printf("Audit Logging Service running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Audit server failed: %v", err)
	}
}

func getAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	query := "SELECT LogID, Action, Timestamp, UserRole, IPAddress FROM AuditLogs ORDER BY Timestamp DESC"
	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var l models.AuditLog
		err := rows.Scan(&l.LogID, &l.Action, &l.Timestamp, &l.UserRole, &l.IPAddress)
		if err == nil {
			logs = append(logs, l)
		}
	}

	c.JSON(http.StatusOK, logs)
}

func backupDatabase(c *gin.Context) {
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backups directory"})
		return
	}

	timestamp := time.Now().Format("20060102_150405")

	if db.DBType == "sqlite" {
		// Perform a simple file copy for SQLite
		srcFile := os.Getenv("SQLITE_PATH")
		if srcFile == "" {
			srcFile = "evoting.db"
		}
		destFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.db", timestamp))

		in, err := os.Open(srcFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open source DB: %v", err)})
			return
		}
		defer in.Close()

		out, err := os.Create(destFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create backup file: %v", err)})
			return
		}
		defer out.Close()

		if _, err = io.Copy(out, in); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to copy DB: %v", err)})
			return
		}

		db.LogEvent("Database backup created: "+destFile, "Administrator", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"message": "SQLite database backup successful",
			"file":    destFile,
		})
		return
	}

	// For PostgreSQL, serialize tables into a JSON structure
	backupData := make(map[string]interface{})

	// Tables to back up
	tables := []string{"Users", "Elections", "Candidates", "Tokens", "Votes", "AuditLogs"}
	for _, table := range tables {
		rows, err := db.DB.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read table %s: %v", table, err)})
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		var tableRows []map[string]interface{}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan row: %v", err)})
				return
			}

			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columns[i]
				b, ok := val.([]byte)
				if ok {
					m[colName] = string(b)
				} else {
					m[colName] = val
				}
			}
			tableRows = append(tableRows, m)
		}
		backupData[table] = tableRows
	}

	destFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.json", timestamp))
	file, err := os.Create(destFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup JSON file"})
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backupData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode backup data to JSON"})
		return
	}

	db.LogEvent("Database backup created (PostgreSQL JSON export): "+destFile, "Administrator", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"message": "PostgreSQL database backup successful (JSON export format)",
		"file":    destFile,
	})
}
