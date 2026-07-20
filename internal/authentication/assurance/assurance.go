package assurance

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

const (
	defaultMaximumAuthenticationAge = 15 * time.Minute
	defaultMaximumClockSkew         = time.Minute
	maximumConfiguredAge            = 24 * time.Hour
	maximumPolicyValues             = 32
	maximumMethodSetSize            = 16
)

type Outcome string

const (
	OutcomeSatisfied                        Outcome = "satisfied"
	OutcomeAdditionalAuthenticationRequired Outcome = "additional_authentication_required"
	OutcomeStepUpRequired                   Outcome = "step_up_required"
	OutcomePhishingResistantRequired        Outcome = "phishing_resistant_required"
	OutcomeDenied                           Outcome = "denied"
)

const (
	ReasonSatisfied                  = "assurance_satisfied"
	ReasonMFARequired                = "mfa_required"
	ReasonPrimaryAuthenticationStale = "primary_authentication_stale"
	ReasonPhishingResistantRequired  = "phishing_resistant_required"
	ReasonInvalidAssurance           = "invalid_assurance"
)

// VerifiedPrincipalHandler is the narrow downstream seam used after assurance
// evaluation. The authenticated-session service implements this contract.
type VerifiedPrincipalHandler interface {
	ServeVerifiedPrincipal(http.ResponseWriter, *http.Request, authentication.Principal)
}

type MethodSet []string

type PolicyConfig struct {
	Version                             string
	RequireMFA                          bool
	MaximumAuthenticationAge            time.Duration
	MaximumClockSkew                    time.Duration
	AcceptedMFAContexts                 []string
	AcceptedMFAMethodSets               []MethodSet
	PhishingResistantRoles              []authz.Role
	AcceptedPhishingResistantContexts   []string
	AcceptedPhishingResistantMethodSets []MethodSet
}

type Policy struct {
	version                             string
	requireMFA                          bool
	maximumAuthenticationAge            time.Duration
	maximumClockSkew                    time.Duration
	acceptedMFAContexts                 map[string]struct{}
	acceptedMFAMethodSets               []map[string]struct{}
	phishingResistantRoles              map[authz.Role]struct{}
	acceptedPhishingResistantContexts   map[string]struct{}
	acceptedPhishingResistantMethodSets []map[string]struct{}
}

type Decision struct {
	Outcome       Outcome
	ReasonCode    string
	PolicyVersion string
	Assurance     authentication.Assurance
}

func (d Decision) Satisfied() bool {
	return d.Outcome == OutcomeSatisfied
}

func NewPolicy(config PolicyConfig) (*Policy, error) {
	version, err := normalizedIdentifier("security policy version", config.Version, 128)
	if err != nil {
		return nil, err
	}
	if !config.RequireMFA {
		return nil, errors.New("authentication assurance policy must require MFA")
	}
	maximumAge := config.MaximumAuthenticationAge
	if maximumAge == 0 {
		maximumAge = defaultMaximumAuthenticationAge
	}
	if maximumAge < time.Minute || maximumAge > maximumConfiguredAge {
		return nil, fmt.Errorf(
			"maximum authentication age must be between one minute and %s",
			maximumConfiguredAge,
		)
	}
	maximumClockSkew := config.MaximumClockSkew
	if maximumClockSkew == 0 {
		maximumClockSkew = defaultMaximumClockSkew
	}
	if maximumClockSkew < 0 || maximumClockSkew > 5*time.Minute {
		return nil, errors.New("maximum clock skew must be between zero and five minutes")
	}

	mfaContexts, err := normalizedSet("accepted MFA context", config.AcceptedMFAContexts, 256)
	if err != nil {
		return nil, err
	}
	mfaMethods, err := normalizedMethodSets("accepted MFA method set", config.AcceptedMFAMethodSets)
	if err != nil {
		return nil, err
	}
	phishingContexts, err := normalizedSet(
		"accepted phishing-resistant context",
		config.AcceptedPhishingResistantContexts,
		256,
	)
	if err != nil {
		return nil, err
	}
	phishingMethods, err := normalizedMethodSets(
		"accepted phishing-resistant method set",
		config.AcceptedPhishingResistantMethodSets,
	)
	if err != nil {
		return nil, err
	}
	roles, err := normalizedRoles(config.PhishingResistantRoles)
	if err != nil {
		return nil, err
	}

	return &Policy{
		version:                             version,
		requireMFA:                          config.RequireMFA,
		maximumAuthenticationAge:            maximumAge,
		maximumClockSkew:                    maximumClockSkew,
		acceptedMFAContexts:                 mfaContexts,
		acceptedMFAMethodSets:               mfaMethods,
		phishingResistantRoles:              roles,
		acceptedPhishingResistantContexts:   phishingContexts,
		acceptedPhishingResistantMethodSets: phishingMethods,
	}, nil
}

