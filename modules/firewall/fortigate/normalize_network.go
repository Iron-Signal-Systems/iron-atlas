package fortigate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func (n *normalizer) parseInterfaces(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	section := childAlias(vdom, "system_interface")
	entries := tableEntries(section, basePath+".system_interface")

	// Register all interface names before resolving parent and zone references.
	for _, entry := range entries {
		id := objectID(scope, "interface", entry.Name)
		n.register(scope, snapshot.ObjectKindInterface, entry.Name, id)
	}

	for _, entry := range entries {
		id := objectID(scope, "interface", entry.Name)
		typeValue := normalizeKey(scalarField(entry.Node, "type"))
		kind := snapshot.InterfacePhysical
		switch typeValue {
		case "vlan":
			kind = snapshot.InterfaceVLAN
		case "aggregate", "8023ad_aggregate":
			kind = snapshot.InterfaceAggregate
		case "redundant":
			kind = snapshot.InterfaceRedundant
		case "loopback":
			kind = snapshot.InterfaceLoopback
		case "tunnel":
			kind = snapshot.InterfaceTunnel
		case "software_switch":
			kind = snapshot.InterfaceSoftware
		case "hard_switch", "hardware_switch":
			kind = snapshot.InterfaceHardware
		case "", "physical":
			kind = snapshot.InterfacePhysical
		default:
			kind = snapshot.InterfaceUnknown
		}

		iface := snapshot.Interface{
			ID: id, Scope: scope, Name: entry.Name,
			Alias: scalarField(entry.Node, "alias"), Comments: firstNonEmpty(scalarField(entry.Node, "description"), scalarField(entry.Node, "comments")),
			Kind: kind, Role: scalarField(entry.Node, "role"), Enabled: enabledField(entry.Node, true, "status"),
			VRF: uint32Field(entry.Node, "vrf"), MTU: intField(entry.Node, "mtu"),
			ManagementAccess: valuesField(entry.Node, "allowaccess"), DHCPMode: scalarField(entry.Node, "mode"),
			Source: location(entry.Node, entry.Path),
		}

		parentName := scalarField(entry.Node, "interface")
		if parentName != "" {
			parent := n.ref(id, "parent-interface", snapshot.ObjectKindInterface, parentName, scope, iface.Source)
			iface.Parent = &parent
		}

		vlanID := uint16Value(scalarField(entry.Node, "vlanid"))
		if vlanID > 0 {
			kind = snapshot.InterfaceVLAN
			iface.Kind = kind
			parent := snapshot.ObjectRef{Kind: snapshot.ObjectKindInterface, VendorName: parentName, Scope: scope, Resolution: snapshot.ReferenceUnresolved, Source: iface.Source}
			if iface.Parent != nil {
				parent = *iface.Parent
			}
			iface.VLAN = &snapshot.VLANAttachment{VLANID: vlanID, Parent: parent}
			vlanIDValue := objectID(scope, "vlan", fmt.Sprintf("%d", vlanID))
			vlan := snapshot.VLAN{
				ID: vlanIDValue, Scope: scope, VLANID: vlanID, Name: entry.Name,
				Interfaces: []snapshot.ObjectRef{{Kind: snapshot.ObjectKindInterface, ID: id, VendorName: entry.Name, Scope: scope, Resolution: snapshot.ReferenceResolved, Source: iface.Source}},
				Source:     iface.Source,
			}
			n.register(scope, snapshot.ObjectKindVLAN, entry.Name, vlanIDValue)
			n.register(scope, snapshot.ObjectKindVLAN, strconv.Itoa(int(vlanID)), vlanIDValue)
			n.out.VLANs = append(n.out.VLANs, vlan)
		}

		if values := valuesField(entry.Node, "ip"); len(values) > 0 {
			if prefix, address, mask, ok := canonicalPrefix(values); ok {
				assignment := snapshot.IPAssignment{Address: address, Prefix: prefix, OriginalAddress: address, OriginalMask: mask, OriginalCIDR: strings.Join(values, "/"), Primary: true}
				iface.Addresses4 = append(iface.Addresses4, assignment)
				subnetID := objectID(scope, "subnet", entry.Name+"-primary")
				subnet := snapshot.Subnet{ID: subnetID, Scope: scope, Name: entry.Name + " primary", Prefix: prefix, OriginalAddress: address, OriginalMask: mask, OriginalCIDR: strings.Join(values, "/"), Source: iface.Source}
				interfaceRef := snapshot.ObjectRef{Kind: snapshot.ObjectKindInterface, ID: id, VendorName: entry.Name, Scope: scope, Resolution: snapshot.ReferenceResolved, Source: iface.Source}
				subnet.Interface = &interfaceRef
				if iface.VLAN != nil {
					vlanRef := snapshot.ObjectRef{Kind: snapshot.ObjectKindVLAN, VendorName: strconv.Itoa(int(iface.VLAN.VLANID)), Scope: scope, Resolution: snapshot.ReferenceResolved, ID: objectID(scope, "vlan", strconv.Itoa(int(iface.VLAN.VLANID))), Source: iface.Source}
					subnet.VLAN = &vlanRef
				}
				n.register(scope, snapshot.ObjectKindSubnet, subnet.Name, subnetID)
				n.out.Subnets = append(n.out.Subnets, subnet)
			} else {
				n.addFinding("warning", "interface", "Invalid interface address", fmt.Sprintf("Interface %s has an IP value Atlas could not normalize: %v", entry.Name, values), id, iface.Source)
			}
		}

		if values := valuesField(entry.Node, "ip6-address", "ip6_address"); len(values) > 0 {
			if prefix, address, mask, ok := canonicalPrefix(values); ok {
				iface.Addresses6 = append(iface.Addresses6, snapshot.IPAssignment{Address: address, Prefix: prefix, OriginalAddress: address, OriginalMask: mask, OriginalCIDR: strings.Join(values, "/"), Primary: true})
			}
		}

		n.out.Interfaces = append(n.out.Interfaces, iface)
	}

	// FortiGate zones are separate from interfaces and can be referenced by policies.
	zoneSection := childAlias(vdom, "system_zone")
	for _, entry := range tableEntries(zoneSection, basePath+".system_zone") {
		id := objectID(scope, "zone", entry.Name)
		n.register(scope, snapshot.ObjectKindZone, entry.Name, id)
		members := n.refs(id, "zone-member", snapshot.ObjectKindInterface, valuesField(entry.Node, "interface"), scope, location(entry.Node, entry.Path))
		for i := range n.out.Interfaces {
			for _, member := range members {
				if member.Resolution == snapshot.ReferenceResolved && n.out.Interfaces[i].ID == member.ID {
					n.out.Interfaces[i].Zones = append(n.out.Interfaces[i].Zones, snapshot.ObjectRef{Kind: snapshot.ObjectKindZone, ID: id, VendorName: entry.Name, Scope: scope, Resolution: snapshot.ReferenceResolved, Source: location(entry.Node, entry.Path)})
				}
			}
		}
	}
}

