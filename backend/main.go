package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"evoting/services/audit"
	"evoting/services/auth"
	"evoting/services/counting"
	"evoting/services/demo"
	"evoting/services/token"
	"evoting/services/voting"
	"evoting/shared/db"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	serviceFlag := flag.String("service", "", "Specify service to run (auth, token, voting, counting, audit, microservices, or empty for unified gateway)")
	portFlag := flag.String("port", "", "Port for unified gateway (default 8080)")
	flag.Parse()

	// Connect to database and run migrations
	db.Connect()

	serviceName := *serviceFlag
	if serviceName == "" {
		serviceName = os.Getenv("SERVICE_NAME")
	}

	switch serviceName {
	case "auth":
		port := os.Getenv("PORT_AUTH")
		if port == "" {
			port = ":8081"
		}
		auth.StartAuthServer(port)

	case "token":
		port := os.Getenv("PORT_TOKEN")
		if port == "" {
			port = ":8082"
		}
		token.StartTokenServer(port)

	case "voting":
		port := os.Getenv("PORT_VOTING")
		if port == "" {
			port = ":8083"
		}
		voting.StartVotingServer(port)

	case "counting":
		port := os.Getenv("PORT_COUNTING")
		if port == "" {
			port = ":8084"
		}
		counting.StartCountingServer(port)

	case "audit":
		port := os.Getenv("PORT_AUDIT")
		if port == "" {
			port = ":8085"
		}
		audit.StartAuditServer(port)

	case "microservices":
		log.Println("Starting AegisVote Microservices Architecture concurrently...")
		go auth.StartAuthServer(":8081")
		go token.StartTokenServer(":8082")
		go voting.StartVotingServer(":8083")
		go counting.StartCountingServer(":8084")
		go audit.StartAuditServer(":8085")

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down microservices...")

	default:
		// Default / Dev / Faculty Review Mode: Unified Gateway
		// All 5 microservice route trees are mounted onto a single high-performance Gin router.
		// Solves Windows socket collisions, port blocking, and enables seamless 1-click faculty presentation.
		port := *portFlag
		if port == "" {
			port = os.Getenv("PORT")
		}
		if port == "" {
			port = "8080"
		}
		addr := ":" + port
		os.Setenv("PORT", port)
		os.Setenv("UNIFIED_MODE", "true")

		r := gin.Default()
		r.Use(corsMiddleware())

		// System Health & Welcome
		r.GET("/api/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"system":  "AegisVote Privacy-Preserving E-Voting System",
				"version": "1.0-IEEE",
				"mode":    "Unified Faculty Presentation & Testing Gateway",
			})
		})

		// Mount all microservice routes
		auth.RegisterRoutes(r)
		token.RegisterRoutes(r)
		voting.RegisterRoutes(r)
		counting.RegisterRoutes(r)
		audit.RegisterRoutes(r)
		demo.RegisterRoutes(r)

		fmt.Println("==========================================================================")
		fmt.Println("  AEGISVOTE: Privacy-Preserving Electronic Voting System (IEEE SRS v1.0)")
		fmt.Printf("  Unified Gateway & Presentation Engine Listening on http://localhost:%s\n", port)
		fmt.Println("  Ready for Faculty Review, Live Walkthrough, and Jira / QA Automation!")
		fmt.Println("==========================================================================")

		if err := r.Run(addr); err != nil {
			log.Fatalf("Unified gateway failed: %v", err)
		}
	}
}
