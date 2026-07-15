package app

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/httpui"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

type Config struct {
	ListenAddress       string
	DevelopmentIdentity bool
}

func ConfigFromEnvironment() Config {
	listen := strings.TrimSpace(os.Getenv("IRON_ATLAS_LISTEN"))
	if listen == "" {
		listen = "127.0.0.1:8080"
	}
	devIdentity := strings.EqualFold(os.Getenv("IRON_ATLAS_DEV_IDENTITY"), "true") || os.Getenv("IRON_ATLAS_DEV_IDENTITY") == ""
	return Config{ListenAddress: listen, DevelopmentIdentity: devIdentity}
}

type Application struct {
	handler http.Handler
}

func New(cfg Config, logger *slog.Logger) (*Application, error) {
	policy := authz.DefaultPolicy()
	changes := change.NewMemoryService(policy)
	changes.Seed(change.Request{
		ID:                "CHG-2026-0001",
		Title:             "Establish first Cisco thirty-day collection profile",
		Risk:              change.RiskModerate,
		Status:            change.StatusPendingApproval,
		Requester:         "network-tech-01",
		RequiredApprovals: 1,
	})

	registry := modules.DefaultRegistry()
	server, err := httpui.New(httpui.Dependencies{
		Logger:              logger,
		Policy:              policy,
		Changes:             changes,
		Modules:             registry,
		DevelopmentIdentity: cfg.DevelopmentIdentity,
	})
	if err != nil {
		return nil, err
	}
	return &Application{handler: server.Handler()}, nil
}

func (a *Application) Handler() http.Handler { return a.handler }
