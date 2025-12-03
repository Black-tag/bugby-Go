// @title Bugby API
// @version 1.0
// @description A bug tracking API written in Go with JWT, PostgreSQL and RBAC.
// @termsOfService http://swagger.io/terms/

// @contact.name Anand Unni
// @contact.url https://github.com/Black-tag
// @contact.email unnianandunni007@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blacktag/bugby-Go/internal/api"
	"github.com/blacktag/bugby-Go/internal/caching"
	"github.com/blacktag/bugby-Go/internal/database"
	_ "github.com/blacktag/bugby-Go/internal/docs"
	"github.com/blacktag/bugby-Go/internal/metrics"
	"github.com/blacktag/bugby-Go/internal/middleware"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	httpswagger "github.com/swaggo/http-swagger"
	// "github.com/ydb-platform/ydb-go-sdk/v3/ratelimiter"
)

func main() {

	godotenv.Load()

	// env := os.Getenv("APP_ENV")
	// if env == "production" {
	// 	godotenv.Load(".env.production")
	// } else {
	// godotenv.Load(".env.development")
	// }
	if os.Getenv("APP_ENV") != "production" {
		godotenv.Load(".env.development")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	secret := os.Getenv("SECRET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = caching.InitCache()
	if err != nil {
		log.Fatal("failed to init cache: ", err)
	}

	cfg := api.APIConfig{
		DB:     dbQueries,
		SECRET: secret,
	}
	enforcer, err := SetupCasbin()
	if err != nil {
		log.Fatal("failed to setup casbin: ", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	reg := prometheus.DefaultRegisterer
	m := metrics.NewMetrics(reg)

	loggingMiddleware := middleware.MetricsMiddleware(m)

	cachingMiddleware := middleware.CachingMiddleware(5 * time.Minute)

	authMiddleware := middleware.Authenticate(cfg.SECRET, cfg.DB)
	authMiddleware2 := middleware.RevokeTokenAthenticate(cfg.DB)

	mux := http.NewServeMux()
	muxWithMetrics := loggingMiddleware(mux)

	protected := authMiddleware(middleware.Authorization(enforcer)(http.HandlerFunc(cfg.DeleteBugByIDHandler)))
	mux.Handle("POST /api/bugs", authMiddleware(http.HandlerFunc(cfg.CreateBugHandler)))
	mux.Handle("DELETE /api/bugs/{bugid}", protected)
	mux.Handle("POST /api/bugs/{bugid}", authMiddleware(http.HandlerFunc(cfg.UpdateBugHandler)))
	mux.Handle("GET /api/bugs/{bugid}", cachingMiddleware(http.HandlerFunc(cfg.GetBugByIDHandler)))
	mux.Handle("GET /api/bugs", cachingMiddleware(http.HandlerFunc(cfg.GetBugsHandler)))
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshTokenHandler)
	mux.Handle("POST /api/revoke", authMiddleware2(http.HandlerFunc(cfg.RevokeTokenHandler)))
	mux.Handle("PUT /api/users", authMiddleware(http.HandlerFunc(cfg.UpdateCredentialsHandler)))
	mux.HandleFunc("/swagger/", httpswagger.WrapHandler)
	mux.Handle("GET /api/users", cachingMiddleware(http.HandlerFunc(cfg.GetUsersHandler)))
	mux.Handle("GET /api/users/me/bugs", authMiddleware(cachingMiddleware(http.HandlerFunc(cfg.GetUserSpecificBugs))))
	mux.Handle("/metrics/", metrics.MetricsHandler())

	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("TEST LOG MESSAGE", "key", "value")
		w.Write([]byte("Check console logs"))
	})

	ratelimiter := middleware.NewRateLimiter(5, 10, time.Minute)
	muxWithLimiter := ratelimiter.Limit(muxWithMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: muxWithLimiter,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	
	go func() {
		if err = server.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}

	}()
	logger.Info("server started", "port", 8080)
	fmt.Println("🌐 starting the server on: http://localhost:8080...")

	<-done
	logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	logger.Info("Server Exited Succesfully")

}

func SetupCasbin() (*casbin.Enforcer, error) {
	m, err := model.NewModelFromFile("rbac_model.conf")
	if err != nil {

		return nil, fmt.Errorf("cannot load model for enforcer: %v", err)
	}

	a := fileadapter.NewAdapter("rbac_policy.csv")

	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {

		return nil, fmt.Errorf("cannot create enforcer: %v", err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {

		return nil, fmt.Errorf("cannot load policy: %v", err)
	}
	return enforcer, nil
}