package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	database "github.com/Iron-Signal-Systems/iron-atlas/internal/database/postgresql"
)

type queryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Resolver maps a verified provider identity to one active governed Atlas actor.
// It reaches PostgreSQL only through the narrowly scoped
// atlas.resolve_governed_actor function and does not require direct SELECT
// privileges on governed identity or role tables.
type Resolver struct {
	database queryRower
}

var _ authentication.ActorResolver = (*Resolver)(nil)

func New(pool *database.Pool) (*Resolver, error) {
	if pool == nil {
		return nil, errors.New("database pool is required")
	}
	return &Resolver{database: pool}, nil
}

func (r *Resolver) Resolve(
	ctx context.Context,
	principal authentication.Principal,
) (authz.Actor, error) {
	if r == nil || r.database == nil {
		return authz.Actor{}, fmt.Errorf(
			"%w: governed actor resolver is unavailable",
			authentication.ErrAuthenticationUnavailable,
		)
	}
	if err := principal.Validate(); err != nil {
		return authz.Actor{}, fmt.Errorf(
			"%w: invalid verified principal",
			authentication.ErrIdentityResolutionFailed,
		)
	}
	if principal.ProviderID != strings.TrimSpace(principal.ProviderID) ||
		principal.Subject != strings.TrimSpace(principal.Subject) {
		return authz.Actor{}, fmt.Errorf(
			"%w: principal identifiers are not normalized",
			authentication.ErrIdentityResolutionFailed,
		)
	}

	var (
		actorID   string
		roleCodes []string
	)
	err := r.database.QueryRow(
		ctx,
		`SELECT actor_id, role_codes
		   FROM atlas.resolve_governed_actor($1, $2)`,
		principal.ProviderID,
		principal.Subject,
	).Scan(&actorID, &roleCodes)
	if errors.Is(err, pgx.ErrNoRows) {
		return authz.Actor{}, authentication.ErrIdentityResolutionFailed
	}
	if err != nil {
		return authz.Actor{}, fmt.Errorf(
			"%w: query governed actor resolution",
			authentication.ErrAuthenticationUnavailable,
		)
	}

	if actorID == "" || actorID != strings.TrimSpace(actorID) {
		return authz.Actor{}, fmt.Errorf(
			"%w: database returned an invalid actor identifier",
			authentication.ErrIdentityResolutionFailed,
		)
	}
	if len(roleCodes) > 32 {
		return authz.Actor{}, fmt.Errorf(
			"%w: database returned too many roles",
			authentication.ErrIdentityResolutionFailed,
		)
	}

	roles := make([]authz.Role, 0, len(roleCodes))
	seen := make(map[authz.Role]struct{}, len(roleCodes))
	for _, code := range roleCodes {
		role, ok := mapRoleCode(code)
		if !ok {
			return authz.Actor{}, fmt.Errorf(
				"%w: database returned an unsupported role",
				authentication.ErrIdentityResolutionFailed,
			)
		}
		if _, duplicate := seen[role]; duplicate {
			return authz.Actor{}, fmt.Errorf(
				"%w: database returned a duplicate role",
				authentication.ErrIdentityResolutionFailed,
			)
		}
		seen[role] = struct{}{}
		roles = append(roles, role)
	}

	return authz.Actor{ID: actorID, Roles: roles}, nil
}

func mapRoleCode(code string) (authz.Role, bool) {
	switch code {
	case "VIEWER":
		return authz.RoleViewer, true
	case "NETWORK_TECHNICIAN":
		return authz.RoleNetworkTech, true
	case "NETWORK_ADMINISTRATOR":
		return authz.RoleNetworkAdmin, true
	case "NETWORK_SECURITY":
		return authz.RoleNetworkSecurity, true
	case "CHANGE_APPROVER":
		return authz.RoleChangeApprover, true
	case "AUDITOR":
		return authz.RoleAuditor, true
	case "PLATFORM_ADMINISTRATOR":
		return authz.RolePlatformAdmin, true
	default:
		return "", false
	}
}