func (n *normalizer) parseRoutes(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	section := childAlias(vdom, "router_static")
	for _, entry := range tableEntries(section, basePath+".router_static") {
		id := objectID(scope, "route", entry.Name)
		destinationValues := valuesField(entry.Node, "dst")
		destination := "0.0.0.0/0"
		if len(destinationValues) > 0 {
			if prefix, _, _, ok := canonicalPrefix(destinationValues); ok {
				destination = prefix
			} else {
				destination = strings.Join(destinationValues, " ")
			}
		}
		route := snapshot.Route{
			ID: id, Scope: scope, Sequence: sequenceNumber(entry), Family: "ipv4", SourceType: "static",
			Destination: destination, Gateway: scalarField(entry.Node, "gateway"), Distance: uint32Field(entry.Node, "distance"),
			Priority: uint32Field(entry.Node, "priority"), VRF: uint32Field(entry.Node, "vrf"),
			Enabled: enabledField(entry.Node, true, "status"), Blackhole: enabledField(entry.Node, false, "blackhole"),
			Comments: scalarField(entry.Node, "comment", "comments"), State: snapshot.EvidenceConfigured, Source: location(entry.Node, entry.Path),
		}
		device := scalarField(entry.Node, "device")
		if device != "" {
			ref := n.ref(id, "route-egress", snapshot.ObjectKindInterfaceOrZone, device, scope, route.Source)
			switch ref.Kind {
			case snapshot.ObjectKindSDWANZone:
				route.SDWAN = &ref
			case snapshot.ObjectKindZone:
				route.Zone = &ref
			default:
				route.Interface = &ref
			}
		}
		if !route.Enabled {
			n.addFinding("info", "routing", "Disabled static route", fmt.Sprintf("Static route %s is disabled", entry.Name), id, route.Source)
		}
		n.out.Routing.Routes = append(n.out.Routing.Routes, route)
	}

	policySection := childAlias(vdom, "router_policy")
	for _, entry := range tableEntries(policySection, basePath+".router_policy") {
		id := objectID(scope, "policy-route", entry.Name)
		policyRoute := snapshot.PolicyRoute{
			ID: id, Scope: scope, Sequence: sequenceNumber(entry), Enabled: enabledField(entry.Node, true, "status"),
			Sources: valuesField(entry.Node, "src"), Destinations: valuesField(entry.Node, "dst"),
			Protocols: valuesField(entry.Node, "protocol"), SourcePorts: parsePortRanges(valuesField(entry.Node, "start-port", "src-port")),
			DestinationPorts: parsePortRanges(valuesField(entry.Node, "end-port", "dst-port")), Gateway: scalarField(entry.Node, "gateway"),
			Comments: scalarField(entry.Node, "comments", "comment"), Source: location(entry.Node, entry.Path),
		}
		policyRoute.Incoming = n.refs(id, "policy-route-ingress", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "input-device", "input_device"), scope, policyRoute.Source)
		if device := scalarField(entry.Node, "output-device", "output_device"); device != "" {
			ref := n.ref(id, "policy-route-egress", snapshot.ObjectKindInterfaceOrZone, device, scope, policyRoute.Source)
			policyRoute.Interface = &ref
		}
		if dscp := parseDSCP(scalarField(entry.Node, "tos", "dscp")); dscp != nil {
			policyRoute.DSCP = &snapshot.DSCPMatch{Enabled: true, Values: []snapshot.DSCPValue{*dscp}}
		}
		n.out.Routing.PolicyRoutes = append(n.out.Routing.PolicyRoutes, policyRoute)
	}

	n.parseRoutingObjects(vdom, scope, basePath)
}

