package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"thugcorp.io/final_bot/auth"
	"thugcorp.io/final_bot/database/postgres"
	"thugcorp.io/final_bot/logger"
)

func initialiseLogger() *logger.Logger{
	logInstance, err := logger.NewLogger("robo.log")
	if err != nil {
		log.Fatal("Failed to initialise logger")
	}
	return logInstance
}

func initialiseOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:	 os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"account:write", "trading"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://app.alpaca.markets/oauth/authorize",
			TokenURL: "https://api.alpaca.markets/oauth/token",
		},
	}
}


func main() {
	// Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialise logger
	logInstance := initialiseLogger()
	defer logInstance.Close()

	// Initialise Postgres DB
	pgDB, err := postgres.NewPostgresDB(logInstance)
	if err != nil {
		logInstance.Error("Failed to initialise Postgres DB", err)
		return
	}

	// Initialise OAuth2 config
	oauthConfig := initialiseOAuthConfig()

	// Register repositories
	authRepo := auth.NewAuthRepository(pgDB.GetDB(), logInstance)
	// Register services
	authService := auth.NewAuthService(oauthConfig, logInstance)
	// Register handlers
	authHandler := auth.NewAuthHandler( logInstance, authRepo, authService)
	
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	// Register routes
	authHandler.RegisterRoutes(*api)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	log.Println("Server stopped")

}

