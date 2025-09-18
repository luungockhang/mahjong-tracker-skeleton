package players

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
