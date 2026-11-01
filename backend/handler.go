package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/rishabhjha/rakshagrid/internal/api"
	"github.com/rishabhjha/rakshagrid/internal/auth"
	"github.com/rishabhjha/rakshagrid/internal/config"
	"github.com/rishabhjha/rakshagrid/internal/db"
	"github.com/rishabhjha/rakshagrid/internal/user"
)

var (
	vercelHandlerOnce sync.Once
	vercelHandler     http.Handler
	vercelHandlerErr  error
)

// Handler is the Vercel Go function entrypoint. Unlike main's long-running
// server, Vercel invokes this handler for each request.
func Handler(w http.ResponseWriter, r *http.Request) {
	vercelHandlerOnce.Do(initializeVercelHandler)
	if vercelHandlerErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "backend initialization failed"})
		log.Printf("Vercel handler initialization failed: %v", vercelHandlerErr)
		return
	}

	// Accommodate rewrites that preserve /api as well as those that strip it.
	if strings.HasPrefix(r.URL.Path, "/api/") {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
	}
	vercelHandler.ServeHTTP(w, r)
}

func initializeVercelHandler() {
	cfg := config.Load()
	authSvc := auth.NewService(cfg.JWTSecret, cfg.JWTExpiry)

	// Vercel services are stateless by default. Use in-memory demo accounts
	// unless database persistence is explicitly enabled and configured.
	var pg *sql.DB
	var userSvc user.UserService
	if os.Getenv("ENABLE_DATABASE") == "true" && cfg.PostgresDSN != "" {
		var err error
		pg, err = db.NewPostgres(cfg.PostgresDSN)
		if err != nil {
			log.Printf("Database unavailable; using demo accounts: %v", err)
			userSvc = user.NewDemoService(authSvc)
		} else if err = db.InitializeSchema(pg); err != nil {
			log.Printf("Database schema unavailable; using demo accounts: %v", err)
			_ = pg.Close()
			pg = nil
			userSvc = user.NewDemoService(authSvc)
		} else {
			userSvc = user.NewService(user.NewRepository(pg), authSvc)
		}
	} else {
		userSvc = user.NewDemoService(authSvc)
	}

	// Redis is not required for the demo login/API handlers and should not
	// delay serving requests during Vercel function initialization.
	vercelHandler = api.NewRouter(pg, nil, authSvc, userSvc, cfg.JWTExpiry, cfg.DispatchRadiusKM, cfg.SOSAckSLA)
}
