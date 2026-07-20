//go:build integration

package postgresql

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication/session"
	database "github.com/Iron-Signal-Systems/iron-atlas/internal/database/postgresql"
)

func integrationStore(t *testing.T) (*Store, *database.Pool, func()) {
	t.Helper()
	url := os.Getenv("IRON_ATLAS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("IRON_ATLAS_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, database.Config{
		URL:               url,
		ApplicationName:   "iron-atlas-session-test",
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
	store, err := New(pool)
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}
	return store, pool, pool.Close
}

func integrationRequest(seed string, actorID string) session.CreateRequest {
	now := time.Now().UTC().Truncate(time.Microsecond)
	digest := sha256.Sum256([]byte(seed))
	return session.CreateRequest{
		IdentifierDigest: digest,
		Principal: authentication.Principal{
			ProviderID:      "dev",
			Subject:         "subject-requester",
			AuthenticatedAt: now.Add(-time.Minute),
			Assurance: authentication.Assurance{
				Context:               "urn:iron-atlas:assurance:provider-mfa",
				Methods:               []string{"pwd", "otp"},
				MFAAuthenticated:      true,
				MFAAuthenticatedAt:    now.Add(-time.Minute),
				SecurityPolicyVersion: "phase-1-step-3-session-v1",
			},
		},
		ActorID:               actorID,
		CreatedAt:             now,
		IdleExpiresAt:         now.Add(30 * time.Minute),
		AbsoluteExpiresAt:     now.Add(8 * time.Hour),
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
	}
}

func TestIntegrationStoreCreatesAndAuthenticatesSession(t *testing.T) {
	store, _, closeStore := integrationStore(t)
	defer closeStore()
	request := integrationRequest(t.Name()+time.Now().String(), "requester")

	created, err := store.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if created.ActorID != "requester" {
		t.Fatalf("created actor = %q", created.ActorID)
	}

	found, err := store.Find(
		context.Background(),
		request.IdentifierDigest,
	)
	if err != nil {
		t.Fatal(err)
	}
	if found.Principal.ProviderID != "dev" ||
		found.Principal.Subject != "subject-requester" ||
		found.ActorID != "requester" ||
		!found.Principal.Assurance.MFAAuthenticated {
		t.Fatal("authenticated session record fields are invalid")
	}
}

func TestIntegrationStoreRejectsActorMismatchAndUnknownLookup(t *testing.T) {
	store, _, closeStore := integrationStore(t)
	defer closeStore()

	mismatch := integrationRequest(t.Name()+"-mismatch", "approver-a")
	_, err := store.Create(context.Background(), mismatch)
	if !errors.Is(err, session.ErrSessionInvalid) {
		t.Fatalf("actor mismatch error = %v", err)
	}

	unknownDigest := sha256.Sum256([]byte(t.Name() + "-unknown"))
	_, err = store.Find(
		context.Background(),
		unknownDigest,
	)
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("unknown lookup error = %v", err)
	}
}

func TestIntegrationApplicationCannotReadSessionTable(t *testing.T) {
	_, pool, closeStore := integrationStore(t)
	defer closeStore()
	var count int
	err := pool.QueryRow(
		context.Background(),
		"SELECT count(*) FROM atlas.authenticated_session",
	).Scan(&count)
	if err == nil {
		t.Fatalf("direct session table read unexpectedly succeeded: %d", count)
	}
}

func TestIntegrationConcurrentSessionLookup(t *testing.T) {
	store, _, closeStore := integrationStore(t)
	defer closeStore()
	request := integrationRequest(t.Name()+time.Now().String(), "requester")
	if _, err := store.Create(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	const operations = 100
	var wait sync.WaitGroup
	errorsFound := make(chan error, operations)
	for range operations {
		wait.Add(1)
		go func() {
			defer wait.Done()
			record, err := store.Find(
				context.Background(),
				request.IdentifierDigest,
			)
			if err != nil {
				errorsFound <- err
				return
			}
			if record.ActorID != "requester" {
				errorsFound <- errors.New("session actor changed during concurrent lookup")
			}
		}()
	}
	wait.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Fatal(err)
	}
}
