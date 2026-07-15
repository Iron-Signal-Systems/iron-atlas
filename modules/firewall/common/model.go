package common

type EvidenceState string

const (
	StateConfigured  EvidenceState = "configured"
	StateObserved    EvidenceState = "observed"
	StateCalculated  EvidenceState = "calculated"
	StateInferred    EvidenceState = "inferred"
	StateUnknown     EvidenceState = "unknown"
	StateConflicting EvidenceState = "conflicting"
)

type Firewall struct {
	Vendor       string
	Platform     string
	Version      string
	Interfaces   []Interface
	StaticRoutes []StaticRoute
	SDWAN        []SDWANZone
	Policies     []Policy
	Warnings     []string
}

type Interface struct {
	ID        string
	Name      string
	Addresses []string
	Zone      string
	Role      string
	Enabled   bool
}

type StaticRoute struct {
	ID          string
	Destination string
	Gateway     string
	Interface   string
	Distance    int
	Priority    int
	State       EvidenceState
}

type SDWANZone struct {
	ID      string
	Name    string
	Members []string
}

type Policy struct {
	ID           string
	Name         string
	Sequence     int
	Ingress      []string
	Egress       []string
	Sources      []string
	Destinations []string
	Services     []string
	Action       string
	NAT          bool
	Enabled      bool
}
