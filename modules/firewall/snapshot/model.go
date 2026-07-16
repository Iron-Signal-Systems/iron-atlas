package snapshot

import "time"

type SnapshotID string
type ObjectID string

type EvidenceState string

const (
	EvidenceConfigured  EvidenceState = "configured"
	EvidenceObserved    EvidenceState = "observed"
	EvidenceCalculated  EvidenceState = "calculated"
	EvidenceInferred    EvidenceState = "inferred"
	EvidenceUnknown     EvidenceState = "unknown"
	EvidenceConflicting EvidenceState = "conflicting"
)

type FirewallSnapshot struct {
	ID       SnapshotID     `json:"id,omitempty"`
	Source   SnapshotSource `json:"source"`
	Device   Device         `json:"device"`
	Captured time.Time      `json:"captured,omitempty"`

	Domains    []RoutingDomain `json:"domains,omitempty"`
	Interfaces []Interface     `json:"interfaces,omitempty"`
	VLANs      []VLAN          `json:"vlans,omitempty"`
	Subnets    []Subnet        `json:"subnets,omitempty"`

	DHCP    DHCPConfiguration    `json:"dhcp"`
	Routing RoutingConfiguration `json:"routing"`
	SDWAN   SDWANConfiguration   `json:"sdwan"`

	Policies []Policy         `json:"policies,omitempty"`
	Objects  ObjectCatalog    `json:"objects"`
	NAT      NATConfiguration `json:"nat"`
	VPNs     []VPN            `json:"vpns,omitempty"`
	QoS      QoSConfiguration `json:"qos"`

	Tags       []Tag          `json:"tags,omitempty"`
	References ReferenceGraph `json:"references"`
	Findings   []Finding      `json:"findings,omitempty"`
	Extensions []Extension    `json:"extensions,omitempty"`
}

type SnapshotSource struct {
	Format         string `json:"format"`
	Vendor         string `json:"vendor"`
	Platform       string `json:"platform"`
	FortiOSVersion string `json:"fortios_version,omitempty"`
	Filename       string `json:"filename,omitempty"`
	SHA256         string `json:"sha256,omitempty"`
	PasswordMasked bool   `json:"password_masked,omitempty"`
}

type Device struct {
	Hostname     string `json:"hostname,omitempty"`
	Alias        string `json:"alias,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	Model        string `json:"model,omitempty"`
	Version      string `json:"version,omitempty"`
	VDOMMode     string `json:"vdom_mode,omitempty"`
}

type ObjectScope struct {
	VDOM string  `json:"vdom,omitempty"`
	VRF  *uint32 `json:"vrf,omitempty"`
}

type SourceLocation struct {
	YAMLPath string `json:"yaml_path,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
}

type ObjectKind string

const (
	ObjectKindInterface       ObjectKind = "interface"
	ObjectKindZone            ObjectKind = "zone"
	ObjectKindInterfaceOrZone ObjectKind = "interface-or-zone"
	ObjectKindVLAN            ObjectKind = "vlan"
	ObjectKindSubnet          ObjectKind = "subnet"
	ObjectKindAddress         ObjectKind = "address"
	ObjectKindAddressGroup    ObjectKind = "address-group"
	ObjectKindService         ObjectKind = "service"
	ObjectKindServiceGroup    ObjectKind = "service-group"
	ObjectKindSchedule        ObjectKind = "schedule"
	ObjectKindUser            ObjectKind = "user"
	ObjectKindUserGroup       ObjectKind = "user-group"
	ObjectKindVIP             ObjectKind = "virtual-ip"
	ObjectKindIPPool          ObjectKind = "ip-pool"
	ObjectKindTrafficShaper   ObjectKind = "traffic-shaper"
	ObjectKindPerIPShaper     ObjectKind = "per-ip-shaper"
	ObjectKindSDWANMember     ObjectKind = "sdwan-member"
	ObjectKindSDWANZone       ObjectKind = "sdwan-zone"
	ObjectKindSDWANHealth     ObjectKind = "sdwan-health-check"
	ObjectKindRouteMap        ObjectKind = "route-map"
	ObjectKindPrefixList      ObjectKind = "prefix-list"
	ObjectKindVPN             ObjectKind = "vpn"
	ObjectKindTag             ObjectKind = "tag"
)

type ReferenceResolution string

const (
	ReferenceResolved   ReferenceResolution = "resolved"
	ReferenceUnresolved ReferenceResolution = "unresolved"
	ReferenceAmbiguous  ReferenceResolution = "ambiguous"
	ReferenceBuiltIn    ReferenceResolution = "built-in"
)