func (p *Policy) Version() string {
	if p == nil {
		return ""
	}
	return p.version
}

func (p *Policy) Evaluate(
	principal authentication.Principal,
	actor authz.Actor,
	now time.Time,
) Decision {
	denied := Decision{Outcome: OutcomeDenied, ReasonCode: ReasonInvalidAssurance}
	if p == nil || now.IsZero() || principal.Validate() != nil || !validActor(actor) {
		return denied
	}
	now = now.UTC()
	if principal.AuthenticatedAt.After(now.Add(p.maximumClockSkew)) {
		return denied
	}
	if now.Sub(principal.AuthenticatedAt) > p.maximumAuthenticationAge+p.maximumClockSkew {
		return Decision{
			Outcome:       OutcomeStepUpRequired,
			ReasonCode:    ReasonPrimaryAuthenticationStale,
			PolicyVersion: p.version,
		}
	}

	methods := make(map[string]struct{}, len(principal.Assurance.Methods))
	for _, method := range principal.Assurance.Methods {
		methods[method] = struct{}{}
	}
	phishingResistant := contextAccepted(
		principal.Assurance.Context,
		p.acceptedPhishingResistantContexts,
	) || methodSetAccepted(methods, p.acceptedPhishingResistantMethodSets)
	mfaAccepted := phishingResistant || contextAccepted(
		principal.Assurance.Context,
		p.acceptedMFAContexts,
	) || methodSetAccepted(methods, p.acceptedMFAMethodSets)

	if p.actorRequiresPhishingResistance(actor) && !phishingResistant {
		return Decision{
			Outcome:       OutcomePhishingResistantRequired,
			ReasonCode:    ReasonPhishingResistantRequired,
			PolicyVersion: p.version,
		}
	}
	if p.requireMFA && !mfaAccepted {
		return Decision{
			Outcome:       OutcomeAdditionalAuthenticationRequired,
			ReasonCode:    ReasonMFARequired,
			PolicyVersion: p.version,
		}
	}

	assurance := principal.Assurance
	assurance.Methods = append([]string(nil), principal.Assurance.Methods...)
	assurance.MFAAuthenticated = mfaAccepted
	if mfaAccepted {
		assurance.MFAAuthenticatedAt = principal.AuthenticatedAt.UTC()
	} else {
		assurance.MFAAuthenticatedAt = time.Time{}
	}
	assurance.SecurityPolicyVersion = p.version
	if assurance.Validate() != nil {
		return denied
	}
	return Decision{
		Outcome:       OutcomeSatisfied,
		ReasonCode:    ReasonSatisfied,
		PolicyVersion: p.version,
		Assurance:     assurance,
	}
}

func (p *Policy) actorRequiresPhishingResistance(actor authz.Actor) bool {
	for _, role := range actor.Roles {
		if _, ok := p.phishingResistantRoles[role]; ok {
			return true
		}
	}
	return false
}

type ServiceConfig struct {
	Resolver authentication.ActorResolver
	Policy   *Policy
	Next     VerifiedPrincipalHandler
	Now      func() time.Time
}

type Service struct {
	resolver authentication.ActorResolver
	policy   *Policy
	next     VerifiedPrincipalHandler
	now      func() time.Time
}

