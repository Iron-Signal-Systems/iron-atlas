package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
)

const (
	CookieName              = "__Host-iron_atlas_session"
	identifierBytes         = 32
	encodedIdentifierBytes  = 43
	defaultIdleLifetime     = 30 * time.Minute
	defaultAbsoluteLifetime = 8 * time.Hour
	maximumAbsoluteLifetime = 24 * time.Hour
)

var (
	ErrSessionInvalid     = errors.New("authenticated session is invalid")
	ErrSessionNotFound    = errors.New("authenticated session was not found")
	ErrSessionConflict    = errors.New("authenticated session conflicts with existing state")
	ErrSessionUnavailable = errors.New("authenticated session store is unavailable")
)

type CreateRequest struct {
	IdentifierDigest      [sha256.Size]byte
	Principal             authentication.Principal
	ActorID               string
	CreatedAt             time.Time
	IdleExpiresAt         time.Time
	AbsoluteExpiresAt     time.Time
	SecurityPolicyVersion string
}

type Record struct {
	Principal             authentication.Principal
	ActorID               string
	CreatedAt             time.Time
	LastActivityAt        time.Time
	IdleExpiresAt         time.Time
	AbsoluteExpiresAt     time.Time
	RevokedAt             time.Time
	RevocationReason      string
	SecurityPolicyVersion string
}

type Store interface {
	Create(context.Context, CreateRequest) (Record, error)
	Find(context.Context, [sha256.Size]byte) (Record, error)
}

type Config struct {
	Store                 Store
	Resolver              authentication.ActorResolver
	Random                io.Reader
	Now                   func() time.Time
	IdleLifetime          time.Duration
	AbsoluteLifetime      time.Duration
	SuccessLocation       string
	SecurityPolicyVersion string
	RequireMFA            bool
}

type Service struct {
	store                 Store
	resolver              authentication.ActorResolver
	random                io.Reader
	now                   func() time.Time
	idleLifetime          time.Duration
	absoluteLifetime      time.Duration
	successLocation       string
	securityPolicyVersion string
	requireMFA            bool
}

var _ authentication.Authenticator = (*Service)(nil)