type ObjectRef struct {
	Kind       ObjectKind          `json:"kind"`
	ID         ObjectID            `json:"id,omitempty"`
	VendorName string              `json:"vendor_name"`
	Scope      ObjectScope         `json:"scope,omitempty"`
	Resolution ReferenceResolution `json:"resolution,omitempty"`
	Source     SourceLocation      `json:"source,omitempty"`
}

type ReferenceGraph struct {
	Edges      []ReferenceEdge       `json:"edges,omitempty"`
	Unresolved []UnresolvedReference `json:"unresolved,omitempty"`
}

type ReferenceEdge struct {
	From   ObjectID       `json:"from"`
	To     ObjectID       `json:"to"`
	Role   string         `json:"role"`
	Source SourceLocation `json:"source,omitempty"`
}

type UnresolvedReference struct {
	From       ObjectID       `json:"from"`
	Kind       ObjectKind     `json:"kind"`
	VendorName string         `json:"vendor_name"`
	Role       string         `json:"role"`
	Scope      ObjectScope    `json:"scope,omitempty"`
	Source     SourceLocation `json:"source,omitempty"`
}

type RoutingDomain struct {
	ID     ObjectID       `json:"id"`
	Name   string         `json:"name"`
	Scope  ObjectScope    `json:"scope"`
	Source SourceLocation `json:"source,omitempty"`
}

type InterfaceKind string

const (
	InterfacePhysical     InterfaceKind = "physical"
	InterfaceVLAN         InterfaceKind = "vlan"
	InterfaceSubinterface InterfaceKind = "subinterface"
	InterfaceAggregate    InterfaceKind = "aggregate"
	InterfaceRedundant    InterfaceKind = "redundant"
	InterfaceLoopback     InterfaceKind = "loopback"
	InterfaceTunnel       InterfaceKind = "tunnel"
	InterfaceIPsec        InterfaceKind = "ipsec"
	InterfaceGRE          InterfaceKind = "gre"
	InterfaceSoftware     InterfaceKind = "software-switch"
	InterfaceHardware     InterfaceKind = "hardware-switch"
	InterfaceUnknown      InterfaceKind = "unknown"
)

type Interface struct {
	ID               ObjectID        `json:"id"`
	Scope            ObjectScope     `json:"scope"`
	Name             string          `json:"name"`
	Alias            string          `json:"alias,omitempty"`
	Comments         string          `json:"comments,omitempty"`
	Kind             InterfaceKind   `json:"kind"`
	Role             string          `json:"role,omitempty"`
	Enabled          bool            `json:"enabled"`
	Parent           *ObjectRef      `json:"parent,omitempty"`
	Zones            []ObjectRef     `json:"zones,omitempty"`
	VRF              uint32          `json:"vrf,omitempty"`
	VLAN             *VLANAttachment `json:"vlan,omitempty"`
	Addresses4       []IPAssignment  `json:"addresses4,omitempty"`
	Addresses6       []IPAssignment  `json:"addresses6,omitempty"`
	MTU              int             `json:"mtu,omitempty"`
	ManagementAccess []string        `json:"management_access,omitempty"`
	DHCPMode         string          `json:"dhcp_mode,omitempty"`
	Tags             []ObjectRef     `json:"tags,omitempty"`
	Source           SourceLocation  `json:"source,omitempty"`
}

type VLANAttachment struct {
	VLANID uint16    `json:"vlan_id"`
	Parent ObjectRef `json:"parent"`
}

type IPAssignment struct {
	Address         string `json:"address"`
	Prefix          string `json:"prefix"`
	OriginalAddress string `json:"original_address,omitempty"`
	OriginalMask    string `json:"original_mask,omitempty"`
	OriginalCIDR    string `json:"original_cidr,omitempty"`
	Primary         bool   `json:"primary,omitempty"`
	Secondary       bool   `json:"secondary,omitempty"`
}

type VLAN struct {
	ID         ObjectID       `json:"id"`
	Scope      ObjectScope    `json:"scope"`
	VLANID     uint16         `json:"vlan_id"`
	Name       string         `json:"name"`
	Interfaces []ObjectRef    `json:"interfaces,omitempty"`
	Subnets    []ObjectRef    `json:"subnets,omitempty"`
	Tags       []ObjectRef    `json:"tags,omitempty"`
	Source     SourceLocation `json:"source,omitempty"`
}

