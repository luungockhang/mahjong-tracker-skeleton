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
)

func main() {
	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Mount("/players", playerRoutes(pool))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

// internal/players/handlers.go
package players

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Player struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	Rating      int       `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

func playerRoutes(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Get("/", listPlayers(pool))
	r.Post("/", createPlayer(pool))
	return r
}

func listPlayers(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := pool.Query(context.Background(), "SELECT id, display_name, rating, created_at FROM players")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var players []Player
		for rows.Next() {
			var p Player
			if err := rows.Scan(&p.ID, &p.DisplayName, &p.Rating, &p.CreatedAt); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			players = append(players, p)
		}
		json.NewEncoder(w).Encode(players)
	}
}

func createPlayer(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var id int
		err := pool.QueryRow(context.Background(),
			"INSERT INTO players(display_name, rating) VALUES($1, 1500) RETURNING id",
			input.DisplayName,
		).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": id})
	}
}

// migrations/0001_init.up.sql
