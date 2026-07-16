package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

type Mode string

const (
	ModeDevelopment Mode = "development"
	ModeProduction  Mode = "production"

	DevelopmentActorHeader = "X-Iron-Atlas-Actor"
	DevelopmentRolesHeader = "X-Iron-Atlas-Roles"
)

var (
	ErrAuthenticationRequired    = errors.New("authentication required")
	ErrAuthenticationInvalid     = errors.New("authentication result is invalid")
	ErrAuthenticationUnavailable = errors.New("authentication service unavailable")
	ErrIdentityResolutionFailed  = errors.New("actor resolution failed")
)

type Principal struct {
	ProviderID      string
	Subject         string
	AuthenticatedAt time.Time
}

type ResolvedIdentity struct {
	Principal Principal
	Actor     authz.Actor
}

type Authenticator interface {
	Authenticate(context.Context, *http.Request) (Principal, error)
}

type AuthenticatorFunc func(context.Context, *http.Request) (Principal, error)

func (f AuthenticatorFunc) Authenticate(
	ctx context.Context,
	request *http.Request,
) (Principal, error) {
	return f(ctx, request)
}

type ActorResolver interface {
	Resolve(context.Context, Principal) (authz.Actor, error)
}

type ActorResolverFunc func(context.Context, Principal) (authz.Actor, error)

func (f ActorResolverFunc) Resolve(
	ctx context.Context,
	principal Principal,
) (authz.Actor, error) {
	return f(ctx, principal)
}

type Options struct {
	Mode             Mode
	Authenticator    Authenticator
	ActorResolver    ActorResolver
	DevelopmentActor authz.Actor
	Now              func() time.Time
}

type Middleware struct {
	mode             Mode
	authenticator    Authenticator
	actorResolver    ActorResolver
	developmentActor authz.Actor
	now              func() time.Time
}

type identityContextKey struct{}

func ParseMode(raw string) (Mode, error) {
	mode := Mode(strings.ToLower(strings.TrimSpace(raw)))
	if err := mode.Validate(); err != nil {
		return "", err
	}
	return mode, nil
}

func (m Mode) Validate() error {
	switch m {
	case ModeDevelopment, ModeProduction:
		return nil
	default:
		return fmt.Errorf(
			"authentication mode must be %q or %q",
			ModeDevelopment,
			ModeProduction,
		)
	}
}

func New(options Options) (*Middleware, error) {
	if err := options.Mode.Validate(); err != nil {
		return nil, err
	}
	if (options.Authenticator == nil) != (options.ActorResolver == nil) {
		return nil, errors.New(
			"authenticator and actor resolver must be configured together",
		)
	}
	if options.Mode == ModeDevelopment &&
		(options.Authenticator != nil || options.ActorResolver != nil) {
		return nil, errors.New(
			"production authentication components are prohibited in development mode",
		)
	}
	if options.Mode == ModeProduction &&
		(options.DevelopmentActor.ID != "" ||
			len(options.DevelopmentActor.Roles) != 0) {
		return nil, errors.New(
			"development actor defaults are prohibited in production mode",
		)
	}

	now := options.Now
	if now == nil {
		now = time.Now
	}

	developmentActor := cloneActor(options.DevelopmentActor)
	if options.Mode == ModeDevelopment && developmentActor.ID == "" {
		developmentActor = authz.Actor{
			ID:    "network-tech-01",
			Roles: []authz.Role{authz.RoleNetworkTech},
		}
	}
	if options.Mode == ModeDevelopment {
		if err := validateActor(developmentActor); err != nil {
			return nil, fmt.Errorf("development actor: %w", err)
		}
	}

	return &Middleware{
		mode:             options.Mode,
		authenticator:    options.Authenticator,
		actorResolver:    options.ActorResolver,
		developmentActor: developmentActor,
		now:              now,
	}, nil
}

