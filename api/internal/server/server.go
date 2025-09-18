package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luungockhang/mahjong-tracker-skeleton/internal/player"
)

type Server struct {
	DB *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{DB: db}
}

func (s *Server) RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/player", player.Route(s.DB))
			// auth routes later
		})
	})
}
