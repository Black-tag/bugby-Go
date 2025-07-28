package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blacktag/bugby-Go/internal/api"
	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "github.com/ydb-platform/ydb-go-sdk/v3/ratelimiter"
)





func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	secret := os.Getenv("SECRET")

	cfg := api.APIConfig{
		DB: dbQueries,
		SECRET: secret,
	}



	
	mux := http.NewServeMux()

	authMiddleware := middleware.Authenticate1(cfg.SECRET)
	authMiddleware2 := middleware.RevokeTokenAthenticate( cfg.DB)
	mux.Handle("POST /api/bugs", authMiddleware(http.HandlerFunc(cfg.CreateBugHandler)))



	mux.HandleFunc("DELETE /api/bugs/{bugid}", cfg.DeleteBugByIDHandler)
	mux.HandleFunc("POST /api/bugs/{bugid}", cfg.UpadteBugHandler)
	mux.HandleFunc("GET /api/bugs/{bugid}", cfg.GetBugByIDHandler)
	mux.HandleFunc("GET /api/bugs", cfg.GetBugsHandler)
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshTokenHandler)
	mux.Handle("POST /api/revoke", authMiddleware2(http.HandlerFunc(cfg.RevokeTokenHandler)))
	mux.Handle("PUT /api/users", authMiddleware(http.HandlerFunc(cfg.UpdateCredentialsHandler)))
	





	ratelimiter := middleware.NewRateLimiter(5,10,time.Minute)
	muxWithLimiter := ratelimiter.Limit(mux)
	server := &http.Server{
		Addr: "localhost:8080",
		Handler: muxWithLimiter,
	}

	fmt.Println("🌐 starting the server on: http://localhost:8080...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}





}