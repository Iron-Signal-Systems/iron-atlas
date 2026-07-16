//go:build integration

package postgresql

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	database "github.com/Iron-Signal-Systems/iron-atlas/internal/database/postgresql"
)

func integrationResolver(t *testing.T) (*Resolver, func()) {
	t.Helper()
	url := os.Getenv("IRON_ATLAS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("IRON_ATLAS_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, database.Config{
		URL:               url,
		ApplicationName:   "iron-atlas-actor-resolution-test",
		MaxConnections:    8,
		MinConnections:    0,
		ConnectTimeout:    5 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
		MaxConnectionIdle: time.Minute,
		HealthInterval:    time.Minute,
		StatementTimeout:  10 * time.Second,
		LockTimeout:       5 * time.Second,
		IdleInTxTimeout:   10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := New(pool)
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}
	return resolver, pool.Close
}

func verifiedPrincipal(subject string) authentication.Principal {
	return authentication.Principal{
		ProviderID:      "dev",
		Subject:         subject,
		AuthenticatedAt: time.Now().UTC(),
	}
}

func TestIntegrationResolverLoadsActiveActorAndRoles(t *testing.T) {
	resolver, closeResolver := integrationResolver(t)
	defer closeResolver()

	actor, err := resolver.Resolve(
		context.Background(),
		verifiedPrincipal("subject-requester"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if actor.ID != "requester" {
		t.Fatalf("actor ID = %q", actor.ID)
	}
	if len(actor.Roles) != 1 || actor.Roles[0] != authz.RoleNetworkTech {
		t.Fatalf("roles = %#v", actor.Roles)
	}
}

func TestIntegrationResolverRejectsInactiveOrUnknownState(t *testing.T) {
	resolver, closeResolver := integrationResolver(t)
	defer closeResolver()

	for _, subject := range []string{
		"subject-disabled-actor",
		"subject-inactive-provider",
		"subject-unmapped",
		"subject-unknown-role",
	} {
		t.Run(subject, func(t *testing.T) {
			_, err := resolver.Resolve(
				context.Background(),
				verifiedPrincipal(subject),
			)
			if !errors.Is(err, authentication.ErrIdentityResolutionFailed) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestIntegrationResolverExcludesExpiredAndInactiveRoles(t *testing.T) {
	resolver, closeResolver := integrationResolver(t)
	defer closeResolver()

	for _, subject := range []string{
		"subject-expired-role",
		"subject-inactive-role",
		"subject-no-role",
	} {
		t.Run(subject, func(t *testing.T) {
			actor, err := resolver.Resolve(
				context.Background(),
				verifiedPrincipal(subject),
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(actor.Roles) != 0 {
				t.Fatalf("roles = %#v", actor.Roles)
			}
		})
	}
}

func TestIntegrationResolverConcurrentIsolation(t *testing.T) {
	resolver, closeResolver := integrationResolver(t)
	defer closeResolver()

	const operations = 200
	var wait sync.WaitGroup
	errorsFound := make(chan error, operations)
	for index := 0; index < operations; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			actor, err := resolver.Resolve(
				context.Background(),
				verifiedPrincipal("subject-requester"),
			)
			if err != nil {
				errorsFound <- err
				return
			}
			if actor.ID != "requester" ||
				len(actor.Roles) != 1 ||
				actor.Roles[0] != authz.RoleNetworkTech {
				errorsFound <- errors.New("resolved identity was not isolated")
			}
		}()
	}
	wait.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Fatal(err)
	}
}
