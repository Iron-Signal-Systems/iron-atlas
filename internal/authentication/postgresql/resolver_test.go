package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/atlas/internal/authz"
)

type fakeQuerier struct {
	row pgx.Row
}

func (f fakeQuerier) QueryRow(
	context.Context,
	string,
	...any,
) pgx.Row {
	return f.row
}

type fakeRow struct {
	scan func(...any) error
}

func (f fakeRow) Scan(destinations ...any) error {
	return f.scan(destinations...)
}

func principal() authentication.Principal {
	return authentication.Principal{
		ProviderID:      "oidc-test",
		Subject:         "subject-123",
		AuthenticatedAt: time.Date(2026, 7, 16, 11, 0, 0, 0, time.UTC),
	}
}

func TestResolverMapsGovernedRoles(t *testing.T) {
	resolver := &Resolver{database: fakeQuerier{row: fakeRow{
		scan: func(destinations ...any) error {
			*destinations[0].(*string) = "actor-123"
			*destinations[1].(*[]string) = []string{
				"AUDITOR",
				"NETWORK_TECHNICIAN",
			}
			return nil
		},
	}}}

	actor, err := resolver.Resolve(context.Background(), principal())
	if err != nil {
		t.Fatal(err)
	}
	if actor.ID != "actor-123" {
		t.Fatalf("actor ID = %q", actor.ID)
	}
	if len(actor.Roles) != 2 ||
		actor.Roles[0] != authz.RoleAuditor ||
		actor.Roles[1] != authz.RoleNetworkTech {
		t.Fatalf("unexpected roles: %#v", actor.Roles)
	}
}

func TestResolverFailsClosedForMissingMapping(t *testing.T) {
	resolver := &Resolver{database: fakeQuerier{row: fakeRow{
		scan: func(...any) error { return pgx.ErrNoRows },
	}}}

	_, err := resolver.Resolve(context.Background(), principal())
	if !errors.Is(err, authentication.ErrIdentityResolutionFailed) {
		t.Fatalf("error = %v", err)
	}
}

func TestResolverClassifiesDatabaseFailureAsUnavailable(t *testing.T) {
	resolver := &Resolver{database: fakeQuerier{row: fakeRow{
		scan: func(...any) error { return errors.New("database offline") },
	}}}

	_, err := resolver.Resolve(context.Background(), principal())
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("error = %v", err)
	}
}

func TestResolverRejectsUnknownOrDuplicateRoles(t *testing.T) {
	for name, roleCodes := range map[string][]string{
		"unknown":   {"UNRESTRICTED_EXECUTION_CONTEXT"},
		"duplicate": {"AUDITOR", "AUDITOR"},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &Resolver{database: fakeQuerier{row: fakeRow{
				scan: func(destinations ...any) error {
					*destinations[0].(*string) = "actor-123"
					*destinations[1].(*[]string) = roleCodes
					return nil
				},
			}}}
			_, err := resolver.Resolve(context.Background(), principal())
			if !errors.Is(err, authentication.ErrIdentityResolutionFailed) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestResolverRejectsUnnormalizedPrincipal(t *testing.T) {
	resolver := &Resolver{database: fakeQuerier{row: fakeRow{
		scan: func(...any) error {
			t.Fatal("database must not be queried")
			return nil
		},
	}}}
	untrusted := principal()
	untrusted.Subject = " subject-123"

	_, err := resolver.Resolve(context.Background(), untrusted)
	if !errors.Is(err, authentication.ErrIdentityResolutionFailed) {
		t.Fatalf("error = %v", err)
	}
}
