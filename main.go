package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/joho/godotenv"
	"github.com/lain-the-coder/ea-qms-backend/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db       *database.Queries
	platform string
	secret   string
	params   *argon2id.Params
	rawDB    *sql.DB
}

func (cfg *apiConfig) WelcomeHome(w http.ResponseWriter, r *http.Request) {
	type WelcomeRequest struct {
		Message string `json:"message"`
	}
	type WelcomeResponse struct {
		Company string `json:"company"`
		Message string `json:"message"`
	}
	reqBody := WelcomeRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		// delegating error structuring to helper function
		respondWithError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	reqBody.Message = strings.TrimSpace(reqBody.Message)
	if reqBody.Message == "" {
		log.Printf("Message is blank")
		respondWithError(w, "Message cannot be blank", http.StatusBadRequest)
		return
	}
	resBody := WelcomeResponse{
		Company: "EA QMS",
		Message: "Welcome! I hope you enjoy this system!",
	}
	respondWithJSON(w, http.StatusOK, resBody)
}

func main() {
	mux := http.NewServeMux()

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// load config struct with env variables
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("JWT_SECRET")
	argonParams := loadArgon2idParams()

	// db setup
	rawDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database initialization failed (check driver registration or URL format): %v", err)
	}
	err = rawDB.Ping()
	if err != nil {
		log.Fatalf("Database connection failed (check network, credentials, or server status): %v", err)
	}

	db := database.New(rawDB)

	cfg := &apiConfig{
		db:       db,
		platform: platform,
		secret:   secret,
		params:   argonParams,
		rawDB:    rawDB,
	}

	// routes
	mux.HandleFunc("POST /", cfg.WelcomeHome)
	server := &http.Server{
		Addr:    ":1304",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
