package change

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Iron-Signal-Systems/atlas/internal/authz"
)

var (
	ErrNotFound  = errors.New("change request not found")
	ErrForbidden = errors.New("change operation forbidden")
	ErrConflict  = errors.New("change state conflict")
	ErrInvalid   = errors.New("invalid change request")
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
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	Risk              Risk       `json:"risk"`
	Status            Status     `json:"status"`
	Requester         string     `json:"requester"`
	RequiredApprovals int        `json:"required_approvals"`
	ApprovalCount     int        `json:"approval_count"`
	Approvals         []Approval `json:"approvals,omitempty"`
}

type Approval struct {
	ActorID  string    `json:"actor_id"`
	Decision string    `json:"decision"`
	Reason   string    `json:"reason"`
	At       time.Time `json:"at"`
}

type CreateInput struct {
	ID                string
	Title             string
	RequiredApprovals int
}

type Service interface {
	List(context.Context) ([]Request, error)
	Get(context.Context, string) (Request, bool, error)
	Create(context.Context, authz.Actor, CreateInput) (Request, error)
	Approve(context.Context, string, authz.Actor, string) (Request, error)
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
	request.ApprovalCount = len(request.Approvals)
	s.items[request.ID] = request
}

func (s *MemoryService) List(context.Context) ([]Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Request, 0, len(s.items))
	for _, item := range s.items {
		item.ApprovalCount = len(item.Approvals)
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *MemoryService) Get(_ context.Context, id string) (Request, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	item.ApprovalCount = len(item.Approvals)
	return item, ok, nil
}

func (s *MemoryService) Create(_ context.Context, actor authz.Actor, input CreateInput) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.policy.Require(actor, authz.PermissionRequestChange); err != nil {
		return Request{}, fmt.Errorf("%w: %v", ErrForbidden, err)
	}
	input.ID = strings.TrimSpace(input.ID)
	input.Title = strings.TrimSpace(input.Title)
	if input.ID == "" || input.Title == "" || input.RequiredApprovals < 1 {
		return Request{}, ErrInvalid
	}
	if _, exists := s.items[input.ID]; exists {
		return Request{}, ErrConflict
	}
	request := Request{
		ID:                input.ID,
		Title:             input.Title,
		Risk:              RiskModerate,
		Status:            StatusPendingApproval,
		Requester:         actor.ID,
		RequiredApprovals: input.RequiredApprovals,
	}
	s.items[request.ID] = request
	return request, nil
}

func (s *MemoryService) Approve(_ context.Context, id string, actor authz.Actor, reason string) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.policy.Require(actor, authz.PermissionApproveChange); err != nil {
		return Request{}, fmt.Errorf("%w: %v", ErrForbidden, err)
	}
	item, ok := s.items[id]
	if !ok {
		return Request{}, ErrNotFound
	}
	if item.Status != StatusPendingApproval {
		return Request{}, fmt.Errorf("%w: change is not pending approval: %s", ErrConflict, item.Status)
	}
	if actor.ID == item.Requester {
		return Request{}, fmt.Errorf("%w: requester cannot approve their own change", ErrForbidden)
	}
	for _, approval := range item.Approvals {
		if approval.ActorID == actor.ID && approval.Decision == "approve" {
			return Request{}, fmt.Errorf("%w: actor has already approved this change", ErrConflict)
		}
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Request{}, fmt.Errorf("%w: approval reason is required", ErrInvalid)
	}

	item.Approvals = append(item.Approvals, Approval{ActorID: actor.ID, Decision: "approve", Reason: reason, At: time.Now().UTC()})
	item.ApprovalCount = len(item.Approvals)
	if item.ApprovalCount >= item.RequiredApprovals {
		item.Status = StatusApproved
	}
	s.items[id] = item
	return item, nil
}
