package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/atlas/internal/authentication/session"
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

func createRequest() session.CreateRequest {
	now := time.Date(2026, 7, 19, 20, 0, 0, 0, time.UTC)
	return session.CreateRequest{
		Principal: authentication.Principal{
			ProviderID:      "oidc-test",
			Subject:         "subject-123",
			AuthenticatedAt: now.Add(-time.Minute),
			Assurance: authentication.Assurance{
				Methods:               []string{"pwd", "otp"},
				MFAAuthenticated:      true,
				MFAAuthenticatedAt:    now.Add(-time.Minute),
				SecurityPolicyVersion: "phase-1-step-3-session-v1",
			},
		},
		ActorID:               "actor-123",
		CreatedAt:             now,
		IdleExpiresAt:         now.Add(30 * time.Minute),
		AbsoluteExpiresAt:     now.Add(8 * time.Hour),
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
	}
}

func TestStoreCreateUsesControlledFunction(t *testing.T) {
	request := createRequest()
	store := &Store{database: fakeQuerier{row: fakeRow{
		scan: func(destinations ...any) error {
			*destinations[0].(*time.Time) = request.CreatedAt
			*destinations[1].(*time.Time) = request.CreatedAt
			*destinations[2].(*time.Time) = request.IdleExpiresAt
			*destinations[3].(*time.Time) = request.AbsoluteExpiresAt
			return nil
		},
	}}}

	record, err := store.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if record.ActorID != "actor-123" ||
		record.SecurityPolicyVersion != "phase-1-step-3-session-v1" {
		t.Fatal("controlled session record fields are invalid")
	}
}

func TestStoreCreateRejectsInvalidRequestBeforeDatabase(t *testing.T) {
	store := &Store{database: fakeQuerier{row: fakeRow{
		scan: func(...any) error {
			t.Fatal("database must not be queried")
			return nil
		},
	}}}
	request := createRequest()
	request.ActorID = " actor-123"
	_, err := store.Create(context.Background(), request)
	if !errors.Is(err, session.ErrSessionInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestStoreFindMapsControlledResult(t *testing.T) {
	now := time.Date(2026, 7, 19, 20, 0, 0, 0, time.UTC)
	store := &Store{database: fakeQuerier{row: fakeRow{
		scan: func(destinations ...any) error {
			*destinations[0].(*string) = "oidc-test"
			*destinations[1].(*string) = "subject-123"
			*destinations[2].(*string) = "actor-123"
			*destinations[3].(*time.Time) = now.Add(-time.Minute)
			*destinations[4].(*time.Time) = now.Add(-time.Minute)
			*destinations[5].(*time.Time) = now.Add(-time.Minute)
			*destinations[6].(*time.Time) = now.Add(30 * time.Minute)
			*destinations[7].(*time.Time) = now.Add(8 * time.Hour)
			*destinations[8].(*pgtype.Text) = pgtype.Text{
				String: "urn:iron-atlas:assurance:provider-mfa",
				Valid:  true,
			}
			*destinations[9].(*[]string) = []string{"pwd", "otp"}
			*destinations[10].(*bool) = true
			*destinations[11].(*pgtype.Timestamptz) = pgtype.Timestamptz{
				Time:  now.Add(-time.Minute),
				Valid: true,
			}
			*destinations[12].(*string) = "phase-1-step-3-session-v1"
			return nil
		},
	}}}

	record, err := store.Find(context.Background(), [32]byte{1})
	if err != nil {
		t.Fatal(err)
	}
	if record.Principal.ProviderID != "oidc-test" ||
		record.ActorID != "actor-123" ||
		!record.Principal.Assurance.MFAAuthenticated {
		t.Fatal("controlled session record fields are invalid")
	}
}

func TestStoreFindClassifiesMissingAndOutage(t *testing.T) {
	for name, testCase := range map[string]struct {
		rowError error
		expected error
	}{
		"missing": {rowError: pgx.ErrNoRows, expected: session.ErrSessionNotFound},
		"outage":  {rowError: errors.New("database offline"), expected: session.ErrSessionUnavailable},
	} {
		t.Run(name, func(t *testing.T) {
			store := &Store{database: fakeQuerier{row: fakeRow{
				scan: func(...any) error { return testCase.rowError },
			}}}
			_, err := store.Find(context.Background(), [32]byte{1})
			if !errors.Is(err, testCase.expected) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
