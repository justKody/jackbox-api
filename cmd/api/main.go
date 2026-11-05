package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/justKody/jackbox-api/internal/config"
	"github.com/justKody/jackbox-api/internal/gameserver"
	"github.com/justKody/jackbox-api/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	games, err := gameserver.Start()
	if err != nil {
		return err
	}
	httpServer := server.New(":"+config.Envs.Port, config.Envs.AssetDir, games.HandleWs)

	errCh := make(chan error, 1)
	go func() { errCh <- httpServer.Start() }()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
		return nil
	case err := <-errCh:
		// e.g. the port is taken
		return fmt.Errorf("http server on %s stopped: %w", httpServer.ListenAddr, err)
	}
}
