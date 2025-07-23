package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blacktag/bugby-Go/internal/api"
	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)





func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cfg := api.APIConfig{
		DB: dbQueries,
	}



	
	mux := http.NewServeMux()
	
	
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/bugs", cfg.CreateBugHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)


	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fmt.Println("🌐 starting the server on: http://localhost:8080...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}





}