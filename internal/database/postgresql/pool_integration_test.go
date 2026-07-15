//go:build integration

package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func integrationPool(t *testing.T, maxConns int32) *Pool {
	t.Helper()
	url := os.Getenv("IRON_ATLAS_TEST_DATABASE_URL")
	if url == "" {
		t.Fatal("IRON_ATLAS_TEST_DATABASE_URL is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := Open(ctx, Config{
		URL:               url,
		ApplicationName:   "iron-atlas-step2-test",
		MaxConnections:    maxConns,
		MinConnections:    0,
		ConnectTimeout:    5 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
		MaxConnectionIdle: time.Minute,
		HealthInterval:    10 * time.Second,
		StatementTimeout:  10 * time.Second,
		LockTimeout:       2 * time.Second,
		IdleInTxTimeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestTransactionLocalActorDoesNotLeakFromSingleConnection(t *testing.T) {
	pool := integrationPool(t, 1)
	ctx := context.Background()

	for i := 0; i < 500; i++ {
		actor := fmt.Sprintf("actor-%d", i%7)
		err := pool.WithActor(ctx, actor, func(ctx context.Context, tx pgx.Tx) error {
			var observed string
			if err := tx.QueryRow(ctx,
				"SELECT current_setting('atlas.actor_id', true)",
			).Scan(&observed); err != nil {
				return err
			}
			if observed != actor {
				return fmt.Errorf("expected actor %q, got %q", actor, observed)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}

		var residual *string
		if err := pool.QueryRow(ctx,
			"SELECT NULLIF(current_setting('atlas.actor_id', true), '')",
		).Scan(&residual); err != nil {
			t.Fatalf("read residual actor: %v", err)
		}
		if residual != nil {
			t.Fatalf("actor context leaked after iteration %d: %q", i, *residual)
		}
	}
}

func TestRollbackClearsActorAndData(t *testing.T) {
	pool := integrationPool(t, 1)
	ctx := context.Background()
	marker := fmt.Sprintf("rollback-%d", time.Now().UnixNano())
	sentinel := errors.New("force rollback")

	err := pool.WithActor(ctx, "requester", func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx,
			"SELECT atlas.create_change_request($1, $2, 1)",
			marker,
			"rollback verification",
		); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected rollback sentinel, got %v", err)
	}

	var count int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM atlas.change_request WHERE change_id=$1",
		marker,
	).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rolled back change persisted: %d", count)
	}

	var residual *string
	if err := pool.QueryRow(ctx,
		"SELECT NULLIF(current_setting('atlas.actor_id', true), '')",
	).Scan(&residual); err != nil {
		t.Fatal(err)
	}
	if residual != nil {
		t.Fatalf("actor context leaked after rollback: %q", *residual)
	}
}

func TestConcurrentActorsRemainTransactionScoped(t *testing.T) {
	pool := integrationPool(t, 4)
	ctx := context.Background()
	const workers = 8
	const iterations = 75

	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for worker := 0; worker < workers; worker++ {
		worker := worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			actor := fmt.Sprintf("concurrent-actor-%d", worker)
			for i := 0; i < iterations; i++ {
				if err := pool.WithActor(ctx, actor, func(ctx context.Context, tx pgx.Tx) error {
					var observed string
					if err := tx.QueryRow(ctx,
						"SELECT current_setting('atlas.actor_id', true)",
					).Scan(&observed); err != nil {
						return err
					}
					if observed != actor {
						return fmt.Errorf("actor mismatch: expected %s got %s", actor, observed)
					}
					return nil
				}); err != nil {
					errs <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}
