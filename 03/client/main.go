package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const serverHost = "http://server:8080"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	clientStopped := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		startServerPolling(ctx)
		log.Info().Msg("end server pollling")
		clientStopped <- struct{}{}
	}(ctx)

	<-shutdown
	cancel()
	<-clientStopped
	log.Info().Msg("bye")
}

func startServerPolling(ctx context.Context) error {
	client := resty.New()

	for {
		if ctx.Err() != nil {
			return nil
		}

		time.Sleep(time.Second)
		createFile(client)
	}
}

func createFile(client *resty.Client) {
	url := serverHost + "/create"
	resp, err := client.R().Post(url)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send request")
	}

	if resp.IsError() {
		log.Error().Str("response", resp.String()).Int("code", resp.StatusCode()).Msg("file not created")
		return
	}

	log.Info().Msg("file created successfully")
}
