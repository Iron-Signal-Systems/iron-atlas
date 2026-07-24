package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Iron-Signal-Systems/atlas/internal/authz"
	"github.com/Iron-Signal-Systems/atlas/internal/change"
	database "github.com/Iron-Signal-Systems/atlas/internal/database/postgresql"
)

// Service persists governed changes through the accepted PostgreSQL function
// boundary. It never writes governed tables directly.
type Service struct {
	database *database.Pool
}

func New(pool *database.Pool) (*Service, error) {
	if pool == nil {
		return nil, errors.New("database pool is required")
	}
	return &Service{database: pool}, nil
}

func (s *Service) List(ctx context.Context) ([]change.Request, error) {
	rows, err := s.database.Query(ctx, `
		SELECT cr.change_id,
		       cr.title,
		       cr.risk,
		       cr.status,
		       cr.requester_actor_id,
		       cr.required_approvals,
		       COALESCE(summary.approval_count, 0)
		FROM atlas.change_request AS cr
		LEFT JOIN atlas.change_approval_summary AS summary
		  ON summary.change_id = cr.change_id
		ORDER BY cr.change_id`)
	if err != nil {
		return nil, fmt.Errorf("list change requests: %w", err)
	}
	defer rows.Close()

	items := make([]change.Request, 0)
	for rows.Next() {
		item, err := scanRequest(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change requests: %w", err)
	}
	return items, nil
}

func (s *Service) Get(ctx context.Context, id string) (change.Request, bool, error) {
	item, err := scanRequest(s.database.QueryRow(ctx, `
		SELECT cr.change_id,
		       cr.title,
		       cr.risk,
		       cr.status,
		       cr.requester_actor_id,
		       cr.required_approvals,
		       COALESCE(summary.approval_count, 0)
		FROM atlas.change_request AS cr
		LEFT JOIN atlas.change_approval_summary AS summary
		  ON summary.change_id = cr.change_id
		WHERE cr.change_id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return change.Request{}, false, nil
	}
	if err != nil {
		return change.Request{}, false, err
	}
	return item, true, nil
}

func (s *Service) Create(ctx context.Context, actor authz.Actor, input change.CreateInput) (change.Request, error) {
	if strings.TrimSpace(actor.ID) == "" {
		return change.Request{}, change.ErrForbidden
	}
	input.ID = strings.TrimSpace(input.ID)
	input.Title = strings.TrimSpace(input.Title)
	if input.ID == "" || input.Title == "" || input.RequiredApprovals < 1 {
		return change.Request{}, change.ErrInvalid
	}

	err := s.database.WithActor(ctx, actor.ID, func(ctx context.Context, tx pgx.Tx) error {
		var createdID string
		if err := tx.QueryRow(ctx,
			"SELECT atlas.create_change_request($1, $2, $3)",
			input.ID,
			input.Title,
			input.RequiredApprovals,
		).Scan(&createdID); err != nil {
			return classifyDatabaseError(err)
		}
		if createdID != input.ID {
			return errors.New("database returned an unexpected change identifier")
		}
		return nil
	})
	if err != nil {
		return change.Request{}, err
	}

	request, found, err := s.Get(ctx, input.ID)
	if err != nil {
		return change.Request{}, err
	}
	if !found {
		return change.Request{}, errors.New("created change was not readable after commit")
	}
	return request, nil
}

func (s *Service) Approve(ctx context.Context, id string, actor authz.Actor, reason string) (change.Request, error) {
	if strings.TrimSpace(actor.ID) == "" {
		return change.Request{}, change.ErrForbidden
	}
	id = strings.TrimSpace(id)
	reason = strings.TrimSpace(reason)
	if id == "" || reason == "" {
		return change.Request{}, change.ErrInvalid
	}

	err := s.database.WithActor(ctx, actor.ID, func(ctx context.Context, tx pgx.Tx) error {
		var actionID int64
		if err := tx.QueryRow(ctx,
			"SELECT atlas.record_approval($1, 'APPROVE', $2)",
			id,
			reason,
		).Scan(&actionID); err != nil {
			return classifyDatabaseError(err)
		}
		if actionID < 1 {
			return errors.New("database returned an invalid approval action identifier")
		}
		return nil
	})
	if err != nil {
		return change.Request{}, err
	}

	request, found, err := s.Get(ctx, id)
	if err != nil {
		return change.Request{}, err
	}
	if !found {
		return change.Request{}, change.ErrNotFound
	}
	return request, nil
}

type rowScanner interface {
	Scan(...any) error
}

func scanRequest(row rowScanner) (change.Request, error) {
	var (
		item          change.Request
		risk          string
		status        string
		approvalCount int
	)
	if err := row.Scan(
		&item.ID,
		&item.Title,
		&risk,
		&status,
		&item.Requester,
		&item.RequiredApprovals,
		&approvalCount,
	); err != nil {
		return change.Request{}, err
	}
	item.Risk = change.Risk(strings.ToLower(risk))
	item.Status = change.Status(strings.ToLower(status))
	item.ApprovalCount = approvalCount
	return item, nil
}

func classifyDatabaseError(err error) error {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "change not found"):
		return fmt.Errorf("%w: change not found", change.ErrNotFound)
	case strings.Contains(message, "lacks"):
		return fmt.Errorf("%w: actor lacks required authority", change.ErrForbidden)
	case strings.Contains(message, "requester cannot"):
		return fmt.Errorf("%w: requester cannot approve own change", change.ErrForbidden)
	case strings.Contains(message, "authenticated actor"):
		return fmt.Errorf("%w: authenticated actor is required", change.ErrForbidden)
	case strings.Contains(message, "already"):
		return fmt.Errorf("%w: actor already has an active approval", change.ErrConflict)
	case strings.Contains(message, "approvable state"):
		return fmt.Errorf("%w: change is not in an approvable state", change.ErrConflict)
	case strings.Contains(message, "duplicate key"):
		return fmt.Errorf("%w: duplicate governed record", change.ErrConflict)
	case strings.Contains(message, "required") || strings.Contains(message, "invalid") || strings.Contains(message, "check constraint"):
		return fmt.Errorf("%w: database rejected invalid change input", change.ErrInvalid)
	default:
		return fmt.Errorf("database change operation failed: %w", err)
	}
}
