package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"free.blr/integrator/internal/bot"
)

var (
	EnvToken = os.Getenv("TOKEN")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := bot.NewBot(EnvToken)
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
