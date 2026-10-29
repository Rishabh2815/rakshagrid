package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rishabhjha/rakshagrid/internal/auth"
	"github.com/rishabhjha/rakshagrid/internal/models"
	"github.com/rishabhjha/rakshagrid/internal/user"
)

type Router struct {
	pg      *sql.DB
	rdb     interface{}
	authSvc *auth.Service
	userSvc user.UserService
}

func NewRouter(pg *sql.DB, rdb interface{}, authSvc *auth.Service, userSvc user.UserService, jwtExpiry time.Duration, radiusKM float64, ackSLA time.Duration) *chi.Mux {
	router := &Router{
		pg:      pg,
		rdb:     rdb,
		authSvc: authSvc,
		userSvc: userSvc,
	}

	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", router.handleLogin())
		r.Post("/signup", router.handleSignup())
	})

	// Protected routes (SOS, responders, incidents)
	r.Route("/sos", func(r chi.Router) {
		r.Post("/raise", router.handleRaiseSOS())
	})

	r.Route("/responders", func(r chi.Router) {
		r.Get("/nearby", router.handleNearbyResponders())
	})

	r.Get("/incidents", router.handleFetchIncidents())
	r.Get("/stream", router.handleWebSocketStream())

	return r
}

// handleLogin authenticates a user
func (router *Router) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "email and password required"})
			return
		}

		// Attempt login
		token, role, name, err := router.userSvc.Login(req.Email, req.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Token: token,
			Role:  role,
			Name:  name,
		})
	}
}

// handleSignup creates a new user account
func (router *Router) handleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" || req.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "email, password, and name required"})
			return
		}

		// Default role to 'user' if not specified
		if req.Role == "" {
			req.Role = "user"
		}

		// Attempt signup
		token, role, name, err := router.userSvc.Signup(req.Email, req.Password, req.Name, req.Role)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.SignupResponse{
			Token: token,
			Role:  role,
			Name:  name,
		})
	}
}

func (router *Router) handleRaiseSOS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"sos_received"}`))
	}
}

func (router *Router) handleNearbyResponders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"responders":[]}`))
	}
}

func (router *Router) handleFetchIncidents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"incidents":[]}`))
	}
}

func (router *Router) handleWebSocketStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("WebSocket endpoint not yet implemented"))
	}
}
