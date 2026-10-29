package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/rishabhjha/rakshagrid/internal/auth"
)

func NewRouter(pg *sql.DB, rdb *redis.Client, authSvc *auth.Service, jwtExpiry time.Duration, radiusKM float64, ackSLA time.Duration) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware - adjust in production
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Basic health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// TODO: Add your SOS, dispatch, responder, audit endpoints here
		// Example:
		r.Post("/sos/raise", handleRaiseSOS(pg, rdb, authSvc))
		r.Get("/responders/nearby", handleNearbyResponders(pg, rdb))
	})

	return r
}

func handleRaiseSOS(pg *sql.DB, rdb *redis.Client, authSvc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement SOS raising logic
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"sos_received"}`))
	}
}

func handleNearbyResponders(pg *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement geospatial responder search
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"responders":[]}`))
	}
}