func (n *normalizer) parseRoutingObjects(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	for _, spec := range []struct {
		section string
		kind    snapshot.ObjectKind
	}{
		{"router_prefix_list", snapshot.ObjectKindPrefixList},
		{"router_prefix_list6", snapshot.ObjectKindPrefixList},
	} {
		for _, entry := range tableEntries(childAlias(vdom, spec.section), basePath+"."+spec.section) {
			id := objectID(scope, string(spec.kind), entry.Name)
			n.register(scope, spec.kind, entry.Name, id)
			prefixList := snapshot.PrefixList{ID: id, Scope: scope, Name: entry.Name, Source: location(entry.Node, entry.Path)}
			rules := childAlias(entry.Node, "rule")
			for _, ruleEntry := range tableEntries(rules, entry.Path+".rule") {
				prefixList.Entries = append(prefixList.Entries, snapshot.PrefixListEntry{
					Sequence: sequenceNumber(ruleEntry), Action: firstNonEmpty(scalarField(ruleEntry.Node, "action"), "permit"),
					Prefix: firstNonEmpty(scalarField(ruleEntry.Node, "prefix", "prefix6"), "0.0.0.0/0"),
				})
			}
			n.out.Routing.PrefixLists = append(n.out.Routing.PrefixLists, prefixList)
		}
	}

	for _, entry := range tableEntries(childAlias(vdom, "router_route_map"), basePath+".router_route_map") {
		id := objectID(scope, "route-map", entry.Name)
		n.register(scope, snapshot.ObjectKindRouteMap, entry.Name, id)
		routeMap := snapshot.RouteMap{ID: id, Scope: scope, Name: entry.Name, Source: location(entry.Node, entry.Path)}
		for _, ruleEntry := range tableEntries(childAlias(entry.Node, "rule"), entry.Path+".rule") {
			rule := snapshot.RouteMapRule{Sequence: sequenceNumber(ruleEntry), Action: firstNonEmpty(scalarField(ruleEntry.Node, "action"), "permit")}
			rule.MatchPrefixLists = n.refs(id, "route-map-prefix-list", snapshot.ObjectKindPrefixList, valuesField(ruleEntry.Node, "match-ip-address", "match_ip_address"), scope, location(ruleEntry.Node, ruleEntry.Path))
			rule.MatchInterfaces = n.refs(id, "route-map-interface", snapshot.ObjectKindInterface, valuesField(ruleEntry.Node, "match-interface", "match_interface"), scope, location(ruleEntry.Node, ruleEntry.Path))
			if value := uint32Field(ruleEntry.Node, "set-metric", "set_metric"); value > 0 {
				rule.SetMetric = &value
			}
			if value := uint32Field(ruleEntry.Node, "set-local-preference", "set_local_preference"); value > 0 {
				rule.SetLocalPreference = &value
			}
			rule.SetNextHop = scalarField(ruleEntry.Node, "set-ip-nexthop", "set_ip_nexthop")
			routeMap.Rules = append(routeMap.Rules, rule)
		}
		n.out.Routing.RouteMaps = append(n.out.Routing.RouteMaps, routeMap)
	}

	if section := childAlias(vdom, "router_bgp"); section != nil {
		n.out.Routing.BGP = &snapshot.RoutingProtocol{Enabled: true, RouterID: scalarField(section, "router-id", "router_id"), ASN: scalarField(section, "as"), Source: location(section, basePath+".router_bgp")}
	}
	if section := childAlias(vdom, "router_ospf"); section != nil {
		n.out.Routing.OSPF = &snapshot.RoutingProtocol{Enabled: true, RouterID: scalarField(section, "router-id", "router_id"), Source: location(section, basePath+".router_ospf")}
	}
	if section := childAlias(vdom, "router_rip"); section != nil {
		n.out.Routing.RIP = &snapshot.RoutingProtocol{Enabled: true, Source: location(section, basePath+".router_rip")}
	}
}

