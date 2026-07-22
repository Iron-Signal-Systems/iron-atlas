//go:build integration

package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/atlas/internal/authz"
	"github.com/Iron-Signal-Systems/atlas/internal/change"
	database "github.com/Iron-Signal-Systems/atlas/internal/database/postgresql"
)

func integrationService(t *testing.T) (*Service, *database.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, database.Config{
		URL:               os.Getenv("IRON_ATLAS_TEST_DATABASE_URL"),
		ApplicationName:   "iron-atlas-change-step2-test",
		MaxConnections:    2,
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
	service, err := New(pool)
	if err != nil {
		t.Fatal(err)
	}
	return service, pool
}

func TestGovernedCreateAndIndependentApproval(t *testing.T) {
	service, _ := integrationService(t)
	ctx := context.Background()
	id := fmt.Sprintf("CHG-GO-%d", time.Now().UnixNano())
	requester := authz.Actor{ID: "requester"}
	approver := authz.Actor{ID: "approver-a"}

	created, err := service.Create(ctx, requester, change.CreateInput{
		ID: id, Title: "Go PostgreSQL integration", RequiredApprovals: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Requester != requester.ID || created.Status != change.StatusPendingApproval {
		t.Fatalf("unexpected created request: %+v", created)
	}

	if _, err := service.Approve(ctx, id, requester, "self approval"); !errors.Is(err, change.ErrForbidden) {
		t.Fatalf("expected self-approval denial, got %v", err)
	}

	approved, err := service.Approve(ctx, id, approver, "independent review completed")
	if err != nil {
		t.Fatal(err)
	}
	if approved.Status != change.StatusApproved || approved.ApprovalCount != 1 {
		t.Fatalf("unexpected approved request: %+v", approved)
	}
}

func TestFailedGovernedCallRollsBackAndNextActorSucceeds(t *testing.T) {
	service, _ := integrationService(t)
	ctx := context.Background()
	id := fmt.Sprintf("CHG-GO-FAIL-%d", time.Now().UnixNano())
	_, err := service.Create(ctx, authz.Actor{ID: "requester"}, change.CreateInput{
		ID: id, Title: "failure isolation", RequiredApprovals: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.Approve(ctx, id, authz.Actor{ID: "unauthorized"}, "must fail"); !errors.Is(err, change.ErrForbidden) {
		t.Fatalf("expected unauthorized actor denial, got %v", err)
	}

	approved, err := service.Approve(ctx, id, authz.Actor{ID: "approver-b"}, "authorized after failed transaction")
	if err != nil {
		t.Fatal(err)
	}
	if approved.Status != change.StatusApproved {
		t.Fatalf("expected approved status, got %s", approved.Status)
	}
}