func NewService(config ServiceConfig) (*Service, error) {
	if config.Resolver == nil {
		return nil, errors.New("governed actor resolver is required")
	}
	if config.Policy == nil {
		return nil, errors.New("authentication assurance policy is required")
	}
	if config.Next == nil {
		return nil, errors.New("verified principal downstream handler is required")
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	return &Service{
		resolver: config.Resolver,
		policy:   config.Policy,
		next:     config.Next,
		now:      now,
	}, nil
}

func (s *Service) ServeVerifiedPrincipal(
	writer http.ResponseWriter,
	request *http.Request,
	principal authentication.Principal,
) {
	browserNoStore(writer)
	if s == nil || request == nil || principal.Validate() != nil || principal.BoundActorID != "" {
		writeFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}
	actor, err := s.resolver.Resolve(request.Context(), principal)
	if err != nil {
		writeAuthenticationError(writer, err)
		return
	}
	decision := s.policy.Evaluate(principal, actor, s.now().UTC())
	if !decision.Satisfied() {
		switch decision.Outcome {
		case OutcomeAdditionalAuthenticationRequired,
			OutcomeStepUpRequired,
			OutcomePhishingResistantRequired:
			writeFailure(writer, http.StatusUnauthorized, "additional authentication required")
		default:
			writeFailure(writer, http.StatusUnauthorized, "authentication failed")
		}
		return
	}
	principal.Assurance = decision.Assurance
	s.next.ServeVerifiedPrincipal(writer, request, principal)
}

func normalizedSet(name string, values []string, maximum int) (map[string]struct{}, error) {
	if len(values) > maximumPolicyValues {
		return nil, fmt.Errorf("%s has too many values", name)
	}
	result := make(map[string]struct{}, len(values))
	for _, raw := range values {
		value, err := normalizedIdentifier(name, raw, maximum)
		if err != nil {
			return nil, err
		}
		if _, duplicate := result[value]; duplicate {
			return nil, fmt.Errorf("%s contains duplicate value %q", name, value)
		}
		result[value] = struct{}{}
	}
	return result, nil
}

func normalizedMethodSets(name string, sets []MethodSet) ([]map[string]struct{}, error) {
	if len(sets) > maximumPolicyValues {
		return nil, fmt.Errorf("%s has too many sets", name)
	}
	result := make([]map[string]struct{}, 0, len(sets))
	seenSets := make(map[string]struct{}, len(sets))
	for _, rawSet := range sets {
		if len(rawSet) == 0 || len(rawSet) > maximumMethodSetSize {
			return nil, fmt.Errorf("%s must contain between one and %d methods", name, maximumMethodSetSize)
		}
		set := make(map[string]struct{}, len(rawSet))
		ordered := make([]string, 0, len(rawSet))
		for _, raw := range rawSet {
			value, err := normalizedIdentifier(name, raw, 64)
			if err != nil {
				return nil, err
			}
			if _, duplicate := set[value]; duplicate {
				return nil, fmt.Errorf("%s contains duplicate method %q", name, value)
			}
			set[value] = struct{}{}
			ordered = append(ordered, value)
		}
		sort.Strings(ordered)
		key := strings.Join(ordered, "\x00")
		if _, duplicate := seenSets[key]; duplicate {
			return nil, fmt.Errorf("%s contains a duplicate method set", name)
		}
		seenSets[key] = struct{}{}
		result = append(result, set)
	}
	return result, nil
}

func normalizedRoles(values []authz.Role) (map[authz.Role]struct{}, error) {
	if len(values) > maximumPolicyValues {
		return nil, errors.New("phishing-resistant role list is too large")
	}
	result := make(map[authz.Role]struct{}, len(values))
	for _, role := range values {
		if !knownRole(role) {
			return nil, fmt.Errorf("unknown phishing-resistant role %q", role)
		}
		if _, duplicate := result[role]; duplicate {
			return nil, fmt.Errorf("duplicate phishing-resistant role %q", role)
		}
		result[role] = struct{}{}
	}
	return result, nil
}

func validActor(actor authz.Actor) bool {
	if actor.ID == "" || actor.ID != strings.TrimSpace(actor.ID) || len(actor.ID) > 256 {
		return false
	}
	seen := make(map[authz.Role]struct{}, len(actor.Roles))
	for _, role := range actor.Roles {
		if !knownRole(role) {
			return false
		}
		if _, duplicate := seen[role]; duplicate {
			return false
		}
		seen[role] = struct{}{}
	}
	return true
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

func normalizedIdentifier(name string, value string, maximum int) (string, error) {
	if value == "" || value != strings.TrimSpace(value) || len(value) > maximum || !utf8.ValidString(value) {
		return "", fmt.Errorf("%s is missing, unnormalized, invalid, or too large", name)
	}
	for _, character := range value {
		if unicode.IsControl(character) {
			return "", fmt.Errorf("%s contains a control character", name)
		}
	}
	return value, nil
}

func contextAccepted(value string, accepted map[string]struct{}) bool {
	if value == "" {
		return false
	}
	_, ok := accepted[value]
	return ok
}

func methodSetAccepted(methods map[string]struct{}, accepted []map[string]struct{}) bool {
	for _, required := range accepted {
		matched := true
		for method := range required {
			if _, ok := methods[method]; !ok {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
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

func writeFailure(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(message + "\n"))
}
