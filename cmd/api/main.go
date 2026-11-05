package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justKody/jackbox-api/internal/config"
	"github.com/justKody/jackbox-api/internal/gameserver"
	"github.com/justKody/jackbox-api/internal/server"
)

// how long open requests and actors get to finish once the server is stopping
const shutdownTimeout = 5 * time.Second

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
	case err := <-errCh:
		// e.g. the port is taken
		games.Stop(context.Background())
		return fmt.Errorf("http server on %s stopped: %w", httpServer.ListenAddr, err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// stop taking clients first. Shutdown leaves websockets alone, stopping
	// the game server is what closes them
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println("http server shutdown:", err)
	}
	games.Stop(shutdownCtx)
	log.Println("server stopped")
	return nil
}