type Subnet struct {
	ID              ObjectID       `json:"id"`
	Scope           ObjectScope    `json:"scope"`
	Name            string         `json:"name"`
	Prefix          string         `json:"prefix"`
	OriginalAddress string         `json:"original_address,omitempty"`
	OriginalMask    string         `json:"original_mask,omitempty"`
	OriginalCIDR    string         `json:"original_cidr,omitempty"`
	Interface       *ObjectRef     `json:"interface,omitempty"`
	VLAN            *ObjectRef     `json:"vlan,omitempty"`
	Source          SourceLocation `json:"source,omitempty"`
}

type DHCPConfiguration struct {
	Servers []DHCPServer `json:"servers,omitempty"`
	Relays  []DHCPRelay  `json:"relays,omitempty"`
}

type DHCPServer struct {
	ID           ObjectID          `json:"id"`
	Scope        ObjectScope       `json:"scope"`
	Name         string            `json:"name"`
	Enabled      bool              `json:"enabled"`
	Interface    ObjectRef         `json:"interface"`
	Network      string            `json:"network,omitempty"`
	Gateway      string            `json:"gateway,omitempty"`
	Netmask      string            `json:"netmask,omitempty"`
	Ranges       []DHCPRange       `json:"ranges,omitempty"`
	Reservations []DHCPReservation `json:"reservations,omitempty"`
	DNSServers   []string          `json:"dns_servers,omitempty"`
	NTPServers   []string          `json:"ntp_servers,omitempty"`
	DomainName   string            `json:"domain_name,omitempty"`
	LeaseSeconds uint64            `json:"lease_seconds,omitempty"`
	Options      []DHCPOption      `json:"options,omitempty"`
	Source       SourceLocation    `json:"source,omitempty"`
}

type DHCPRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type DHCPReservation struct {
	Name        string `json:"name"`
	MACAddress  string `json:"mac_address,omitempty"`
	IP          string `json:"ip,omitempty"`
	Description string `json:"description,omitempty"`
}
type DHCPOption struct {
	Code      uint16 `json:"code"`
	Value     string `json:"value"`
	ValueType string `json:"value_type,omitempty"`
}
type DHCPRelay struct {
	ID        ObjectID       `json:"id"`
	Scope     ObjectScope    `json:"scope"`
	Interface ObjectRef      `json:"interface"`
	ServerIPs []string       `json:"server_ips,omitempty"`
	Source    SourceLocation `json:"source,omitempty"`
}

type RoutingConfiguration struct {
	Routes       []Route          `json:"routes,omitempty"`
	PolicyRoutes []PolicyRoute    `json:"policy_routes,omitempty"`
	PrefixLists  []PrefixList     `json:"prefix_lists,omitempty"`
	RouteMaps    []RouteMap       `json:"route_maps,omitempty"`
	BGP          *RoutingProtocol `json:"bgp,omitempty"`
	OSPF         *RoutingProtocol `json:"ospf,omitempty"`
	RIP          *RoutingProtocol `json:"rip,omitempty"`
}

type Route struct {
	ID          ObjectID       `json:"id"`
	Scope       ObjectScope    `json:"scope"`
	Sequence    int            `json:"sequence"`
	Family      string         `json:"family,omitempty"`
	SourceType  string         `json:"source_type"`
	Destination string         `json:"destination"`
	Gateway     string         `json:"gateway,omitempty"`
	Interface   *ObjectRef     `json:"interface,omitempty"`
	Zone        *ObjectRef     `json:"zone,omitempty"`
	SDWAN       *ObjectRef     `json:"sdwan,omitempty"`
	Distance    uint32         `json:"distance,omitempty"`
	Priority    uint32         `json:"priority,omitempty"`
	Metric      uint32         `json:"metric,omitempty"`
	VRF         uint32         `json:"vrf,omitempty"`
	Table       uint32         `json:"table,omitempty"`
	Enabled     bool           `json:"enabled"`
	Blackhole   bool           `json:"blackhole,omitempty"`
	Comments    string         `json:"comments,omitempty"`
	Tags        []ObjectRef    `json:"tags,omitempty"`
	State       EvidenceState  `json:"state"`
	Source      SourceLocation `json:"source,omitempty"`
}

