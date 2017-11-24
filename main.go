package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	router *chi.Mux
}

func New() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
}

func (s *Server) Middleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.CloseNotify)
	s.router.Use(middleware.Timeout(time.Second * 60))
}

func (s *Server) Router() {

}

func main() {
	s := New()
	s.Init()
	s.Middleware()
	s.Router()
	log.Println("Starting app")
	http.ListenAndServe(":3000", s.router)
}
