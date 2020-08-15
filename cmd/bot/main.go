package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"

	"free.blr/integrator/internal/bot"
	"free.blr/integrator/internal/repository"
)

var (
	EnvToken = os.Getenv("TOKEN")
	EnvDB    = os.Getenv("DB")
)

const (
	maxIdleCons = 10
	maxOpenCons = 10
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := setupDBConn(EnvDB)
	if err != nil {
		log.Panic(err)
	}

	tagRepo := repository.NewTag(conn)
	requestRepo := repository.NewRequest(conn)

	bot, err := bot.NewBot(EnvToken, tagRepo, requestRepo)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug(true)
	go func() {
		err := bot.Run(ctx, 0)
		if err != nil {
			log.Panic(err)
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	s := <-terminate
	log.Printf("service closed with signal: %s", s)
}

func setupDBConn(config string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", config)
	if err != nil {
		return nil, errors.Wrap(err, "could not open DB connections")
	}
	conn.SetMaxIdleConns(maxIdleCons)
	conn.SetMaxOpenConns(maxOpenCons)
	return conn, err
}
