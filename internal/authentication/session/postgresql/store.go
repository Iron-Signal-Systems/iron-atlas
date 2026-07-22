package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication/session"
	database "github.com/Iron-Signal-Systems/atlas/internal/database/postgresql"
)

type queryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type Store struct {
	database queryRower
}

var _ session.Store = (*Store)(nil)

func New(pool *database.Pool) (*Store, error) {
	if pool == nil {
		return nil, errors.New("database pool is required")
	}
	return &Store{database: pool}, nil
}

func (s *Store) Create(
	ctx context.Context,
	request session.CreateRequest,
) (session.Record, error) {
	if s == nil || s.database == nil {
		return session.Record{}, session.ErrSessionUnavailable
	}
	if err := validateCreateRequest(request); err != nil {
		return session.Record{}, fmt.Errorf("%w: %v", session.ErrSessionInvalid, err)
	}

	var record session.Record
	record.Principal = request.Principal
	record.ActorID = request.ActorID
	record.SecurityPolicyVersion = request.SecurityPolicyVersion

	idleSeconds := durationSeconds(request.IdleExpiresAt.Sub(request.CreatedAt))
	absoluteSeconds := durationSeconds(request.AbsoluteExpiresAt.Sub(request.CreatedAt))
	err := s.database.QueryRow(
		ctx,
		`SELECT created_at,
		        last_activity_at,
		        idle_expires_at,
		        absolute_expires_at
		   FROM atlas.create_authenticated_session(
		        $1, $2, $3, $4, $5, $6,
		        $7, $8, $9, $10, $11, $12
		   )`,
		request.IdentifierDigest[:],
		request.Principal.ProviderID,
		request.Principal.Subject,
		request.ActorID,
		request.Principal.AuthenticatedAt.UTC(),
		idleSeconds,
		absoluteSeconds,
		nullableText(request.Principal.Assurance.Context),
		append([]string(nil), request.Principal.Assurance.Methods...),
		request.Principal.Assurance.MFAAuthenticated,
		nullableTime(request.Principal.Assurance.MFAAuthenticatedAt),
		request.SecurityPolicyVersion,
	).Scan(
		&record.CreatedAt,
		&record.LastActivityAt,
		&record.IdleExpiresAt,
		&record.AbsoluteExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return session.Record{}, session.ErrSessionInvalid
	}
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return session.Record{}, session.ErrSessionConflict
		}
		return session.Record{}, fmt.Errorf(
			"%w: create authenticated session",
			session.ErrSessionUnavailable,
		)
	}
	return record, nil
}

func (s *Store) Find(
	ctx context.Context,
	digest [32]byte,
) (session.Record, error) {
	if s == nil || s.database == nil {
		return session.Record{}, session.ErrSessionUnavailable
	}
	var (
		record             session.Record
		contextValue       pgtype.Text
		methods            []string
		mfaAuthenticated   bool
		mfaAuthenticatedAt pgtype.Timestamptz
	)
	err := s.database.QueryRow(
		ctx,
		`SELECT provider_id,
		        provider_subject,
		        actor_id,
		        authenticated_at,
		        created_at,
		        last_activity_at,
		        idle_expires_at,
		        absolute_expires_at,
		        authentication_context,
		        authentication_methods,
		        mfa_authenticated,
		        mfa_authenticated_at,
		        security_policy_version
		   FROM atlas.authenticate_session($1)`,
		digest[:],
	).Scan(
		&record.Principal.ProviderID,
		&record.Principal.Subject,
		&record.ActorID,
		&record.Principal.AuthenticatedAt,
		&record.CreatedAt,
		&record.LastActivityAt,
		&record.IdleExpiresAt,
		&record.AbsoluteExpiresAt,
		&contextValue,
		&methods,
		&mfaAuthenticated,
		&mfaAuthenticatedAt,
		&record.SecurityPolicyVersion,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return session.Record{}, session.ErrSessionNotFound
	}
	if err != nil {
		return session.Record{}, fmt.Errorf(
			"%w: authenticate server-side session",
			session.ErrSessionUnavailable,
		)
	}

	if contextValue.Valid {
		record.Principal.Assurance.Context = contextValue.String
	}
	record.Principal.Assurance.Methods = append([]string(nil), methods...)
	record.Principal.Assurance.MFAAuthenticated = mfaAuthenticated
	if mfaAuthenticatedAt.Valid {
		record.Principal.Assurance.MFAAuthenticatedAt = mfaAuthenticatedAt.Time.UTC()
	}
	record.Principal.Assurance.SecurityPolicyVersion = record.SecurityPolicyVersion
	return record, nil
}

func validateCreateRequest(request session.CreateRequest) error {
	if err := request.Principal.Validate(); err != nil {
		return err
	}
	if request.ActorID == "" || request.ActorID != strings.TrimSpace(request.ActorID) || len(request.ActorID) > 256 {
		return errors.New("actor identifier is invalid")
	}
	if request.CreatedAt.IsZero() || request.IdleExpiresAt.IsZero() || request.AbsoluteExpiresAt.IsZero() {
		return errors.New("session timestamps are required")
	}
	if !request.IdleExpiresAt.After(request.CreatedAt) ||
		!request.AbsoluteExpiresAt.After(request.CreatedAt) ||
		request.IdleExpiresAt.After(request.AbsoluteExpiresAt) {
		return errors.New("session lifetime is invalid")
	}
	if request.SecurityPolicyVersion == "" ||
		request.SecurityPolicyVersion != strings.TrimSpace(request.SecurityPolicyVersion) ||
		len(request.SecurityPolicyVersion) > 128 {
		return errors.New("security policy version is invalid")
	}
	if !request.Principal.Assurance.MFAAuthenticated ||
		request.Principal.Assurance.MFAAuthenticatedAt.IsZero() {
		return errors.New("session MFA assurance is required")
	}
	if request.Principal.Assurance.SecurityPolicyVersion != request.SecurityPolicyVersion {
		return errors.New("principal and session security policy versions differ")
	}
	return nil
}

func durationSeconds(value time.Duration) int32 {
	seconds := value / time.Second
	if value%time.Second != 0 {
		seconds++
	}
	return int32(seconds)
}

func nullableText(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value.UTC()
}
