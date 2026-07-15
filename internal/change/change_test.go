package change

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

func TestRequesterCannotApproveOwnChange(t *testing.T) {
	policy := authz.DefaultPolicy()
	service := NewMemoryService(policy)
	service.Seed(Request{ID: "CHG-1", Title: "test", Requester: "same", Status: StatusPendingApproval, Risk: RiskModerate, RequiredApprovals: 1})
	_, err := service.Approve(context.Background(), "CHG-1", authz.Actor{ID: "same", Roles: []authz.Role{authz.RoleNetworkAdmin}}, "reason")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestConcurrentDuplicateApprovalHasOneWinner(t *testing.T) {
	policy := authz.DefaultPolicy()
	service := NewMemoryService(policy)
	service.Seed(Request{ID: "CHG-2", Title: "test", Requester: "requester", Status: StatusPendingApproval, Risk: RiskModerate, RequiredApprovals: 2})
	actor := authz.Actor{ID: "approver", Roles: []authz.Role{authz.RoleChangeApprover}}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := service.Approve(context.Background(), "CHG-2", actor, "independent review")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	success, failure := 0, 0
	for err := range errs {
		if err == nil {
			success++
		} else {
			failure++
		}
	}
	if success != 1 || failure != 1 {
		t.Fatalf("expected one success and one failure, got success=%d failure=%d", success, failure)
	}
}

func TestCreateUsesAuthenticatedActor(t *testing.T) {
	service := NewMemoryService(authz.DefaultPolicy())
	actor := authz.Actor{ID: "requester", Roles: []authz.Role{authz.RoleNetworkTech}}
	request, err := service.Create(context.Background(), actor, CreateInput{ID: "CHG-3", Title: "new change", RequiredApprovals: 1})
	if err != nil {
		t.Fatal(err)
	}
	if request.Requester != actor.ID {
		t.Fatalf("expected requester %q, got %q", actor.ID, request.Requester)
	}
}
