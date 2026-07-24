package fortigate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func (n *normalizer) parseObjects(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	addressEntries := tableEntries(childAlias(vdom, "firewall_address"), basePath+".firewall_address")
	for _, entry := range addressEntries {
		id := objectID(scope, "address", entry.Name)
		n.register(scope, snapshot.ObjectKindAddress, entry.Name, id)
	}
	for _, entry := range addressEntries {
		id := objectID(scope, "address", entry.Name)
		source := location(entry.Node, entry.Path)
		obj := snapshot.AddressObject{
			ID: id, Scope: scope, Name: entry.Name, Kind: firstNonEmpty(scalarField(entry.Node, "type"), "subnet"),
			FQDNs: valuesField(entry.Node, "fqdn"), Wildcards: valuesField(entry.Node, "wildcard"), MACs: valuesField(entry.Node, "macaddr", "mac"),
			Countries: valuesField(entry.Node, "country"), Comments: scalarField(entry.Node, "comment", "comments"), Source: source,
		}
		if values := valuesField(entry.Node, "subnet"); len(values) > 0 {
			if prefix, _, _, ok := canonicalPrefix(values); ok {
				obj.Prefixes = append(obj.Prefixes, prefix)
			} else {
				obj.Prefixes = append(obj.Prefixes, strings.Join(values, " "))
			}
		}
		start := scalarField(entry.Node, "start-ip", "start_ip")
		end := scalarField(entry.Node, "end-ip", "end_ip")
		if start != "" || end != "" {
			obj.IPRanges = append(obj.IPRanges, snapshot.IPRange{Start: start, End: end})
		}
		if iface := scalarField(entry.Node, "associated-interface", "associated_interface"); iface != "" {
			ref := n.ref(id, "address-associated-interface", snapshot.ObjectKindInterfaceOrZone, iface, scope, source)
			obj.Interface = &ref
		}
		n.out.Objects.Addresses = append(n.out.Objects.Addresses, obj)
	}

	groupEntries := tableEntries(childAlias(vdom, "firewall_addrgrp"), basePath+".firewall_addrgrp")
	for _, entry := range groupEntries {
		id := objectID(scope, "address-group", entry.Name)
		n.register(scope, snapshot.ObjectKindAddressGroup, entry.Name, id)
	}
	for _, entry := range groupEntries {
		id := objectID(scope, "address-group", entry.Name)
		source := location(entry.Node, entry.Path)
		group := snapshot.AddressGroup{
			ID: id, Scope: scope, Name: entry.Name,
			Members:        n.refs(id, "address-group-member", snapshot.ObjectKindAddress, valuesField(entry.Node, "member"), scope, source),
			ExcludeMembers: n.refs(id, "address-group-exclude-member", snapshot.ObjectKindAddress, valuesField(entry.Node, "exclude-member", "exclude_member"), scope, source),
			Comments:       scalarField(entry.Node, "comment", "comments"), Source: source,
		}
		if len(group.Members) == 0 {
			n.addFinding("warning", "object", "Empty address group", fmt.Sprintf("Address group %s contains no members", entry.Name), id, source)
		}
		n.out.Objects.AddressGroups = append(n.out.Objects.AddressGroups, group)
	}

	serviceEntries := tableEntries(childAlias(vdom, "firewall_service_custom"), basePath+".firewall_service_custom")
	for _, entry := range serviceEntries {
		id := objectID(scope, "service", entry.Name)
		n.register(scope, snapshot.ObjectKindService, entry.Name, id)
	}
	for _, entry := range serviceEntries {
		id := objectID(scope, "service", entry.Name)
		obj := snapshot.ServiceObject{
			ID: id, Scope: scope, Name: entry.Name, Protocol: scalarField(entry.Node, "protocol"),
			TCPPorts:  parsePortRanges(valuesField(entry.Node, "tcp-portrange", "tcp_portrange")),
			UDPPorts:  parsePortRanges(valuesField(entry.Node, "udp-portrange", "udp_portrange")),
			SCTPPorts: parsePortRanges(valuesField(entry.Node, "sctp-portrange", "sctp_portrange")),
			ICMPTypes: valuesField(entry.Node, "icmptype", "icmp-type"), Comments: scalarField(entry.Node, "comment", "comments"),
			Source: location(entry.Node, entry.Path),
		}
		if value := uint32Field(entry.Node, "protocol-number", "protocol_number"); value <= 255 && value > 0 {
			v := uint8(value)
			obj.ProtocolNumber = &v
		}
		n.out.Objects.Services = append(n.out.Objects.Services, obj)
	}

	serviceGroupEntries := tableEntries(childAlias(vdom, "firewall_service_group"), basePath+".firewall_service_group")
	for _, entry := range serviceGroupEntries {
		id := objectID(scope, "service-group", entry.Name)
		n.register(scope, snapshot.ObjectKindServiceGroup, entry.Name, id)
	}
	for _, entry := range serviceGroupEntries {
		id := objectID(scope, "service-group", entry.Name)
		source := location(entry.Node, entry.Path)
		group := snapshot.ServiceGroup{ID: id, Scope: scope, Name: entry.Name, Members: n.refs(id, "service-group-member", snapshot.ObjectKindService, valuesField(entry.Node, "member"), scope, source), Comments: scalarField(entry.Node, "comment", "comments"), Source: source}
		n.out.Objects.ServiceGroups = append(n.out.Objects.ServiceGroups, group)
	}

	for _, spec := range []struct{ section, kind string }{{"firewall_schedule_onetime", "one-time"}, {"firewall_schedule_recurring", "recurring"}} {
		entries := tableEntries(childAlias(vdom, spec.section), basePath+"."+spec.section)
		for _, entry := range entries {
			id := objectID(scope, "schedule", entry.Name)
			n.register(scope, snapshot.ObjectKindSchedule, entry.Name, id)
			n.out.Objects.Schedules = append(n.out.Objects.Schedules, snapshot.Schedule{ID: id, Scope: scope, Name: entry.Name, Kind: spec.kind, Start: scalarField(entry.Node, "start"), End: scalarField(entry.Node, "end"), Days: valuesField(entry.Node, "day"), Source: location(entry.Node, entry.Path)})
		}
	}
}

