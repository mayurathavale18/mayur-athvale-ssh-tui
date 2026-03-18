package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/analytics"
	"mayur-athavale-tui/internal/config"
	"mayur-athavale-tui/internal/server"
)

func main() {
	cfg := config.Load()

	portfolio, err := content.LoadPortfolio()
	if err != nil {
		log.Fatal("Failed to load portfolio content", "error", err)
	}

	// Initialize analytics
	store, err := analytics.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize analytics", "error", err)
	}
	defer store.Close()

	tracker := analytics.NewTracker(store)

	srv, err := server.New(cfg, portfolio, tracker)
	if err != nil {
		log.Fatal("Failed to create server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Server error", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Shutdown error", "error", err)
	}
}
