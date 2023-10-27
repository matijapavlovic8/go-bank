package main

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	. "go-bank-v2/internal/api"
	. "go-bank-v2/internal/infrastructure/postgres"
	"log"
)

type config struct {
	db  DbConfig
	api RestApiConfig
}

func loadConfig() config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var config config
	var dbConfig DbConfig
	var apiConfig RestApiConfig

	if err := env.Parse(&dbConfig); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	if err := env.Parse(&apiConfig); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	config.db = dbConfig
	config.api = apiConfig

	return config
}