type PolicyRoute struct {
	ID               ObjectID       `json:"id"`
	Scope            ObjectScope    `json:"scope"`
	Sequence         int            `json:"sequence"`
	Enabled          bool           `json:"enabled"`
	Incoming         []ObjectRef    `json:"incoming,omitempty"`
	Sources          []string       `json:"sources,omitempty"`
	Destinations     []string       `json:"destinations,omitempty"`
	Protocols        []string       `json:"protocols,omitempty"`
	SourcePorts      []PortRange    `json:"source_ports,omitempty"`
	DestinationPorts []PortRange    `json:"destination_ports,omitempty"`
	Gateway          string         `json:"gateway,omitempty"`
	Interface        *ObjectRef     `json:"interface,omitempty"`
	DSCP             *DSCPMatch     `json:"dscp,omitempty"`
	Comments         string         `json:"comments,omitempty"`
	Source           SourceLocation `json:"source,omitempty"`
}

type PrefixList struct {
	ID      ObjectID          `json:"id"`
	Scope   ObjectScope       `json:"scope"`
	Name    string            `json:"name"`
	Entries []PrefixListEntry `json:"entries,omitempty"`
	Source  SourceLocation    `json:"source,omitempty"`
}
type PrefixListEntry struct {
	Sequence  int    `json:"sequence"`
	Action    string `json:"action"`
	Prefix    string `json:"prefix"`
	MinLength *uint8 `json:"min_length,omitempty"`
	MaxLength *uint8 `json:"max_length,omitempty"`
}
type RouteMap struct {
	ID     ObjectID       `json:"id"`
	Scope  ObjectScope    `json:"scope"`
	Name   string         `json:"name"`
	Rules  []RouteMapRule `json:"rules,omitempty"`
	Source SourceLocation `json:"source,omitempty"`
}
type RouteMapRule struct {
	Sequence           int         `json:"sequence"`
	Action             string      `json:"action"`
	MatchPrefixLists   []ObjectRef `json:"match_prefix_lists,omitempty"`
	MatchInterfaces    []ObjectRef `json:"match_interfaces,omitempty"`
	SetMetric          *uint32     `json:"set_metric,omitempty"`
	SetLocalPreference *uint32     `json:"set_local_preference,omitempty"`
	SetNextHop         string      `json:"set_next_hop,omitempty"`
	SetCommunities     []string    `json:"set_communities,omitempty"`
	SetRouteTag        *uint32     `json:"set_route_tag,omitempty"`
}
type RoutingProtocol struct {
	Enabled        bool           `json:"enabled"`
	RouterID       string         `json:"router_id,omitempty"`
	ASN            string         `json:"asn,omitempty"`
	RawObjectNames []string       `json:"raw_object_names,omitempty"`
	Source         SourceLocation `json:"source,omitempty"`
}

type SDWANConfiguration struct {
	Enabled      bool               `json:"enabled"`
	Zones        []SDWANZone        `json:"zones,omitempty"`
	Members      []SDWANMember      `json:"members,omitempty"`
	HealthChecks []SDWANHealthCheck `json:"health_checks,omitempty"`
	Rules        []SDWANRule        `json:"rules,omitempty"`
}

type SDWANZone struct {
	ID      ObjectID       `json:"id"`
	Scope   ObjectScope    `json:"scope"`
	Name    string         `json:"name"`
	Members []ObjectRef    `json:"members,omitempty"`
	Source  SourceLocation `json:"source,omitempty"`
}
type SDWANMember struct {
	ID        ObjectID       `json:"id"`
	Scope     ObjectScope    `json:"scope"`
	MemberID  uint32         `json:"member_id"`
	Interface ObjectRef      `json:"interface"`
	Zone      *ObjectRef     `json:"zone,omitempty"`
	Gateway   string         `json:"gateway,omitempty"`
	Cost      uint32         `json:"cost,omitempty"`
	Priority  uint32         `json:"priority,omitempty"`
	Weight    uint32         `json:"weight,omitempty"`
	SourceIP  string         `json:"source_ip,omitempty"`
	Enabled   bool           `json:"enabled"`
	Source    SourceLocation `json:"source,omitempty"`
}
type SDWANHealthCheck struct {
	ID                  ObjectID       `json:"id"`
	Scope               ObjectScope    `json:"scope"`
	Name                string         `json:"name"`
	Server              []string       `json:"server,omitempty"`
	Protocol            string         `json:"protocol,omitempty"`
	Members             []ObjectRef    `json:"members,omitempty"`
	LatencyThreshold    uint32         `json:"latency_threshold,omitempty"`
	JitterThreshold     uint32         `json:"jitter_threshold,omitempty"`
	PacketLossThreshold uint32         `json:"packet_loss_threshold,omitempty"`
	Source              SourceLocation `json:"source,omitempty"`
}
type SDWANRule struct {
	ID                   ObjectID       `json:"id"`
	Scope                ObjectScope    `json:"scope"`
	Sequence             int            `json:"sequence"`
	Name                 string         `json:"name"`
	Enabled              bool           `json:"enabled"`
	SourceAddresses      []ObjectRef    `json:"source_addresses,omitempty"`
	DestinationAddresses []ObjectRef    `json:"destination_addresses,omitempty"`
	Services             []ObjectRef    `json:"services,omitempty"`
	Applications         []string       `json:"applications,omitempty"`
	InternetServices     []string       `json:"internet_services,omitempty"`
	DSCP                 *DSCPMatch     `json:"dscp,omitempty"`
	Strategy             string         `json:"strategy,omitempty"`
	PreferredMembers     []ObjectRef    `json:"preferred_members,omitempty"`
	RequiredSLAs         []ObjectRef    `json:"required_slas,omitempty"`
	Source               SourceLocation `json:"source,omitempty"`
}

