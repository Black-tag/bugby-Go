package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	
)

type apiConfig struct {
	db *database.Queries

} 



func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgre", dbURL)
	dbQueries := database.New(db)

	cfg := apiapiConfig{
		db: dbQueries,
	}



	mainFileServer := http.FileServer((http.Dir(".")))
	mux := http.NewServeMux()
	mux.Handle("/", mainFileServer)

	
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