func (n *normalizer) parseNAT(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	vipEntries := tableEntries(childAlias(vdom, "firewall_vip"), basePath+".firewall_vip")
	for _, entry := range vipEntries {
		id := objectID(scope, "virtual-ip", entry.Name)
		n.register(scope, snapshot.ObjectKindVIP, entry.Name, id)
	}
	for _, entry := range vipEntries {
		id := objectID(scope, "virtual-ip", entry.Name)
		source := location(entry.Node, entry.Path)
		vip := snapshot.VirtualIP{
			ID: id, Scope: scope, Name: entry.Name, Enabled: enabledField(entry.Node, true, "status"),
			ExternalAddress: scalarField(entry.Node, "extip"), MappedAddresses: valuesField(entry.Node, "mappedip"), Protocol: scalarField(entry.Node, "protocol"),
			ExternalPorts: parsePortRanges(valuesField(entry.Node, "extport")), MappedPorts: parsePortRanges(valuesField(entry.Node, "mappedport")),
			PortForwarding: enabledField(entry.Node, false, "portforward"), SourceFilters: valuesField(entry.Node, "src-filter", "src_filter"),
			Comments: scalarField(entry.Node, "comment", "comments"), Source: source,
		}
		if iface := scalarField(entry.Node, "extintf"); iface != "" {
			ref := n.ref(id, "vip-external-interface", snapshot.ObjectKindInterfaceOrZone, iface, scope, source)
			vip.ExternalInterface = &ref
		}
		n.out.NAT.VirtualIPs = append(n.out.NAT.VirtualIPs, vip)
	}

	vipGroupEntries := tableEntries(childAlias(vdom, "firewall_vipgrp"), basePath+".firewall_vipgrp")
	for _, entry := range vipGroupEntries {
		id := objectID(scope, "virtual-ip-group", entry.Name)
		members := n.refs(id, "vip-group-member", snapshot.ObjectKindVIP, valuesField(entry.Node, "member"), scope, location(entry.Node, entry.Path))
		n.out.NAT.VirtualIPGroups = append(n.out.NAT.VirtualIPGroups, snapshot.VirtualIPGroup{ID: id, Scope: scope, Name: entry.Name, Members: members, Source: location(entry.Node, entry.Path)})
	}

	poolEntries := tableEntries(childAlias(vdom, "firewall_ippool"), basePath+".firewall_ippool")
	for _, entry := range poolEntries {
		id := objectID(scope, "ip-pool", entry.Name)
		n.register(scope, snapshot.ObjectKindIPPool, entry.Name, id)
		n.out.NAT.IPPools = append(n.out.NAT.IPPools, snapshot.IPPool{ID: id, Scope: scope, Name: entry.Name, StartIP: scalarField(entry.Node, "startip"), EndIP: scalarField(entry.Node, "endip"), Type: scalarField(entry.Node, "type"), Source: location(entry.Node, entry.Path)})
	}

	for _, entry := range tableEntries(childAlias(vdom, "firewall_central_snat_map"), basePath+".firewall_central_snat_map") {
		id := objectID(scope, "central-snat", entry.Name)
		source := location(entry.Node, entry.Path)
		rule := snapshot.CentralSNATRule{
			ID: id, Scope: scope, Sequence: sequenceNumber(entry),
			SourceInterfaces:      n.refs(id, "central-snat-source-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "srcintf"), scope, source),
			DestinationInterfaces: n.refs(id, "central-snat-destination-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "dstintf"), scope, source),
			SourceAddresses:       n.refs(id, "central-snat-source-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "orig-addr", "orig_addr"), scope, source),
			DestinationAddresses:  n.refs(id, "central-snat-destination-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "dst-addr", "dst_addr"), scope, source),
			NATEnabled:            enabledField(entry.Node, true, "nat"), Source: source,
		}
		if pool := scalarField(entry.Node, "nat-ippool", "nat_ippool"); pool != "" {
			ref := n.ref(id, "central-snat-ip-pool", snapshot.ObjectKindIPPool, pool, scope, source)
			rule.IPPool = &ref
		}
		n.out.NAT.CentralSNAT = append(n.out.NAT.CentralSNAT, rule)
	}
}

func (n *normalizer) parseQoS(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	shaperEntries := tableEntries(childAlias(vdom, "firewall_shaper_traffic_shaper"), basePath+".firewall_shaper_traffic_shaper")
	for _, entry := range shaperEntries {
		id := objectID(scope, "traffic-shaper", entry.Name)
		n.register(scope, snapshot.ObjectKindTrafficShaper, entry.Name, id)
	}
	for _, entry := range shaperEntries {
		id := objectID(scope, "traffic-shaper", entry.Name)
		shaper := snapshot.TrafficShaper{
			ID: id, Scope: scope, Name: entry.Name,
			GuaranteedBandwidth: parseBandwidth(scalarField(entry.Node, "guaranteed-bandwidth", "guaranteed_bandwidth"), scalarField(entry.Node, "bandwidth-unit", "bandwidth_unit")),
			MaximumBandwidth:    parseBandwidth(scalarField(entry.Node, "maximum-bandwidth", "maximum_bandwidth"), scalarField(entry.Node, "bandwidth-unit", "bandwidth_unit")),
			BurstBytes:          uint64Field(entry.Node, "burst-in-msec", "burst_in_msec"), Priority: scalarField(entry.Node, "priority"), Mode: scalarField(entry.Node, "mode"),
			PerPolicy: enabledField(entry.Node, false, "per-policy", "per_policy"), Comments: scalarField(entry.Node, "comments", "comment"), Source: location(entry.Node, entry.Path),
		}
		if dscp := parseDSCP(scalarField(entry.Node, "dscp-marking", "dscp_marking")); dscp != nil {
			shaper.DSCPRemark = dscp
		}
		n.out.QoS.TrafficShapers = append(n.out.QoS.TrafficShapers, shaper)
	}

	perIPEntries := tableEntries(childAlias(vdom, "firewall_shaper_per_ip_shaper"), basePath+".firewall_shaper_per_ip_shaper")
	for _, entry := range perIPEntries {
		id := objectID(scope, "per-ip-shaper", entry.Name)
		n.register(scope, snapshot.ObjectKindPerIPShaper, entry.Name, id)
		n.out.QoS.PerIPShapers = append(n.out.QoS.PerIPShapers, snapshot.PerIPShaper{ID: id, Scope: scope, Name: entry.Name, MaximumBandwidth: parseBandwidth(scalarField(entry.Node, "max-bandwidth", "max_bandwidth"), scalarField(entry.Node, "bandwidth-unit", "bandwidth_unit")), ConcurrentLimit: uint32Field(entry.Node, "max-concurrent-session", "max_concurrent_session"), Source: location(entry.Node, entry.Path)})
	}

	for _, entry := range tableEntries(childAlias(vdom, "firewall_shaping_policy"), basePath+".firewall_shaping_policy") {
		id := objectID(scope, "shaping-policy", entry.Name)
		source := location(entry.Node, entry.Path)
		policy := snapshot.ShapingPolicy{
			ID: id, Scope: scope, Sequence: sequenceNumber(entry), Enabled: enabledField(entry.Node, true, "status"),
			SourceInterfaces:      n.refs(id, "shaping-source-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "srcintf"), scope, source),
			DestinationInterfaces: n.refs(id, "shaping-destination-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "dstintf"), scope, source),
			SourceAddresses:       n.refs(id, "shaping-source-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "srcaddr"), scope, source),
			DestinationAddresses:  n.refs(id, "shaping-destination-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "dstaddr"), scope, source),
			Services:              n.refs(id, "shaping-service", snapshot.ObjectKindService, valuesField(entry.Node, "service"), scope, source),
			Applications:          valuesField(entry.Node, "application"), Source: source,
		}
		if name := scalarField(entry.Node, "traffic-shaper", "traffic_shaper"); name != "" {
			ref := n.ref(id, "shaping-forward-shaper", snapshot.ObjectKindTrafficShaper, name, scope, source)
			policy.ForwardShaper = &ref
		}
		if name := scalarField(entry.Node, "traffic-shaper-reverse", "traffic_shaper_reverse"); name != "" {
			ref := n.ref(id, "shaping-reverse-shaper", snapshot.ObjectKindTrafficShaper, name, scope, source)
			policy.ReverseShaper = &ref
		}
		if name := scalarField(entry.Node, "per-ip-shaper", "per_ip_shaper"); name != "" {
			ref := n.ref(id, "shaping-per-ip-shaper", snapshot.ObjectKindPerIPShaper, name, scope, source)
			policy.PerIPShaper = &ref
		}
		if dscp := parseDSCP(scalarField(entry.Node, "dscp")); dscp != nil {
			policy.DSCP = &snapshot.DSCPMatch{Enabled: true, Values: []snapshot.DSCPValue{*dscp}}
		}
		n.out.QoS.ShapingPolicies = append(n.out.QoS.ShapingPolicies, policy)
	}
}

func parseBandwidth(value, unit string) snapshot.Bandwidth {
	value = strings.TrimSpace(value)
	unit = strings.ToLower(strings.TrimSpace(unit))
	if value == "" {
		return snapshot.Bandwidth{}
	}
	numeric, _ := strconv.ParseUint(value, 10, 64)
	multiplier := uint64(1000)
	switch unit {
	case "bps", "bit", "bits":
		multiplier = 1
	case "kbps", "kbit", "kbits", "":
		multiplier = 1000
	case "mbps", "mbit", "mbits":
		multiplier = 1000 * 1000
	case "gbps", "gbit", "gbits":
		multiplier = 1000 * 1000 * 1000
	}
	return snapshot.Bandwidth{BitsPerSecond: numeric * multiplier, OriginalValue: value, OriginalUnit: unit}
}

func (n *normalizer) parseVPN(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	phase1Entries := tableEntries(childAlias(vdom, "vpn_ipsec_phase1_interface"), basePath+".vpn_ipsec_phase1_interface")
	for _, entry := range phase1Entries {
		id := objectID(scope, "vpn", entry.Name)
		n.register(scope, snapshot.ObjectKindVPN, entry.Name, id)
		source := location(entry.Node, entry.Path)
		vpn := snapshot.VPN{ID: id, Scope: scope, Name: entry.Name, Kind: "ipsec-interface", RemoteGateway: scalarField(entry.Node, "remote-gw", "remote_gw"), Enabled: enabledField(entry.Node, true, "status"), Source: source}
		if iface := scalarField(entry.Node, "interface"); iface != "" {
			ref := n.ref(id, "vpn-interface", snapshot.ObjectKindInterface, iface, scope, source)
			vpn.Interface = &ref
		}
		n.out.VPNs = append(n.out.VPNs, vpn)
	}
}

func (n *normalizer) parsePolicies(vdom *YAMLNode, scope snapshot.ObjectScope, basePath string) {
	entries := tableEntries(childAlias(vdom, "firewall_policy"), basePath+".firewall_policy")
	for _, entry := range entries {
		id := objectID(scope, "policy", entry.Name)
		source := location(entry.Node, entry.Path)
		policy := snapshot.Policy{
			ID: id, Scope: scope, VendorID: entry.Name, UUID: scalarField(entry.Node, "uuid"), Name: firstNonEmpty(scalarField(entry.Node, "name"), entry.Name),
			Comments: scalarField(entry.Node, "comments", "comment"), Enabled: enabledField(entry.Node, true, "status"), Sequence: sequenceNumber(entry),
			Action: firstNonEmpty(scalarField(entry.Node, "action"), "deny"), Source: source,
		}
		policy.Match.SourceInterfaces = n.refs(id, "policy-source-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "srcintf"), scope, source)
		policy.Match.DestinationInterfaces = n.refs(id, "policy-destination-interface", snapshot.ObjectKindInterfaceOrZone, valuesField(entry.Node, "dstintf"), scope, source)
		policy.Match.SourceAddresses = n.refs(id, "policy-source-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "srcaddr"), scope, source)
		policy.Match.DestinationAddresses = n.refs(id, "policy-destination-address", snapshot.ObjectKindAddress, valuesField(entry.Node, "dstaddr"), scope, source)
		policy.Match.SourceAddressesNegated = enabledField(entry.Node, false, "srcaddr-negate", "srcaddr_negate")
		policy.Match.DestinationAddressesNegated = enabledField(entry.Node, false, "dstaddr-negate", "dstaddr_negate")
		policy.Match.Services = n.refs(id, "policy-service", snapshot.ObjectKindService, valuesField(entry.Node, "service"), scope, source)
		policy.Match.ServicesNegated = enabledField(entry.Node, false, "service-negate", "service_negate")
		policy.Match.Schedules = n.refs(id, "policy-schedule", snapshot.ObjectKindSchedule, valuesField(entry.Node, "schedule"), scope, source)
		policy.Match.Users = n.refs(id, "policy-user", snapshot.ObjectKindUser, valuesField(entry.Node, "users"), scope, source)
		policy.Match.UserGroups = n.refs(id, "policy-user-group", snapshot.ObjectKindUserGroup, valuesField(entry.Node, "groups"), scope, source)
		policy.Match.InternetServices = valuesField(entry.Node, "internet-service-name", "internet_service_name")
		policy.Match.Applications = valuesField(entry.Node, "application-list", "application_list")
		if dscp := parseDSCP(scalarField(entry.Node, "diffservcode", "dscp")); dscp != nil {
			policy.Match.DSCP = &snapshot.DSCPMatch{Enabled: enabledField(entry.Node, true, "diffserv"), Values: []snapshot.DSCPValue{*dscp}}
		}

		policy.NAT.Enabled = enabledField(entry.Node, false, "nat")
		policy.NAT.SourceNAT = policy.NAT.Enabled
		policy.NAT.FixedPort = enabledField(entry.Node, false, "fixedport")
		policy.NAT.PreserveSourcePort = policy.NAT.FixedPort
		policy.NAT.IPPools = n.refs(id, "policy-ip-pool", snapshot.ObjectKindIPPool, valuesField(entry.Node, "poolname"), scope, source)
		for _, ref := range policy.Match.DestinationAddresses {
			if ref.Kind == snapshot.ObjectKindVIP {
				policy.NAT.VIPs = append(policy.NAT.VIPs, ref)
			}
		}

		if name := scalarField(entry.Node, "traffic-shaper", "traffic_shaper"); name != "" {
			ref := n.ref(id, "policy-forward-shaper", snapshot.ObjectKindTrafficShaper, name, scope, source)
			policy.QoS.ForwardShaper = &ref
		}
		if name := scalarField(entry.Node, "traffic-shaper-reverse", "traffic_shaper_reverse"); name != "" {
			ref := n.ref(id, "policy-reverse-shaper", snapshot.ObjectKindTrafficShaper, name, scope, source)
			policy.QoS.ReverseShaper = &ref
		}
		if name := scalarField(entry.Node, "per-ip-shaper", "per_ip_shaper"); name != "" {
			ref := n.ref(id, "policy-per-ip-shaper", snapshot.ObjectKindPerIPShaper, name, scope, source)
			policy.QoS.PerIPShaper = &ref
		}
		if value := parseDSCP(scalarField(entry.Node, "diffservcode-forward", "diffservcode_forward")); value != nil {
			policy.QoS.ForwardDSCP = &snapshot.DSCPRewrite{Enabled: true, Value: value}
		}
		if value := parseDSCP(scalarField(entry.Node, "diffservcode-rev", "diffservcode_rev")); value != nil {
			policy.QoS.ReverseDSCP = &snapshot.DSCPRewrite{Enabled: true, Value: value}
		}

		policy.Security.InspectionMode = scalarField(entry.Node, "inspection-mode", "inspection_mode")
		policy.Security.AVProfile = scalarField(entry.Node, "av-profile", "av_profile")
		policy.Security.WebFilterProfile = scalarField(entry.Node, "webfilter-profile", "webfilter_profile")
		policy.Security.IPSProfile = scalarField(entry.Node, "ips-sensor", "ips_sensor")
		policy.Security.ApplicationList = scalarField(entry.Node, "application-list", "application_list")
		policy.Security.SSLSSHProfile = scalarField(entry.Node, "ssl-ssh-profile", "ssl_ssh_profile")
		logTraffic := strings.ToLower(scalarField(entry.Node, "logtraffic"))
		policy.Logging.Enabled = logTraffic != "" && logTraffic != "disable"
		policy.Logging.AllSessions = logTraffic == "all"
		policy.Logging.UTM = logTraffic == "utm"
		policy.Logging.Start = enabledField(entry.Node, false, "logtraffic-start", "logtraffic_start")

		if !policy.Enabled {
			n.addFinding("info", "policy", "Disabled firewall policy", fmt.Sprintf("Policy %s (%s) is disabled", policy.Name, policy.VendorID), id, source)
		}
		n.out.Policies = append(n.out.Policies, policy)
	}
}