type Policy struct {
	ID       ObjectID       `json:"id"`
	Scope    ObjectScope    `json:"scope"`
	VendorID string         `json:"vendor_id,omitempty"`
	UUID     string         `json:"uuid,omitempty"`
	Name     string         `json:"name"`
	Comments string         `json:"comments,omitempty"`
	Enabled  bool           `json:"enabled"`
	Sequence int            `json:"sequence"`
	Action   string         `json:"action"`
	Match    PolicyMatch    `json:"match"`
	NAT      PolicyNAT      `json:"nat"`
	QoS      PolicyQoS      `json:"qos"`
	Security PolicySecurity `json:"security"`
	Logging  PolicyLogging  `json:"logging"`
	Tags     []ObjectRef    `json:"tags,omitempty"`
	Source   SourceLocation `json:"source,omitempty"`
}

type PolicyMatch struct {
	SourceInterfaces            []ObjectRef `json:"source_interfaces,omitempty"`
	DestinationInterfaces       []ObjectRef `json:"destination_interfaces,omitempty"`
	SourceAddresses             []ObjectRef `json:"source_addresses,omitempty"`
	DestinationAddresses        []ObjectRef `json:"destination_addresses,omitempty"`
	SourceAddressesNegated      bool        `json:"source_addresses_negated,omitempty"`
	DestinationAddressesNegated bool        `json:"destination_addresses_negated,omitempty"`
	Services                    []ObjectRef `json:"services,omitempty"`
	ServicesNegated             bool        `json:"services_negated,omitempty"`
	Schedules                   []ObjectRef `json:"schedules,omitempty"`
	Users                       []ObjectRef `json:"users,omitempty"`
	UserGroups                  []ObjectRef `json:"user_groups,omitempty"`
	InternetServices            []string    `json:"internet_services,omitempty"`
	Applications                []string    `json:"applications,omitempty"`
	SourcePorts                 []PortRange `json:"source_ports,omitempty"`
	DestinationPorts            []PortRange `json:"destination_ports,omitempty"`
	Protocols                   []string    `json:"protocols,omitempty"`
	DSCP                        *DSCPMatch  `json:"dscp,omitempty"`
	VLAN                        *VLANMatch  `json:"vlan,omitempty"`
}

type PolicyNAT struct {
	Enabled            bool        `json:"enabled"`
	SourceNAT          bool        `json:"source_nat,omitempty"`
	FixedPort          bool        `json:"fixed_port,omitempty"`
	IPPools            []ObjectRef `json:"ip_pools,omitempty"`
	VIPs               []ObjectRef `json:"vips,omitempty"`
	PreserveSourcePort bool        `json:"preserve_source_port,omitempty"`
	CentralNATApplied  bool        `json:"central_nat_applied,omitempty"`
}
type PolicyQoS struct {
	ForwardShaper  *ObjectRef   `json:"forward_shaper,omitempty"`
	ReverseShaper  *ObjectRef   `json:"reverse_shaper,omitempty"`
	PerIPShaper    *ObjectRef   `json:"per_ip_shaper,omitempty"`
	ForwardDSCP    *DSCPRewrite `json:"forward_dscp,omitempty"`
	ReverseDSCP    *DSCPRewrite `json:"reverse_dscp,omitempty"`
	ForwardVLANCoS *uint8       `json:"forward_vlan_cos,omitempty"`
	ReverseVLANCoS *uint8       `json:"reverse_vlan_cos,omitempty"`
}
type PolicySecurity struct {
	InspectionMode   string `json:"inspection_mode,omitempty"`
	AVProfile        string `json:"av_profile,omitempty"`
	WebFilterProfile string `json:"web_filter_profile,omitempty"`
	IPSProfile       string `json:"ips_profile,omitempty"`
	ApplicationList  string `json:"application_list,omitempty"`
	SSLSSHProfile    string `json:"ssl_ssh_profile,omitempty"`
}
type PolicyLogging struct {
	Enabled     bool `json:"enabled"`
	AllSessions bool `json:"all_sessions,omitempty"`
	UTM         bool `json:"utm,omitempty"`
	Start       bool `json:"start,omitempty"`
}

