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

	// CORS middleware - permissive for Vercel services (same origin)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check - no prefix required
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes - served under /api by Vercel rewrite
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handleLogin(authSvc))
	})

	r.Route("/sos", func(r chi.Router) {
		r.Post("/raise", handleRaiseSOS(pg, rdb, authSvc))
	})

	r.Route("/responders", func(r chi.Router) {
		r.Get("/nearby", handleNearbyResponders(pg, rdb))
	})

	r.Get("/incidents", handleFetchIncidents(pg))

	// WebSocket stream for real-time updates
	r.Get("/stream", handleWebSocketStream(pg, rdb, authSvc))

	return r
}

func handleLogin(authSvc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement login logic
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"dummy-token","role":"user"}`))
	}
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

func handleFetchIncidents(pg *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement incident fetching
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"incidents":[]}`))
	}
}

func handleWebSocketStream(pg *sql.DB, rdb *redis.Client, authSvc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement WebSocket stream for real-time updates
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("WebSocket endpoint not yet implemented"))
	}
}
