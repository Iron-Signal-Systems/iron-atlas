package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/atlas/internal/authz"
	"github.com/Iron-Signal-Systems/atlas/internal/change"
	changepostgresql "github.com/Iron-Signal-Systems/atlas/internal/change/postgresql"
	database "github.com/Iron-Signal-Systems/atlas/internal/database/postgresql"
	"github.com/Iron-Signal-Systems/atlas/internal/health"
	"github.com/Iron-Signal-Systems/atlas/internal/httpui"
	"github.com/Iron-Signal-Systems/atlas/internal/modules"
)

const (
	ChangeStoreMemory     = "memory"
	ChangeStorePostgreSQL = "postgresql"
)

type Config struct {
	ListenAddress      string
	AuthenticationMode authentication.Mode
	ChangeStore        string
	Database           database.Config
	StartupTimeout     time.Duration
}

func ConfigFromEnvironment() (Config, error) {
	listen := strings.TrimSpace(os.Getenv("atlas_LISTEN"))
	if listen == "" {
		listen = "127.0.0.1:8080"
	}
	store := strings.ToLower(strings.TrimSpace(os.Getenv("atlas_CHANGE_STORE")))
	if store == "" {
		store = ChangeStoreMemory
	}
	if store != ChangeStoreMemory && store != ChangeStorePostgreSQL {
		return Config{}, fmt.Errorf("unsupported change store %q", store)
	}

	authenticationMode, err := authenticationModeFromEnvironment(store)
	if err != nil {
		return Config{}, err
	}

	maxConns, err := envInt32("atlas_DATABASE_MAX_CONNECTIONS", 8)
	if err != nil {
		return Config{}, err
	}
	minConns, err := envInt32("atlas_DATABASE_MIN_CONNECTIONS", 0)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ListenAddress:      listen,
		AuthenticationMode: authenticationMode,
		ChangeStore:        store,
		StartupTimeout:     10 * time.Second,
		Database: database.Config{
			URL:               strings.TrimSpace(os.Getenv("atlas_DATABASE_URL")),
			ApplicationName:   "atlas",
			MaxConnections:    maxConns,
			MinConnections:    minConns,
			ConnectTimeout:    5 * time.Second,
			MaxConnectionAge:  30 * time.Minute,
			MaxConnectionIdle: 5 * time.Minute,
			HealthInterval:    30 * time.Second,
			StatementTimeout:  15 * time.Second,
			LockTimeout:       5 * time.Second,
			IdleInTxTimeout:   15 * time.Second,
		},
	}
	if store == ChangeStorePostgreSQL && cfg.Database.URL == "" {
		return Config{}, errors.New(
			"atlas_DATABASE_URL is required when " +
				"atlas_CHANGE_STORE=postgresql",
		)
	}
	return cfg, nil
}

func authenticationModeFromEnvironment(
	store string,
) (authentication.Mode, error) {
	if legacy := strings.TrimSpace(
		os.Getenv("atlas_DEV_IDENTITY"),
	); legacy != "" {
		return "", errors.New(
			"atlas_DEV_IDENTITY is no longer supported; use " +
				"atlas_AUTHENTICATION_MODE=development or production",
		)
	}

	raw := strings.TrimSpace(
		os.Getenv("atlas_AUTHENTICATION_MODE"),
	)
	if raw == "" {
		if store == ChangeStoreMemory {
			return authentication.ModeDevelopment, nil
		}
		return authentication.ModeProduction, nil
	}
	mode, err := authentication.ParseMode(raw)
	if err != nil {
		return "", fmt.Errorf(
			"atlas_AUTHENTICATION_MODE: %w",
			err,
		)
	}
	return mode, nil
}

func envInt32(name string, fallback int32) (int32, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	return int32(value), nil
}

type Application struct {
	handler http.Handler
	close   func()
}

func New(cfg Config, logger *slog.Logger) (*Application, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	authenticationMiddleware, err := authentication.New(
		authentication.Options{Mode: cfg.AuthenticationMode},
	)
	if err != nil {
		return nil, fmt.Errorf("initialize authentication boundary: %w", err)
	}

	policy := authz.DefaultPolicy()
	var (
		changes   change.Service
		readiness health.Checker
		closeFn   = func() {}
	)

	switch cfg.ChangeStore {
	case ChangeStoreMemory:
		memory := change.NewMemoryService(policy)
		memory.Seed(change.Request{
			ID:                "CHG-2026-0001",
			Title:             "Establish first Cisco thirty-day collection profile",
			Risk:              change.RiskModerate,
			Status:            change.StatusPendingApproval,
			Requester:         "network-tech-01",
			RequiredApprovals: 1,
		})
		changes = memory
		readiness = health.Static{DependencyName: ChangeStoreMemory}
	case ChangeStorePostgreSQL:
		startupCtx, cancel := context.WithTimeout(
			context.Background(),
			cfg.StartupTimeout,
		)
		defer cancel()
		pool, err := database.Open(startupCtx, cfg.Database)
		if err != nil {
			return nil, fmt.Errorf(
				"initialize PostgreSQL runtime: %w",
				err,
			)
		}
		service, err := changepostgresql.New(pool)
		if err != nil {
			pool.Close()
			return nil, err
		}
		changes = service
		readiness = pool
		closeFn = pool.Close
	default:
		return nil, fmt.Errorf(
			"unsupported change store %q",
			cfg.ChangeStore,
		)
	}

	registry := modules.DefaultRegistry()
	server, err := httpui.New(httpui.Dependencies{
		Logger:         logger,
		Policy:         policy,
		Changes:        changes,
		Modules:        registry,
		Readiness:      readiness,
		Authentication: authenticationMiddleware,
	})
	if err != nil {
		closeFn()
		return nil, err
	}
	return &Application{
		handler: server.Handler(),
		close:   closeFn,
	}, nil
}

func (a *Application) Handler() http.Handler { return a.handler }

func (a *Application) Close() {
	if a != nil && a.close != nil {
		a.close()
	}
}
