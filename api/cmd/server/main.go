// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/luungockhang/mahjong-tracker-skeleton/internal/server" // 👈 replace "yourmodule" with your module path from go.mod
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// Connect to DB
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Router setup
	r := chi.NewRouter()

	// Register routes via Server struct
	srv := server.NewServer(pool)
	srv.RegisterRoutes(r)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run server
	go func() {
		log.Println("Starting server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed:", err)
	}

	log.Println("Server gracefully stopped")
}