func New(config Config) (*Service, error) {
	if config.Store == nil {
		return nil, errors.New("session store is required")
	}
	if config.Resolver == nil {
		return nil, errors.New("governed actor resolver is required")
	}
	if config.Random == nil {
		config.Random = rand.Reader
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.IdleLifetime == 0 {
		config.IdleLifetime = defaultIdleLifetime
	}
	if config.AbsoluteLifetime == 0 {
		config.AbsoluteLifetime = defaultAbsoluteLifetime
	}
	if config.IdleLifetime <= 0 {
		return nil, errors.New("session idle lifetime must be positive")
	}
	if config.AbsoluteLifetime <= 0 || config.AbsoluteLifetime > maximumAbsoluteLifetime {
		return nil, fmt.Errorf(
			"session absolute lifetime must be positive and no greater than %s",
			maximumAbsoluteLifetime,
		)
	}
	if config.IdleLifetime > config.AbsoluteLifetime {
		return nil, errors.New("session idle lifetime must not exceed absolute lifetime")
	}
	if strings.TrimSpace(config.SecurityPolicyVersion) == "" ||
		len(config.SecurityPolicyVersion) > 128 ||
		config.SecurityPolicyVersion != strings.TrimSpace(config.SecurityPolicyVersion) {
		return nil, errors.New("security policy version is required and must be normalized")
	}
	if !config.RequireMFA {
		return nil, errors.New("authenticated sessions require MFA assurance enforcement")
	}
	location, err := safeLocalLocation(config.SuccessLocation)
	if err != nil {
		return nil, err
	}

	return &Service{
		store:                 config.Store,
		resolver:              config.Resolver,
		random:                config.Random,
		now:                   config.Now,
		idleLifetime:          config.IdleLifetime,
		absoluteLifetime:      config.AbsoluteLifetime,
		successLocation:       location,
		securityPolicyVersion: config.SecurityPolicyVersion,
		requireMFA:            config.RequireMFA,
	}, nil
}

func (s *Service) ServeVerifiedPrincipal(
	writer http.ResponseWriter,
	request *http.Request,
	principal authentication.Principal,
) {
	browserNoStore(writer)
	if s == nil || request == nil {
		writeFailure(writer, http.StatusServiceUnavailable, "authentication service unavailable")
		return
	}
	if err := principal.Validate(); err != nil || principal.BoundActorID != "" {
		writeFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}

	now := s.now().UTC()
	if principal.AuthenticatedAt.After(now.Add(2*time.Minute)) ||
		principal.Assurance.SecurityPolicyVersion != s.securityPolicyVersion ||
		(s.requireMFA && (!principal.Assurance.MFAAuthenticated ||
			principal.Assurance.MFAAuthenticatedAt.IsZero() ||
			principal.Assurance.MFAAuthenticatedAt.After(now.Add(2*time.Minute)))) {
		writeFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}

	actor, err := s.resolver.Resolve(request.Context(), principal)
	if err != nil {
		writeAuthenticationError(writer, err)
		return
	}

	identifier, digest, err := generateIdentifier(s.random)
	if err != nil {
		writeFailure(writer, http.StatusServiceUnavailable, "authentication service unavailable")
		return
	}
	requestRecord := CreateRequest{
		IdentifierDigest:      digest,
		Principal:             principal,
		ActorID:               actor.ID,
		CreatedAt:             now,
		IdleExpiresAt:         now.Add(s.idleLifetime),
		AbsoluteExpiresAt:     now.Add(s.absoluteLifetime),
		SecurityPolicyVersion: s.securityPolicyVersion,
	}
	created, err := s.store.Create(request.Context(), requestRecord)
	if err != nil {
		writeStoreError(writer, err)
		return
	}
	if err := created.Validate(now); err != nil {
		writeFailure(writer, http.StatusServiceUnavailable, "authentication service unavailable")
		return
	}

	http.SetCookie(writer, sessionCookie(identifier, created.AbsoluteExpiresAt, now))
	writer.Header().Set("Location", s.successLocation)
	writer.WriteHeader(http.StatusSeeOther)
}

func (s *Service) Authenticate(
	ctx context.Context,
	request *http.Request,
) (authentication.Principal, error) {
	if s == nil || s.store == nil || request == nil {
		return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
	}
	cookies := request.CookiesNamed(CookieName)
	if len(cookies) != 1 {
		return authentication.Principal{}, authentication.ErrAuthenticationRequired
	}
	identifier, err := decodeIdentifier(cookies[0].Value)
	if err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationRequired
	}
	digest := sha256.Sum256(identifier)
	now := s.now().UTC()
	record, err := s.store.Find(ctx, digest)
	if err != nil {
		if errors.Is(err, ErrSessionUnavailable) {
			return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
		}
		return authentication.Principal{}, authentication.ErrAuthenticationRequired
	}
	if err := record.Validate(now); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationRequired
	}
	principal := record.Principal
	principal.BoundActorID = record.ActorID
	principal.Assurance.SecurityPolicyVersion = record.SecurityPolicyVersion
	if err := principal.Validate(); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationRequired
	}
	return principal, nil
}

