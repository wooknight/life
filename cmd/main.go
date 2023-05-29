package main

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/wooknight/life/config"
)

type server struct {
	Router chi.Router
}

func main() {
	config, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load configurations")
	}
	//start DB
	db, err := setupDB(config.DBUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
}

func setupDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
