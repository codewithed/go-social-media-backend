package main

import (
	_ "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	ListenAddr string
	Store      Storage
}

func NewApiServer(addr string, store Storage) (*ApiServer, error) {
	return &ApiServer{
		ListenAddr: addr,
		Store:      store,
	}, nil
}

func (s *ApiServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	http.ListenAndServe(s.ListenAddr, r)
}
