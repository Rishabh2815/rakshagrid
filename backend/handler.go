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

	var pg *sql.DB
	var userSvc user.UserService
	if cfg.PostgresDSN != "" {
		var err error
		pg, err = db.NewPostgres(cfg.PostgresDSN)
		if err != nil {
			vercelHandlerErr = err
			return
		}
		if err = db.InitializeSchema(pg); err != nil {
			_ = pg.Close()
			vercelHandlerErr = err
			return
		}
		userSvc = user.NewService(user.NewRepository(pg), authSvc)
	} else {
		userSvc = user.NewDemoService(authSvc)
	}

	// Redis-backed routes are not needed for demo login; only initialize Redis
	// if a deployment explicitly supplies its address.
	var rdb interface{}
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		client, err := db.NewRedis(addr, cfg.RedisPassword)
		if err != nil {
			log.Printf("Redis unavailable; continuing without Redis: %v", err)
		} else {
			rdb = client
		}
	}

	vercelHandler = api.NewRouter(pg, rdb, authSvc, userSvc, cfg.JWTExpiry, cfg.DispatchRadiusKM, cfg.SOSAckSLA)
}
