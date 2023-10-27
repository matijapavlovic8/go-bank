package main

import (
	"go-bank-v2/internal/api"
	"go-bank-v2/internal/infrastructure/postgres"
	"log"
)

type app struct {
	config config
}

func newApp(config config) *app {
	return &app{
		config: config,
	}
}

func (a *app) Run() {
	store, err := postgres.NewPostgresStore(a.config.db)
	if err != nil {
		log.Fatal(err)
	}

	if err = store.Migrate(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(a.config.api, store)
	server.Run()

}
