package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"attendance-service/routes"
	"attendance-service/storage"
	"shared/middleware"
	"shared/pkgs/jwtmanager"
	"shared/pkgs/keys"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Connect to MongoDB

	// if err := keys.InitKeyPair(); err != nil {
	// 	log.Fatalf("failed to initiate key pair %v", err)
	// }

	if err := storage.InitMongo(); err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}
	fmt.Println("MongoDB connected")

	// Initialize Vault-backed JWT signer and cache public key
	if err := jwtmanager.InitVaultJWT(); err != nil {
		log.Printf("⚠️ JWT initialization failed (Vault missing): %v", err)
		// Try local generated keypair for dev if Vault not available
		if kerr := keys.GenerateLocalKeyPair(2048); kerr == nil {
			localPub := keys.GetPublicKey()
			if localPub != nil {
				jwtmanager.SetPublicKey(localPub)
				fmt.Println("⚠️ Using generated local public key fallback for JWT verification")
			} else {
				log.Fatalf("JWT initialization failed and local public key unavailable: %v", err)
			}
		} else {
			log.Fatalf("JWT initialization failed and local key generation failed: %v / %v", err, kerr)
		}
	} else {
		fmt.Println("✅ JWT signer initialized")
	}

	// Create main app router
	app := gin.Default()

	// Enable CORS globally
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type, Accept, Origin, X-Requested-With, X-CSRF-Token, X-Company-Code, Authorization, X-Forwarded-Proto"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Logging middleware
	app.Use(middleware.LogMiddleware())
	app.Use(gin.Recovery())

	// API group
	api := app.Group("/api/v1")

	// Internal JWT health endpoint (no auth, dev-only)
	api.GET("/internal/jwt", routes.JwtStatus)

	// Load routes
	routes.Routes(api)

	// Start server
	server := &http.Server{
		Addr:    ":8082",
		Handler: app,
	}

	fmt.Println("🚀 Server running on :8082")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