func (r Record) Validate(now time.Time) error {
	if now.IsZero() {
		return errors.New("session validation time is required")
	}
	if err := r.Principal.Validate(); err != nil {
		return fmt.Errorf("principal: %w", err)
	}
	if strings.TrimSpace(r.ActorID) == "" || r.ActorID != strings.TrimSpace(r.ActorID) || len(r.ActorID) > 256 {
		return errors.New("session actor identifier is invalid")
	}
	if r.CreatedAt.IsZero() || r.LastActivityAt.IsZero() ||
		r.IdleExpiresAt.IsZero() || r.AbsoluteExpiresAt.IsZero() {
		return errors.New("session timestamps are incomplete")
	}
	if r.LastActivityAt.Before(r.CreatedAt) ||
		!r.IdleExpiresAt.After(r.LastActivityAt) ||
		!r.AbsoluteExpiresAt.After(r.CreatedAt) ||
		r.IdleExpiresAt.After(r.AbsoluteExpiresAt) {
		return errors.New("session timestamps are inconsistent")
	}
	if !r.RevokedAt.IsZero() {
		return errors.New("session is revoked")
	}
	if !now.Before(r.IdleExpiresAt) || !now.Before(r.AbsoluteExpiresAt) {
		return errors.New("session is expired")
	}
	if r.SecurityPolicyVersion == "" ||
		r.SecurityPolicyVersion != strings.TrimSpace(r.SecurityPolicyVersion) ||
		len(r.SecurityPolicyVersion) > 128 {
		return errors.New("session security policy version is invalid")
	}
	if !r.Principal.Assurance.MFAAuthenticated ||
		r.Principal.Assurance.MFAAuthenticatedAt.IsZero() {
		return errors.New("session MFA assurance is required")
	}
	if r.Principal.Assurance.SecurityPolicyVersion != r.SecurityPolicyVersion {
		return errors.New("session security policy version is inconsistent")
	}
	return nil
}

func generateIdentifier(reader io.Reader) (string, [sha256.Size]byte, error) {
	var raw [identifierBytes]byte
	if _, err := io.ReadFull(reader, raw[:]); err != nil {
		return "", [sha256.Size]byte{}, fmt.Errorf("generate session identifier: %w", err)
	}
	identifier := base64.RawURLEncoding.EncodeToString(raw[:])
	return identifier, sha256.Sum256(raw[:]), nil
}

func decodeIdentifier(value string) ([]byte, error) {
	if len(value) != encodedIdentifierBytes || value != strings.TrimSpace(value) {
		return nil, ErrSessionInvalid
	}
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil || len(raw) != identifierBytes {
		return nil, ErrSessionInvalid
	}
	if base64.RawURLEncoding.EncodeToString(raw) != value {
		return nil, ErrSessionInvalid
	}
	return raw, nil
}

func sessionCookie(identifier string, absoluteExpiry time.Time, now time.Time) *http.Cookie {
	remaining := absoluteExpiry.Sub(now)
	maxAge := int(remaining / time.Second)
	if remaining%time.Second != 0 {
		maxAge++
	}
	if maxAge < 1 {
		maxAge = 1
	}
	return &http.Cookie{
		Name:     CookieName,
		Value:    identifier,
		Path:     "/",
		Expires:  absoluteExpiry.UTC(),
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func safeLocalLocation(raw string) (string, error) {
	if raw == "" {
		return "/", nil
	}
	if raw != strings.TrimSpace(raw) || len(raw) > 2048 ||
		!strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") {
		return "", errors.New("session success location must be an absolute local path")
	}
	for _, character := range raw {
		if !safeLocalPathCharacter(character) {
			return "", errors.New("session success location contains an unsafe character")
		}
	}
	return raw, nil
}

func safeLocalPathCharacter(character rune) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' ||
		character == '/' || character == '-' || character == '_' ||
		character == '.' || character == '~'
}

func browserNoStore(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Referrer-Policy", "no-referrer")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
}

func writeAuthenticationError(writer http.ResponseWriter, err error) {
	if errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		writeFailure(writer, http.StatusServiceUnavailable, "authentication service unavailable")
		return
	}
	writeFailure(writer, http.StatusUnauthorized, "authentication failed")
}

func writeStoreError(writer http.ResponseWriter, err error) {
	if errors.Is(err, ErrSessionUnavailable) || errors.Is(err, ErrSessionConflict) {
		writeFailure(writer, http.StatusServiceUnavailable, "authentication service unavailable")
		return
	}
	writeFailure(writer, http.StatusUnauthorized, "authentication failed")
}

func writeFailure(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(message + "\n"))
}
