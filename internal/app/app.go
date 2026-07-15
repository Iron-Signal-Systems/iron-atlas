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

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	changepostgresql "github.com/Iron-Signal-Systems/iron-atlas/internal/change/postgresql"
	database "github.com/Iron-Signal-Systems/iron-atlas/internal/database/postgresql"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/health"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/httpui"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

const (
	ChangeStoreMemory     = "memory"
	ChangeStorePostgreSQL = "postgresql"
)

type Config struct {
	ListenAddress       string
	DevelopmentIdentity bool
	ChangeStore         string
	Database            database.Config
	StartupTimeout      time.Duration
}

func ConfigFromEnvironment() (Config, error) {
	listen := strings.TrimSpace(os.Getenv("IRON_ATLAS_LISTEN"))
	if listen == "" {
		listen = "127.0.0.1:8080"
	}
	store := strings.ToLower(strings.TrimSpace(os.Getenv("IRON_ATLAS_CHANGE_STORE")))
	if store == "" {
		store = ChangeStoreMemory
	}
	if store != ChangeStoreMemory && store != ChangeStorePostgreSQL {
		return Config{}, fmt.Errorf("unsupported change store %q", store)
	}

	devIdentity, err := developmentIdentityFromEnvironment(store)
	if err != nil {
		return Config{}, err
	}

	maxConns, err := envInt32("IRON_ATLAS_DATABASE_MAX_CONNECTIONS", 8)
	if err != nil {
		return Config{}, err
	}
	minConns, err := envInt32("IRON_ATLAS_DATABASE_MIN_CONNECTIONS", 0)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ListenAddress:       listen,
		DevelopmentIdentity: devIdentity,
		ChangeStore:         store,
		StartupTimeout:      10 * time.Second,
		Database: database.Config{
			URL:               strings.TrimSpace(os.Getenv("IRON_ATLAS_DATABASE_URL")),
			ApplicationName:   "iron-atlas",
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
		return Config{}, errors.New("IRON_ATLAS_DATABASE_URL is required when IRON_ATLAS_CHANGE_STORE=postgresql")
	}
	return cfg, nil
}

func developmentIdentityFromEnvironment(store string) (bool, error) {
	raw := strings.TrimSpace(os.Getenv("IRON_ATLAS_DEV_IDENTITY"))
	if raw == "" {
		// Preserve the simple Phase 0 memory-mode demonstration while making
		// persistent mode fail closed unless development headers are explicitly
		// enabled for a controlled test environment.
		return store == ChangeStoreMemory, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("IRON_ATLAS_DEV_IDENTITY must be a boolean: %w", err)
	}
	return value, nil
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
		startupCtx, cancel := context.WithTimeout(context.Background(), cfg.StartupTimeout)
		defer cancel()
		pool, err := database.Open(startupCtx, cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("initialize PostgreSQL runtime: %w", err)
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
		return nil, fmt.Errorf("unsupported change store %q", cfg.ChangeStore)
	}

	registry := modules.DefaultRegistry()
	server, err := httpui.New(httpui.Dependencies{
		Logger:              logger,
		Policy:              policy,
		Changes:             changes,
		Modules:             registry,
		Readiness:           readiness,
		DevelopmentIdentity: cfg.DevelopmentIdentity,
	})
	if err != nil {
		closeFn()
		return nil, err
	}
	return &Application{handler: server.Handler(), close: closeFn}, nil
}

func (a *Application) Handler() http.Handler { return a.handler }

func (a *Application) Close() {
	if a != nil && a.close != nil {
		a.close()
	}
}