type PortRange struct {
	Start uint16 `json:"start"`
	End   uint16 `json:"end"`
}
type DSCPValue struct {
	Value         uint8  `json:"value"`
	CanonicalName string `json:"canonical_name,omitempty"`
	VendorName    string `json:"vendor_name,omitempty"`
	Original      string `json:"original,omitempty"`
}
type DSCPMatch struct {
	Enabled bool        `json:"enabled"`
	Negated bool        `json:"negated,omitempty"`
	Values  []DSCPValue `json:"values,omitempty"`
}
type DSCPRewrite struct {
	Enabled bool       `json:"enabled"`
	Copy    bool       `json:"copy,omitempty"`
	Value   *DSCPValue `json:"value,omitempty"`
}
type VLANMatch struct {
	Enabled bool     `json:"enabled"`
	VLANIDs []uint16 `json:"vlan_ids,omitempty"`
	Negated bool     `json:"negated,omitempty"`
}

type ObjectCatalog struct {
	Addresses        []AddressObject `json:"addresses,omitempty"`
	AddressGroups    []AddressGroup  `json:"address_groups,omitempty"`
	Services         []ServiceObject `json:"services,omitempty"`
	ServiceGroups    []ServiceGroup  `json:"service_groups,omitempty"`
	Schedules        []Schedule      `json:"schedules,omitempty"`
	Users            []NamedObject   `json:"users,omitempty"`
	UserGroups       []ObjectGroup   `json:"user_groups,omitempty"`
	InternetServices []NamedObject   `json:"internet_services,omitempty"`
	SecurityProfiles []NamedObject   `json:"security_profiles,omitempty"`
}

