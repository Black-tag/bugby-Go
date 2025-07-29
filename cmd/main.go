// @title Bugby API
// @version 1.0
// @description A bug tracking API written in Go with JWT, PostgreSQL and RBAC.
// @termsOfService http://swagger.io/terms/

// @contact.name Anand Unni
// @contact.url https://github.com/Black-tag
// @contact.email your-email@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
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
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpswagger "github.com/swaggo/http-swagger"
	_ "github.com/blacktag/bugby-Go/docs"
	
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
	enforcer, err := SetupCasbin()
	if err != nil {
		log.Fatal("failed to setup casbin: %w", err)
	}



	
	
	

	 
	authMiddleware := middleware.Authenticate1(cfg.SECRET, cfg.DB)
	authMiddleware2 := middleware.RevokeTokenAthenticate( cfg.DB)


	mux := http.NewServeMux()

	protected := authMiddleware(middleware.Authorization(enforcer)(http.HandlerFunc(cfg.DeleteBugByIDHandler)))
	mux.Handle("POST /api/bugs", authMiddleware(http.HandlerFunc(cfg.CreateBugHandler)))
	mux.Handle("DELETE /api/bugs/{bugid}", protected)
	mux.Handle("POST /api/bugs/{bugid}", authMiddleware(http.HandlerFunc(cfg.UpadteBugHandler)))
	mux.HandleFunc("GET /api/bugs/{bugid}", cfg.GetBugByIDHandler)
	mux.HandleFunc("GET /api/bugs", cfg.GetBugsHandler)
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshTokenHandler)
	mux.Handle("POST /api/revoke", authMiddleware2(http.HandlerFunc(cfg.RevokeTokenHandler)))
	mux.Handle("PUT /api/users", authMiddleware(http.HandlerFunc(cfg.UpdateCredentialsHandler)))
	mux.HandleFunc("/swagger/", httpswagger.WrapHandler)





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


func SetupCasbin() (*casbin.Enforcer, error) {
		m, err := model.NewModelFromFile("rbac_model.conf")
		if err != nil {
			
			return nil, fmt.Errorf("cannot load model for enforcer: %w", err)
		}	

		a := fileadapter.NewAdapter("rbac_policy.csv")

		enforcer, err := casbin.NewEnforcer(m, a)
		if err != nil {
			
			return nil, fmt.Errorf("cannot create enforcer: %w", err)
		} 
		err = enforcer.LoadPolicy()
		if err != nil {
			
			return nil, fmt.Errorf("cannot load policy: %w", err)
		}
		return enforcer, nil
	}