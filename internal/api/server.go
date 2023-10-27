package api

import (
	. "go-bank-v2/internal/infrastructure/postgres"
)

type Server struct {
	listenAddr string
	store      Store
}

func NewServer(config RestApiConfig, store Store) *Server {
	return &Server{
		listenAddr: ":" + config.Port,
		store:      store,
	}
}
