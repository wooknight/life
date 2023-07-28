package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wooknight/life/config"
	"github.com/wooknight/life/internal/database/repository"
	"github.com/wooknight/life/internal/domain/entity"
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
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DB)
	db, err := setupDB(psqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	defer db.Close()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	httplog.Configure(httplog.Options{Concise: true, TimeFieldFormat: time.RFC3339})
	router := chi.NewRouter()
	router.Use(httplog.RequestLogger(log.Logger))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Recoverer)
	//new db entity
	thgt := repository.NewThoughtRepository(db)
	idea := entity.Thought{
		Thought:            "What a day",
		ThoughtDescription: "a brnd new duy",
		ThoughtTags:        []string{"wango", "pongo", "bongo"},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	id, err := thgt.Create(context.Background(), &idea)
	if err != nil {
		log.Fatal().Msgf("%s", err.Error())
	}
	dat, err := thgt.Get(context.Background(), id)
	if err != nil {
		log.Printf("Did the damned create succeed - %d\n", id)
		log.Fatal().Msgf("%s", err.Error())
	}
	log.Info().Msgf("%+v", dat)

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
