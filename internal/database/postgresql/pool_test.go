package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func validConfig() Config {
	return Config{
		URL:               "postgres://example",
		ApplicationName:   "iron-atlas-test",
		MaxConnections:    4,
		MinConnections:    1,
		ConnectTimeout:    time.Second,
		MaxConnectionAge:  time.Minute,
		MaxConnectionIdle: time.Minute,
		HealthInterval:    time.Second,
		StatementTimeout:  time.Second,
		LockTimeout:       time.Second,
		IdleInTxTimeout:   time.Second,
	}
}

func TestConfigValidation(t *testing.T) {
	valid := validConfig()
	if err := valid.validate(); err != nil {
		t.Fatalf("valid configuration rejected: %v", err)
	}

	cases := []Config{
		{},
		func() Config { c := validConfig(); c.MaxConnections = 0; return c }(),
		func() Config { c := validConfig(); c.MinConnections = 5; return c }(),
		func() Config { c := validConfig(); c.StatementTimeout = 0; return c }(),
		func() Config {
			c := validConfig()
			c.URL = "postgres://example?options=-c%20atlas.actor_id%3Devil"
			return c
		}(),
	}
	for _, cfg := range cases {
		if err := cfg.validate(); err == nil {
			t.Fatalf("invalid configuration accepted: %+v", cfg)
		}
	}
}

func TestActorValidation(t *testing.T) {
	if err := validateActorID("network-admin@example.test"); err != nil {
		t.Fatal(err)
	}
	for _, actor := range []string{"", "actor\nother"} {
		if err := validateActorID(actor); err == nil {
			t.Fatalf("invalid actor accepted: %q", actor)
		}
	}
}

func TestDurationMilliseconds(t *testing.T) {
	if got := durationMilliseconds(1500 * time.Millisecond); got != "1500ms" {
		t.Fatalf("unexpected duration: %s", got)
	}
	if got := durationMilliseconds(0); got != "0" {
		t.Fatalf("unexpected zero duration: %s", got)
	}
}

func TestEncodedActorRuntimeParameterIsRejected(t *testing.T) {
	cfg := validConfig()
	cfg.URL = "postgres://example?options=-c%20%61%74%6c%61%73%2e%61%63%74%6f%72%5f%69%64%3devil"
	if err := cfg.validate(); err == nil {
		t.Fatal("encoded session actor configuration was accepted")
	}
}

func TestNilPoolRejectsGovernedTransaction(t *testing.T) {
	var pool *Pool
	if err := pool.WithActor(context.Background(), "actor", func(context.Context, pgx.Tx) error { return nil }); err == nil {
		t.Fatal("nil pool accepted governed transaction")
	}
}
