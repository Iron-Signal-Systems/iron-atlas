package common

type InterfaceClass string

const (
	ClassAccessEndpoint      InterfaceClass = "access_endpoint"
	ClassInfrastructureTrunk InterfaceClass = "infrastructure_trunk"
	ClassPortChannelMember   InterfaceClass = "port_channel_member"
	ClassPortChannel         InterfaceClass = "port_channel"
	ClassRouted              InterfaceClass = "routed_interface"
	ClassStack               InterfaceClass = "stack_interface"
	ClassFabric              InterfaceClass = "fabric_interconnect"
	ClassExcluded            InterfaceClass = "explicitly_excluded"
	ClassUnknown             InterfaceClass = "unknown"
)

type Interface struct {
	DeviceID           string
	Name               string
	Description        string
	AdministrativeMode string
	OperationalMode    string
	Routed             bool
	PortChannelID      string
	StackInterface     bool
	FabricInterconnect bool
	ExplicitlyExcluded bool
	NativeVLAN         string
	AllowedVLANs       []string
	PrunedVLANs        []string
	CDPNeighbor        *Neighbor
	LLDPNeighbor       *Neighbor
	SpanningTree       []SpanningTreeState
	ACLs               []ACLAttachment
	MACAddresses       []string
}

type Neighbor struct{ DeviceID, LocalInterface, RemoteInterface, Platform, ManagementAddress string }
type SpanningTreeState struct {
	VLAN, Role, State string
	Cost              int
}
type ACLAttachment struct{ Name, Direction, AddressFamily, Scope string }

type Classification struct {
	Class               InterfaceClass
	EndpointAttribution bool
	Reasons             []string
}