func (n *normalizer) parseDHCP(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	section := childAlias(vdom, "system_dhcp_server")
	for _, entry := range tableEntries(section, basePath+".system_dhcp_server") {
		id := objectID(scope, "dhcp-server", entry.Name)
		source := location(entry.Node, entry.Path)
		interfaceName := scalarField(entry.Node, "interface")
		server := snapshot.DHCPServer{
			ID: id, Scope: scope, Name: entry.Name, Enabled: enabledField(entry.Node, true, "status"),
			Interface: n.ref(id, "dhcp-interface", snapshot.ObjectKindInterface, interfaceName, scope, source),
			Gateway:   scalarField(entry.Node, "default-gateway", "default_gateway"), Netmask: scalarField(entry.Node, "netmask"),
			DNSServers: cleanValues([]string{scalarField(entry.Node, "dns-server1"), scalarField(entry.Node, "dns-server2"), scalarField(entry.Node, "dns-server3")}),
			NTPServers: valuesField(entry.Node, "ntp-server1", "ntp_server1"), DomainName: scalarField(entry.Node, "domain"),
			LeaseSeconds: uint64Field(entry.Node, "lease-time", "lease_time"), Source: source,
		}
		for _, rangeEntry := range tableEntries(childAlias(entry.Node, "ip-range", "ip_range"), entry.Path+".ip-range") {
			server.Ranges = append(server.Ranges, snapshot.DHCPRange{Start: scalarField(rangeEntry.Node, "start-ip", "start_ip"), End: scalarField(rangeEntry.Node, "end-ip", "end_ip")})
		}
		for _, reservationEntry := range tableEntries(childAlias(entry.Node, "reserved-address", "reserved_address"), entry.Path+".reserved-address") {
			server.Reservations = append(server.Reservations, snapshot.DHCPReservation{Name: reservationEntry.Name, MACAddress: scalarField(reservationEntry.Node, "mac"), IP: scalarField(reservationEntry.Node, "ip"), Description: scalarField(reservationEntry.Node, "description")})
		}
		if server.Gateway != "" && server.Netmask != "" {
			if prefix, ok := canonicalIPv4(server.Gateway, server.Netmask); ok {
				server.Network = prefix
			}
		}
		n.out.DHCP.Servers = append(n.out.DHCP.Servers, server)
	}

	for _, iface := range n.out.Interfaces {
		if iface.Scope.VDOM != scope.VDOM {
			continue
		}
		for _, relay := range []string{"dhcp-relay-ip", "dhcp_relay_ip"} {
			_ = relay
		}
	}
}

