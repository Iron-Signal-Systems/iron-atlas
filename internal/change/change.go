package change

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

type Status string

type Risk string

const (
	StatusDraft           Status = "draft"
	StatusPendingApproval Status = "pending_approval"
	StatusApproved        Status = "approved"
	StatusRejected        Status = "rejected"
	StatusImplementing    Status = "implementing"
	StatusValidating      Status = "validating"
	StatusAccepted        Status = "accepted"
)

const (
	RiskLow      Risk = "low"
	RiskModerate Risk = "moderate"
	RiskHigh     Risk = "high"
	RiskCritical Risk = "critical"
)

type Request struct {
	ID                string
	Title             string
	Risk              Risk
	Status            Status
	Requester         string
	RequiredApprovals int
	Approvals         []Approval
}

type Approval struct {
	ActorID  string
	Decision string
	Reason   string
	At       time.Time
}

type Service interface {
	List() []Request
	Get(id string) (Request, bool)
	Approve(id string, actor authz.Actor, reason string) (Request, error)
}

type MemoryService struct {
	mu     sync.Mutex
	policy *authz.Policy
	items  map[string]Request
}

func NewMemoryService(policy *authz.Policy) *MemoryService {
	return &MemoryService{policy: policy, items: make(map[string]Request)}
}

func (s *MemoryService) Seed(request Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[request.ID] = request
}

func (s *MemoryService) List() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Request, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *MemoryService) Get(id string) (Request, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *MemoryService) Approve(id string, actor authz.Actor, reason string) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.policy.Require(actor, authz.PermissionApproveChange); err != nil {
		return Request{}, err
	}
	item, ok := s.items[id]
	if !ok {
		return Request{}, errors.New("change request not found")
	}
	if item.Status != StatusPendingApproval {
		return Request{}, fmt.Errorf("change is not pending approval: %s", item.Status)
	}
	if actor.ID == item.Requester {
		return Request{}, errors.New("requester cannot approve their own change")
	}
	for _, approval := range item.Approvals {
		if approval.ActorID == actor.ID && approval.Decision == "approve" {
			return Request{}, errors.New("actor has already approved this change")
		}
	}
	if reason == "" {
		return Request{}, errors.New("approval reason is required")
	}

	item.Approvals = append(item.Approvals, Approval{ActorID: actor.ID, Decision: "approve", Reason: reason, At: time.Now().UTC()})
	if len(item.Approvals) >= item.RequiredApprovals {
		item.Status = StatusApproved
	}
	s.items[id] = item
	return item, nil
}