type AddressObject struct {
	ID        ObjectID       `json:"id"`
	Scope     ObjectScope    `json:"scope"`
	Name      string         `json:"name"`
	Kind      string         `json:"kind"`
	Prefixes  []string       `json:"prefixes,omitempty"`
	IPRanges  []IPRange      `json:"ip_ranges,omitempty"`
	FQDNs     []string       `json:"fqdns,omitempty"`
	Wildcards []string       `json:"wildcards,omitempty"`
	MACs      []string       `json:"macs,omitempty"`
	Countries []string       `json:"countries,omitempty"`
	Interface *ObjectRef     `json:"interface,omitempty"`
	Tags      []ObjectRef    `json:"tags,omitempty"`
	Comments  string         `json:"comments,omitempty"`
	Source    SourceLocation `json:"source,omitempty"`
}
type AddressGroup struct {
	ID             ObjectID       `json:"id"`
	Scope          ObjectScope    `json:"scope"`
	Name           string         `json:"name"`
	Members        []ObjectRef    `json:"members,omitempty"`
	ExcludeMembers []ObjectRef    `json:"exclude_members,omitempty"`
	Comments       string         `json:"comments,omitempty"`
	Source         SourceLocation `json:"source,omitempty"`
}
type ServiceObject struct {
	ID             ObjectID       `json:"id"`
	Scope          ObjectScope    `json:"scope"`
	Name           string         `json:"name"`
	Protocol       string         `json:"protocol,omitempty"`
	TCPPorts       []PortRange    `json:"tcp_ports,omitempty"`
	UDPPorts       []PortRange    `json:"udp_ports,omitempty"`
	SCTPPorts      []PortRange    `json:"sctp_ports,omitempty"`
	ICMPTypes      []string       `json:"icmp_types,omitempty"`
	ProtocolNumber *uint8         `json:"protocol_number,omitempty"`
	SourcePorts    []PortRange    `json:"source_ports,omitempty"`
	Comments       string         `json:"comments,omitempty"`
	Source         SourceLocation `json:"source,omitempty"`
}
type ServiceGroup struct {
	ID       ObjectID       `json:"id"`
	Scope    ObjectScope    `json:"scope"`
	Name     string         `json:"name"`
	Members  []ObjectRef    `json:"members,omitempty"`
	Comments string         `json:"comments,omitempty"`
	Source   SourceLocation `json:"source,omitempty"`
}
type Schedule struct {
	ID     ObjectID       `json:"id"`
	Scope  ObjectScope    `json:"scope"`
	Name   string         `json:"name"`
	Kind   string         `json:"kind,omitempty"`
	Start  string         `json:"start,omitempty"`
	End    string         `json:"end,omitempty"`
	Days   []string       `json:"days,omitempty"`
	Source SourceLocation `json:"source,omitempty"`
}
type NamedObject struct {
	ID     ObjectID       `json:"id"`
	Scope  ObjectScope    `json:"scope"`
	Name   string         `json:"name"`
	Kind   string         `json:"kind,omitempty"`
	Source SourceLocation `json:"source,omitempty"`
}
type ObjectGroup struct {
	ID      ObjectID       `json:"id"`
	Scope   ObjectScope    `json:"scope"`
	Name    string         `json:"name"`
	Members []ObjectRef    `json:"members,omitempty"`
	Source  SourceLocation `json:"source,omitempty"`
}
type IPRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type NATConfiguration struct {
	VirtualIPs      []VirtualIP       `json:"virtual_ips,omitempty"`
	VirtualIPGroups []VirtualIPGroup  `json:"virtual_ip_groups,omitempty"`
	IPPools         []IPPool          `json:"ip_pools,omitempty"`
	CentralSNAT     []CentralSNATRule `json:"central_snat,omitempty"`
}
type VirtualIP struct {
	ID                ObjectID       `json:"id"`
	Scope             ObjectScope    `json:"scope"`
	Name              string         `json:"name"`
	Enabled           bool           `json:"enabled"`
	ExternalAddress   string         `json:"external_address,omitempty"`
	MappedAddresses   []string       `json:"mapped_addresses,omitempty"`
	ExternalInterface *ObjectRef     `json:"external_interface,omitempty"`
	Protocol          string         `json:"protocol,omitempty"`
	ExternalPorts     []PortRange    `json:"external_ports,omitempty"`
	MappedPorts       []PortRange    `json:"mapped_ports,omitempty"`
	PortForwarding    bool           `json:"port_forwarding,omitempty"`
	SourceFilters     []string       `json:"source_filters,omitempty"`
	Comments          string         `json:"comments,omitempty"`
	Tags              []ObjectRef    `json:"tags,omitempty"`
	Source            SourceLocation `json:"source,omitempty"`
}
type VirtualIPGroup struct {
	ID      ObjectID       `json:"id"`
	Scope   ObjectScope    `json:"scope"`
	Name    string         `json:"name"`
	Members []ObjectRef    `json:"members,omitempty"`
	Source  SourceLocation `json:"source,omitempty"`
}
type IPPool struct {
	ID      ObjectID       `json:"id"`
	Scope   ObjectScope    `json:"scope"`
	Name    string         `json:"name"`
	StartIP string         `json:"start_ip,omitempty"`
	EndIP   string         `json:"end_ip,omitempty"`
	Type    string         `json:"type,omitempty"`
	Source  SourceLocation `json:"source,omitempty"`
}
type CentralSNATRule struct {
	ID                    ObjectID       `json:"id"`
	Scope                 ObjectScope    `json:"scope"`
	Sequence              int            `json:"sequence"`
	SourceInterfaces      []ObjectRef    `json:"source_interfaces,omitempty"`
	DestinationInterfaces []ObjectRef    `json:"destination_interfaces,omitempty"`
	SourceAddresses       []ObjectRef    `json:"source_addresses,omitempty"`
	DestinationAddresses  []ObjectRef    `json:"destination_addresses,omitempty"`
	IPPool                *ObjectRef     `json:"ip_pool,omitempty"`
	NATEnabled            bool           `json:"nat_enabled"`
	Source                SourceLocation `json:"source,omitempty"`
}

type VPN struct {
	ID             ObjectID       `json:"id"`
	Scope          ObjectScope    `json:"scope"`
	Name           string         `json:"name"`
	Kind           string         `json:"kind"`
	Interface      *ObjectRef     `json:"interface,omitempty"`
	RemoteGateway  string         `json:"remote_gateway,omitempty"`
	LocalNetworks  []ObjectRef    `json:"local_networks,omitempty"`
	RemoteNetworks []ObjectRef    `json:"remote_networks,omitempty"`
	Enabled        bool           `json:"enabled"`
	Source         SourceLocation `json:"source,omitempty"`
}