func (n *normalizer) parseSDWAN(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	section := childAlias(vdom, "system_sdwan", "system_virtual_wan_link")
	if section == nil {
		return
	}
	n.out.SDWAN.Enabled = enabledField(section, true, "status")

	zoneEntries := tableEntries(childAlias(section, "zone"), basePath+".system_sdwan.zone")
	for _, entry := range zoneEntries {
		id := objectID(scope, "sdwan-zone", entry.Name)
		n.register(scope, snapshot.ObjectKindSDWANZone, entry.Name, id)
		n.register(scope, snapshot.ObjectKindZone, entry.Name, id)
		n.out.SDWAN.Zones = append(n.out.SDWAN.Zones, snapshot.SDWANZone{ID: id, Scope: scope, Name: entry.Name, Source: location(entry.Node, entry.Path)})
	}
	// FortiGate's historical implicit name is often used even when no explicit zone table exists.
	if len(zoneEntries) == 0 {
		id := objectID(scope, "sdwan-zone", "virtual-wan-link")
		n.register(scope, snapshot.ObjectKindSDWANZone, "virtual-wan-link", id)
		n.register(scope, snapshot.ObjectKindZone, "virtual-wan-link", id)
		n.out.SDWAN.Zones = append(n.out.SDWAN.Zones, snapshot.SDWANZone{ID: id, Scope: scope, Name: "virtual-wan-link", Source: location(section, basePath+".system_sdwan")})
	}

	memberEntries := tableEntries(childAlias(section, "members"), basePath+".system_sdwan.members")
	for _, entry := range memberEntries {
		id := objectID(scope, "sdwan-member", entry.Name)
		n.register(scope, snapshot.ObjectKindSDWANMember, entry.Name, id)
	}
	for _, entry := range memberEntries {
		id := objectID(scope, "sdwan-member", entry.Name)
		source := location(entry.Node, entry.Path)
		interfaceName := scalarField(entry.Node, "interface")
		member := snapshot.SDWANMember{
			ID: id, Scope: scope, MemberID: uint32(sequenceNumber(entry)), Interface: n.ref(id, "sdwan-member-interface", snapshot.ObjectKindInterface, interfaceName, scope, source),
			Gateway: scalarField(entry.Node, "gateway"), Cost: uint32Field(entry.Node, "cost"), Priority: uint32Field(entry.Node, "priority"),
			Weight: uint32Field(entry.Node, "weight"), SourceIP: scalarField(entry.Node, "source"), Enabled: enabledField(entry.Node, true, "status"), Source: source,
		}
		if explicitID := uint32Field(entry.Node, "seq-num", "seq_num"); explicitID > 0 {
			member.MemberID = explicitID
		}
		zoneName := firstNonEmpty(scalarField(entry.Node, "zone"), "virtual-wan-link")
		zoneRef := n.ref(id, "sdwan-member-zone", snapshot.ObjectKindSDWANZone, zoneName, scope, source)
		member.Zone = &zoneRef
		n.out.SDWAN.Members = append(n.out.SDWAN.Members, member)
		for i := range n.out.SDWAN.Zones {
			if n.out.SDWAN.Zones[i].ID == zoneRef.ID {
				n.out.SDWAN.Zones[i].Members = append(n.out.SDWAN.Zones[i].Members, snapshot.ObjectRef{Kind: snapshot.ObjectKindSDWANMember, ID: id, VendorName: entry.Name, Scope: scope, Resolution: snapshot.ReferenceResolved, Source: source})
			}
		}
	}

	healthEntries := tableEntries(childAlias(section, "health-check", "health_check"), basePath+".system_sdwan.health-check")
	for _, entry := range healthEntries {
		id := objectID(scope, "sdwan-health-check", entry.Name)
		n.register(scope, snapshot.ObjectKindSDWANHealth, entry.Name, id)
	}
	for _, entry := range healthEntries {
		id := objectID(scope, "sdwan-health-check", entry.Name)
		source := location(entry.Node, entry.Path)
		health := snapshot.SDWANHealthCheck{
			ID: id, Scope: scope, Name: entry.Name, Server: valuesField(entry.Node, "server"), Protocol: scalarField(entry.Node, "protocol"),
			Members:          n.refs(id, "sdwan-health-member", snapshot.ObjectKindSDWANMember, valuesField(entry.Node, "members"), scope, source),
			LatencyThreshold: uint32Field(entry.Node, "latency-threshold", "latency_threshold"), JitterThreshold: uint32Field(entry.Node, "jitter-threshold", "jitter_threshold"),
			PacketLossThreshold: uint32Field(entry.Node, "packetloss-threshold", "packetloss_threshold"), Source: source,
		}
		n.out.SDWAN.HealthChecks = append(n.out.SDWAN.HealthChecks, health)
	}

	ruleEntries := tableEntries(childAlias(section, "service", "rules"), basePath+".system_sdwan.service")
	for _, entry := range ruleEntries {
		id := objectID(scope, "sdwan-rule", entry.Name)
		source := location(entry.Node, entry.Path)
		rule := snapshot.SDWANRule{
			ID: id, Scope: scope, Sequence: sequenceNumber(entry), Name: firstNonEmpty(scalarField(entry.Node, "name"), entry.Name), Enabled: enabledField(entry.Node, true, "status"),
			SourceAddresses:      n.refs(id, "sdwan-source-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "src", "srcaddr"), scope, source),
			DestinationAddresses: n.refs(id, "sdwan-destination-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "dst", "dstaddr"), scope, source),
			Services:             n.refs(id, "sdwan-service", snapshot.ObjectKindService, valuesField(entry.Node, "service"), scope, source),
			Applications:         valuesField(entry.Node, "application"), InternetServices: valuesField(entry.Node, "internet-service", "internet_service"),
			Strategy:         firstNonEmpty(scalarField(entry.Node, "mode"), scalarField(entry.Node, "strategy")),
			PreferredMembers: n.refs(id, "sdwan-preferred-member", snapshot.ObjectKindSDWANMember, valuesField(entry.Node, "priority-members", "priority_members"), scope, source),
			RequiredSLAs:     n.refs(id, "sdwan-required-sla", snapshot.ObjectKindSDWANHealth, valuesField(entry.Node, "health-check", "health_check"), scope, source),
			Source:           source,
		}
		if dscp := parseDSCP(scalarField(entry.Node, "dscp-forward", "dscp_forward", "dscp")); dscp != nil {
			rule.DSCP = &snapshot.DSCPMatch{Enabled: true, Values: []snapshot.DSCPValue{*dscp}}
		}
		n.out.SDWAN.Rules = append(n.out.SDWAN.Rules, rule)
	}
}

func sequenceNumber(entry namedNode) int {
	if parsed, err := strconv.Atoi(entry.Name); err == nil {
		return parsed
	}
	return entry.Sequence
}
