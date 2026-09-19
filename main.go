// Command server is the single entrypoint for the RakshaGrid modular
// monolith. It wires config -> DB connections -> the HTTP router, and
// runs a background sweeper that auto-escalates incidents that breach
// their ACK SLA. Splitting any of internal/* into an independently
// deployed service later means moving its directory and adding a new
// main.go that wires just that service — the packages themselves don't
// change (see docs/system-design "Growth path").
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rishabhjha/rakshagrid/internal/api"
	"github.com/rishabhjha/rakshagrid/internal/auth"
	"github.com/rishabhjha/rakshagrid/internal/config"
	"github.com/rishabhjha/rakshagrid/internal/db"
)

func main() {
	cfg := config.Load()

	pg, err := db.NewPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer pg.Close()

	rdb, err := db.NewRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer rdb.Close()

	authSvc := auth.NewService(cfg.JWTSecret, cfg.JWTExpiry)
	router := api.NewRouter(pg, rdb, authSvc, cfg.JWTExpiry, cfg.DispatchRadiusKM, cfg.SOSAckSLA)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second, // longer than default to accommodate the WS stream endpoint
	}

	sweeperCtx, cancelSweeper := context.WithCancel(context.Background())
	go runEscalationSweeper(sweeperCtx, pg, cfg.SOSAckSLA)

	go func() {
		log.Printf("rakshagrid server listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	cancelSweeper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// runEscalationSweeper polls for incidents stuck in "raised" or "matched"
// status past the ACK SLA and marks them escalated so an admin gets
// pulled in. A polling loop is the pragmatic v1 choice for a single-
// instance deployment; once dispatch is split into its own service (see
// growth-path table in the design doc), this becomes a proper delayed-job
// queue instead of a poll.
func runEscalationSweeper(ctx context.Context, pg *sql.DB, ackSLA time.Duration) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().Add(-ackSLA)
			rows, err := pg.QueryContext(ctx, `
				SELECT id FROM incidents
				WHERE status IN ('raised', 'matched') AND created_at < $1
			`, cutoff)
			if err != nil {
				log.Printf("escalation sweeper: query failed: %v", err)
				continue
			}

			var staleIDs []string
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err == nil {
					staleIDs = append(staleIDs, id)
				}
			}
			rows.Close()

			for _, id := range staleIDs {
				_, err := pg.ExecContext(ctx, `
					UPDATE incidents SET status = 'escalated' WHERE id = $1
				`, id)
				if err != nil {
					log.Printf("escalation sweeper: failed to escalate %s: %v", id, err)
					continue
				}
				_, err = pg.ExecContext(ctx, `
					INSERT INTO incident_events (incident_id, event_type, actor_id, detail)
					VALUES ($1, 'escalated', NULL, 'auto-escalated: no ACK within SLA')
				`, id)
				if err != nil {
					log.Printf("escalation sweeper: failed to log event for %s: %v", id, err)
				}
				log.Printf("escalation sweeper: auto-escalated incident %s", id)
			}
		}
	}
}