type QoSConfiguration struct {
	TrafficShapers  []TrafficShaper `json:"traffic_shapers,omitempty"`
	PerIPShapers    []PerIPShaper   `json:"per_ip_shapers,omitempty"`
	ShapingPolicies []ShapingPolicy `json:"shaping_policies,omitempty"`
	DSCPMappings    []DSCPMapping   `json:"dscp_mappings,omitempty"`
}
type TrafficShaper struct {
	ID                  ObjectID       `json:"id"`
	Scope               ObjectScope    `json:"scope"`
	Name                string         `json:"name"`
	GuaranteedBandwidth Bandwidth      `json:"guaranteed_bandwidth"`
	MaximumBandwidth    Bandwidth      `json:"maximum_bandwidth"`
	BurstBytes          uint64         `json:"burst_bytes,omitempty"`
	Priority            string         `json:"priority,omitempty"`
	Mode                string         `json:"mode,omitempty"`
	PerPolicy           bool           `json:"per_policy,omitempty"`
	DSCPRemark          *DSCPValue     `json:"dscp_remark,omitempty"`
	Comments            string         `json:"comments,omitempty"`
	Source              SourceLocation `json:"source,omitempty"`
}
type PerIPShaper struct {
	ID               ObjectID       `json:"id"`
	Scope            ObjectScope    `json:"scope"`
	Name             string         `json:"name"`
	MaximumBandwidth Bandwidth      `json:"maximum_bandwidth"`
	ConcurrentLimit  uint32         `json:"concurrent_limit,omitempty"`
	DSCPRemark       *DSCPValue     `json:"dscp_remark,omitempty"`
	Source           SourceLocation `json:"source,omitempty"`
}
type Bandwidth struct {
	BitsPerSecond uint64 `json:"bits_per_second,omitempty"`
	OriginalValue string `json:"original_value,omitempty"`
	OriginalUnit  string `json:"original_unit,omitempty"`
}
type ShapingPolicy struct {
	ID                    ObjectID       `json:"id"`
	Scope                 ObjectScope    `json:"scope"`
	Sequence              int            `json:"sequence"`
	Enabled               bool           `json:"enabled"`
	SourceInterfaces      []ObjectRef    `json:"source_interfaces,omitempty"`
	DestinationInterfaces []ObjectRef    `json:"destination_interfaces,omitempty"`
	SourceAddresses       []ObjectRef    `json:"source_addresses,omitempty"`
	DestinationAddresses  []ObjectRef    `json:"destination_addresses,omitempty"`
	Services              []ObjectRef    `json:"services,omitempty"`
	Applications          []string       `json:"applications,omitempty"`
	DSCP                  *DSCPMatch     `json:"dscp,omitempty"`
	ForwardShaper         *ObjectRef     `json:"forward_shaper,omitempty"`
	ReverseShaper         *ObjectRef     `json:"reverse_shaper,omitempty"`
	PerIPShaper           *ObjectRef     `json:"per_ip_shaper,omitempty"`
	Source                SourceLocation `json:"source,omitempty"`
}
type DSCPMapping struct {
	Name  string    `json:"name"`
	Value DSCPValue `json:"value"`
}

type Tag struct {
	ID         ObjectID       `json:"id"`
	Scope      ObjectScope    `json:"scope"`
	Name       string         `json:"name"`
	Category   string         `json:"category,omitempty"`
	Value      string         `json:"value,omitempty"`
	Color      string         `json:"color,omitempty"`
	VendorName string         `json:"vendor_name,omitempty"`
	Source     SourceLocation `json:"source,omitempty"`
}

type Finding struct {
	ID       string         `json:"id"`
	Severity string         `json:"severity"`
	Category string         `json:"category"`
	Title    string         `json:"title"`
	Detail   string         `json:"detail"`
	ObjectID ObjectID       `json:"object_id,omitempty"`
	State    EvidenceState  `json:"state"`
	Source   SourceLocation `json:"source,omitempty"`
}
type Extension struct {
	Scope  ObjectScope `json:"scope,omitempty"`
	Path   string      `json:"path"`
	Name   string      `json:"name,omitempty"`
	Reason string      `json:"reason,omitempty"`
}