func (m *Middleware) Mode() Mode {
	if m == nil {
		return ""
	}
	return m.mode
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	if next == nil {
		panic("authentication middleware requires a next handler")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if _, ok := identityFromContext(r.Context()); ok {
			writeFailure(
				w,
				http.StatusInternalServerError,
				"request identity boundary violation",
			)
			return
		}

		var (
			identity ResolvedIdentity
			err      error
		)
		switch m.mode {
		case ModeDevelopment:
			identity, err = m.developmentIdentity(r)
			if err != nil {
				writeFailure(w, http.StatusBadRequest, "invalid development identity")
				return
			}
		case ModeProduction:
			if hasDevelopmentHeaders(r.Header) {
				writeFailure(
					w,
					http.StatusBadRequest,
					"development identity headers are prohibited",
				)
				return
			}
			identity, err = m.productionIdentity(r)
			if err != nil {
				status := http.StatusUnauthorized
				message := "authentication required"
				if errors.Is(err, ErrAuthenticationUnavailable) {
					status = http.StatusServiceUnavailable
					message = "authentication service unavailable"
				}
				writeFailure(w, status, message)
				return
			}
		default:
			writeFailure(
				w,
				http.StatusInternalServerError,
				"authentication mode is invalid",
			)
			return
		}

		if err := identity.Validate(); err != nil {
			writeFailure(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(
			w,
			r.WithContext(withIdentity(r.Context(), identity)),
		)
	})
}

func (m *Middleware) developmentIdentity(
	request *http.Request,
) (ResolvedIdentity, error) {
	actorID, err := singleHeader(request.Header, DevelopmentActorHeader)
	if err != nil {
		return ResolvedIdentity{}, err
	}
	rawRoles, err := singleHeader(request.Header, DevelopmentRolesHeader)
	if err != nil {
		return ResolvedIdentity{}, err
	}

	actor := cloneActor(m.developmentActor)
	if actorID != "" {
		actor.ID = actorID
	}
	if rawRoles != "" {
		actor.Roles, err = parseDevelopmentRoles(rawRoles)
		if err != nil {
			return ResolvedIdentity{}, err
		}
	}
	if err := validateActor(actor); err != nil {
		return ResolvedIdentity{}, err
	}

	return ResolvedIdentity{
		Principal: Principal{
			ProviderID:      "development",
			Subject:         actor.ID,
			AuthenticatedAt: m.now().UTC(),
		},
		Actor: actor,
	}, nil
}

func (m *Middleware) productionIdentity(
	request *http.Request,
) (ResolvedIdentity, error) {
	if m.authenticator == nil || m.actorResolver == nil {
		return ResolvedIdentity{}, ErrAuthenticationRequired
	}

	principal, err := m.authenticator.Authenticate(
		request.Context(),
		request,
	)
	if err != nil {
		if errors.Is(err, ErrAuthenticationUnavailable) {
			return ResolvedIdentity{}, ErrAuthenticationUnavailable
		}
		return ResolvedIdentity{}, ErrAuthenticationRequired
	}
	if err := principal.Validate(); err != nil {
		return ResolvedIdentity{}, ErrAuthenticationInvalid
	}

	actor, err := m.actorResolver.Resolve(request.Context(), principal)
	if err != nil {
		if errors.Is(err, ErrAuthenticationUnavailable) {
			return ResolvedIdentity{}, ErrAuthenticationUnavailable
		}
		return ResolvedIdentity{}, ErrIdentityResolutionFailed
	}
	identity := ResolvedIdentity{Principal: principal, Actor: actor}
	if err := identity.Validate(); err != nil {
		return ResolvedIdentity{}, ErrIdentityResolutionFailed
	}
	return identity, nil
}

func (p Principal) Validate() error {
	if err := validateBoundedIdentifier("provider ID", p.ProviderID, 256); err != nil {
		return err
	}
	if err := validateBoundedIdentifier("provider subject", p.Subject, 512); err != nil {
		return err
	}
	if p.AuthenticatedAt.IsZero() {
		return errors.New("authentication time is required")
	}
	return nil
}

func (i ResolvedIdentity) Validate() error {
	if err := i.Principal.Validate(); err != nil {
		return err
	}
	if err := validateActor(i.Actor); err != nil {
		return err
	}
	return nil
}

func ResolvedIdentityFromContext(
	ctx context.Context,
) (ResolvedIdentity, bool) {
	identity, ok := identityFromContext(ctx)
	if !ok {
		return ResolvedIdentity{}, false
	}
	identity.Actor = cloneActor(identity.Actor)
	return identity, true
}

func ActorFromContext(ctx context.Context) (authz.Actor, bool) {
	identity, ok := ResolvedIdentityFromContext(ctx)
	if !ok {
		return authz.Actor{}, false
	}
	return identity.Actor, true
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	identity, ok := ResolvedIdentityFromContext(ctx)
	if !ok {
		return Principal{}, false
	}
	return identity.Principal, true
}

func withIdentity(
	ctx context.Context,
	identity ResolvedIdentity,
) context.Context {
	identity.Actor = cloneActor(identity.Actor)
	return context.WithValue(ctx, identityContextKey{}, identity)
}

func identityFromContext(
	ctx context.Context,
) (ResolvedIdentity, bool) {
	identity, ok := ctx.Value(identityContextKey{}).(ResolvedIdentity)
	return identity, ok
}

func publicPath(path string) bool {
	return path == "/healthz" ||
		path == "/readyz" ||
		strings.HasPrefix(path, "/static/")
}

func hasDevelopmentHeaders(header http.Header) bool {
	return len(header.Values(DevelopmentActorHeader)) != 0 ||
		len(header.Values(DevelopmentRolesHeader)) != 0
}

func singleHeader(header http.Header, name string) (string, error) {
	values := header.Values(name)
	if len(values) > 1 {
		return "", fmt.Errorf("%s must not be repeated", name)
	}
	if len(values) == 0 {
		return "", nil
	}
	value := strings.TrimSpace(values[0])
	if len(value) > 1024 {
		return "", fmt.Errorf("%s is too large", name)
	}
	return value, nil
}

func parseDevelopmentRoles(raw string) ([]authz.Role, error) {
	parts := strings.Split(raw, ",")
	if len(parts) > 32 {
		return nil, errors.New("too many development roles")
	}
	roles := make([]authz.Role, 0, len(parts))
	seen := make(map[authz.Role]struct{}, len(parts))
	for _, part := range parts {
		role := authz.Role(strings.TrimSpace(part))
		if role == "" {
			continue
		}
		if !knownRole(role) {
			return nil, fmt.Errorf("unknown development role %q", role)
		}
		if _, exists := seen[role]; exists {
			return nil, fmt.Errorf("duplicate development role %q", role)
		}
		seen[role] = struct{}{}
		roles = append(roles, role)
	}
	if len(roles) == 0 {
		return nil, errors.New("at least one development role is required")
	}
	return roles, nil
}

func knownRole(role authz.Role) bool {
	switch role {
	case authz.RoleViewer,
		authz.RoleNetworkTech,
		authz.RoleNetworkAdmin,
		authz.RoleNetworkSecurity,
		authz.RoleChangeApprover,
		authz.RoleAuditor,
		authz.RolePlatformAdmin:
		return true
	default:
		return false
	}
}

func validateActor(actor authz.Actor) error {
	if err := validateBoundedIdentifier("actor ID", actor.ID, 256); err != nil {
		return err
	}
	if len(actor.Roles) > 32 {
		return errors.New("actor has too many roles")
	}
	for _, role := range actor.Roles {
		if !knownRole(role) {
			return fmt.Errorf("actor contains unknown role %q", role)
		}
	}
	return nil
}

func validateBoundedIdentifier(
	name string,
	value string,
	maxBytes int,
) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	if !utf8.ValidString(value) {
		return fmt.Errorf("%s must be valid UTF-8", name)
	}
	if len(value) > maxBytes {
		return fmt.Errorf("%s exceeds %d bytes", name, maxBytes)
	}
	if strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return fmt.Errorf("%s contains a control character", name)
	}
	return nil
}

func cloneActor(actor authz.Actor) authz.Actor {
	cloned := actor
	cloned.Roles = append([]authz.Role(nil), actor.Roles...)
	return cloned
}

func writeFailure(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
