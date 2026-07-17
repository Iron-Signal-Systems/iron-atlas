package snapshot

// RecordCount is an aggregate count for a fixed normalized snapshot record
// kind. Kind values are Atlas-owned labels and never contain vendor data.
type RecordCount struct {
	Kind  string
	Count int
}

// NormalizedRecordCounts returns every supported snapshot record kind in a
// stable order, including kinds with a zero count. This makes coverage gaps
// visible without exposing source-derived names or values.
func NormalizedRecordCounts(value *FirewallSnapshot) []RecordCount {
	if value == nil {
		value = &FirewallSnapshot{}
	}
	return []RecordCount{
		{Kind: "address", Count: len(value.Objects.Addresses)},
		{Kind: "address-group", Count: len(value.Objects.AddressGroups)},
		{Kind: "bgp-configuration", Count: boolCount(value.Routing.BGP != nil)},
		{Kind: "central-snat-rule", Count: len(value.NAT.CentralSNAT)},
		{Kind: "dhcp-relay", Count: len(value.DHCP.Relays)},
		{Kind: "dhcp-server", Count: len(value.DHCP.Servers)},
		{Kind: "dscp-mapping", Count: len(value.QoS.DSCPMappings)},
		{Kind: "extension", Count: len(value.Extensions)},
		{Kind: "interface", Count: len(value.Interfaces)},
		{Kind: "internet-service", Count: len(value.Objects.InternetServices)},
		{Kind: "ip-pool", Count: len(value.NAT.IPPools)},
		{Kind: "ospf-configuration", Count: boolCount(value.Routing.OSPF != nil)},
		{Kind: "per-ip-shaper", Count: len(value.QoS.PerIPShapers)},
		{Kind: "policy", Count: len(value.Policies)},
		{Kind: "policy-route", Count: len(value.Routing.PolicyRoutes)},
		{Kind: "prefix-list", Count: len(value.Routing.PrefixLists)},
		{Kind: "rip-configuration", Count: boolCount(value.Routing.RIP != nil)},
		{Kind: "route", Count: len(value.Routing.Routes)},
		{Kind: "route-map", Count: len(value.Routing.RouteMaps)},
		{Kind: "routing-domain", Count: len(value.Domains)},
		{Kind: "schedule", Count: len(value.Objects.Schedules)},
		{Kind: "sdwan-health-check", Count: len(value.SDWAN.HealthChecks)},
		{Kind: "sdwan-member", Count: len(value.SDWAN.Members)},
		{Kind: "sdwan-rule", Count: len(value.SDWAN.Rules)},
		{Kind: "sdwan-zone", Count: len(value.SDWAN.Zones)},
		{Kind: "security-profile", Count: len(value.Objects.SecurityProfiles)},
		{Kind: "service", Count: len(value.Objects.Services)},
		{Kind: "service-group", Count: len(value.Objects.ServiceGroups)},
		{Kind: "shaping-policy", Count: len(value.QoS.ShapingPolicies)},
		{Kind: "subnet", Count: len(value.Subnets)},
		{Kind: "tag", Count: len(value.Tags)},
		{Kind: "traffic-shaper", Count: len(value.QoS.TrafficShapers)},
		{Kind: "user", Count: len(value.Objects.Users)},
		{Kind: "user-group", Count: len(value.Objects.UserGroups)},
		{Kind: "virtual-ip", Count: len(value.NAT.VirtualIPs)},
		{Kind: "virtual-ip-group", Count: len(value.NAT.VirtualIPGroups)},
		{Kind: "vlan", Count: len(value.VLANs)},
		{Kind: "vpn", Count: len(value.VPNs)},
	}
}

// TotalNormalizedRecords returns the sum of all fixed record-kind counts.
func TotalNormalizedRecords(value *FirewallSnapshot) int {
	total := 0
	for _, record := range NormalizedRecordCounts(value) {
		total += record.Count
	}
	return total
}

func boolCount(value bool) int {
	if value {
		return 1
	}
	return 0
}
