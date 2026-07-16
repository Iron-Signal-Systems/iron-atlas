package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/app"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
	)
	cfg, err := app.ConfigFromEnvironment()
	if err != nil {
		logger.Error("configuration failed", "error", err)
		os.Exit(1)
	}

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("application initialization failed", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	server := &http.Server{
		Addr:              cfg.ListenAddress,
		Handler:           application.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info(
			"iron atlas starting",
			"listen",
			cfg.ListenAddress,
			"authentication_mode",
			cfg.AuthenticationMode,
			"change_store",
			cfg.ChangeStore,
		)
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("iron atlas stopped")
}
