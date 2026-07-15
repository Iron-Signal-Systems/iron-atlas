package change

import (
	"sync"
	"testing"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

func TestRequesterCannotApproveOwnChange(t *testing.T) {
	service := NewMemoryService(authz.DefaultPolicy())
	service.Seed(Request{ID: "CHG-1", Title: "test", Status: StatusPendingApproval, Requester: "alice", RequiredApprovals: 1})

	_, err := service.Approve("CHG-1", authz.Actor{ID: "alice", Roles: []authz.Role{authz.RoleNetworkAdmin}}, "looks good")
	if err == nil {
		t.Fatal("expected requester independence failure")
	}
}

func TestIndependentApproverCanApprove(t *testing.T) {
	service := NewMemoryService(authz.DefaultPolicy())
	service.Seed(Request{ID: "CHG-1", Title: "test", Status: StatusPendingApproval, Requester: "alice", RequiredApprovals: 1})

	request, err := service.Approve("CHG-1", authz.Actor{ID: "bob", Roles: []authz.Role{authz.RoleChangeApprover}}, "reviewed scope and rollback")
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if request.Status != StatusApproved {
		t.Fatalf("expected approved, got %s", request.Status)
	}
}

func TestConcurrentDuplicateApprovalProducesOneApproval(t *testing.T) {
	service := NewMemoryService(authz.DefaultPolicy())
	service.Seed(Request{ID: "CHG-1", Title: "test", Status: StatusPendingApproval, Requester: "alice", RequiredApprovals: 2})
	actor := authz.Actor{ID: "bob", Roles: []authz.Role{authz.RoleChangeApprover}}

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			_, _ = service.Approve("CHG-1", actor, "reviewed")
		}()
	}
	wg.Wait()

	request, _ := service.Get("CHG-1")
	if len(request.Approvals) != 1 {
		t.Fatalf("expected one effective approval, got %d", len(request.Approvals))
	}
}
