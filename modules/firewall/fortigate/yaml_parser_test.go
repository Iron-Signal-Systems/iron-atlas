package fortigate

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func TestParseYAMLDocumentPreservesFortiGateHierarchy(t *testing.T) {
	input := `vdom:
  - root:
      system_interface:
        - wan1:
            ip: [192.0.2.2, 255.255.255.0]
            role: wan
`
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	vdom := doc.Root.Child("vdom")
	if vdom == nil || len(vdom.Seq) != 1 {
		t.Fatalf("unexpected vdom parse: %#v", vdom)
	}
	root := vdom.Seq[0].Child("root")
	interfaces := childAlias(root, "system_interface")
	entries := tableEntries(interfaces, "$.vdom[0].root.system_interface")
	if len(entries) != 1 || entries[0].Name != "wan1" {
		t.Fatalf("unexpected interface entries: %#v", entries)
	}
	if got := valuesField(entries[0].Node, "ip"); len(got) != 2 || got[0] != "192.0.2.2" {
		t.Fatalf("unexpected interface address: %#v", got)
	}
}

func TestParseFortiGateYAMLNormalizesRequestedDomains(t *testing.T) {
	file, err := os.Open("testdata/fortigate-sanitized.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	got, err := ParseFortiGateYAML(file)
	if err != nil {
		t.Fatal(err)
	}
	if got.Device.Hostname != "atlas-fgt-test" || got.Source.FortiOSVersion == "" {
		t.Fatalf("unexpected device metadata: %#v %#v", got.Device, got.Source)
	}
	counts := []struct {
		name string
		got  int
		want int
	}{
		{name: "domains", got: len(got.Domains), want: 1},
		{name: "interfaces", got: len(got.Interfaces), want: 4},
		{name: "VLANs", got: len(got.VLANs), want: 1},
		{name: "subnets", got: len(got.Subnets), want: 3},
		{name: "DHCP servers", got: len(got.DHCP.Servers), want: 1},
		{name: "address objects", got: len(got.Objects.Addresses), want: 3},
		{name: "address groups", got: len(got.Objects.AddressGroups), want: 1},
		{name: "service objects", got: len(got.Objects.Services), want: 2},
		{name: "service groups", got: len(got.Objects.ServiceGroups), want: 1},
		{name: "policies", got: len(got.Policies), want: 2},
		{name: "static routes", got: len(got.Routing.Routes), want: 2},
		{name: "policy routes", got: len(got.Routing.PolicyRoutes), want: 0},
		{name: "SD-WAN zones", got: len(got.SDWAN.Zones), want: 1},
		{name: "SD-WAN members", got: len(got.SDWAN.Members), want: 2},
		{name: "SD-WAN health checks", got: len(got.SDWAN.HealthChecks), want: 1},
		{name: "SD-WAN rules", got: len(got.SDWAN.Rules), want: 1},
		{name: "virtual IPs", got: len(got.NAT.VirtualIPs), want: 1},
		{name: "IP pools", got: len(got.NAT.IPPools), want: 1},
		{name: "VPNs", got: len(got.VPNs), want: 1},
		{name: "traffic shapers", got: len(got.QoS.TrafficShapers), want: 1},
		{name: "reference edges", got: len(got.References.Edges), want: 38},
		{name: "built-in references", got: len(got.References.BuiltIns), want: 6},
		{name: "unresolved references", got: len(got.References.Unresolved), want: 0},
		{name: "findings", got: len(got.Findings), want: 0},
	}
	for _, count := range counts {
		if count.got != count.want {
			t.Errorf("unexpected %s count: got %d want %d", count.name, count.got, count.want)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
	if len(got.Interfaces) != 4 || len(got.VLANs) != 1 || len(got.Subnets) != 3 {
		t.Fatalf("unexpected network counts: interfaces=%d vlans=%d subnets=%d", len(got.Interfaces), len(got.VLANs), len(got.Subnets))
	}
	if got.VLANs[0].VLANID != 120 {
		t.Fatalf("unexpected VLAN: %#v", got.VLANs[0])
	}
	var usersInterface *snapshot.Interface
	for i := range got.Interfaces {
		if got.Interfaces[i].Name == "users.120" {
			usersInterface = &got.Interfaces[i]
			break
		}
	}
	if usersInterface == nil || usersInterface.VLAN == nil || usersInterface.Addresses4[0].Prefix != "10.120.0.0/24" {
		t.Fatalf("unexpected VLAN interface: %#v", usersInterface)
	}
	if len(got.DHCP.Servers) != 1 || got.DHCP.Servers[0].Network != "10.120.0.0/24" || len(got.DHCP.Servers[0].Reservations) != 1 {
		t.Fatalf("unexpected DHCP parse: %#v", got.DHCP)
	}
	if len(got.Routing.Routes) != 2 || got.Routing.Routes[0].Destination != "0.0.0.0/0" {
		t.Fatalf("unexpected route parse: %#v", got.Routing.Routes)
	}
	if len(got.SDWAN.Zones) != 1 || len(got.SDWAN.Members) != 2 || len(got.SDWAN.HealthChecks) != 1 || len(got.SDWAN.Rules) != 1 {
		t.Fatalf("unexpected SD-WAN parse: %#v", got.SDWAN)
	}
	if len(got.Policies) != 2 || len(got.Policies[0].Match.SourceInterfaces) != 1 {
		t.Fatalf("unexpected policy parse: %#v", got.Policies)
	}
	if got.Policies[0].Match.DSCP == nil || got.Policies[0].Match.DSCP.Values[0].CanonicalName != "AF31" {
		t.Fatalf("unexpected policy DSCP: %#v", got.Policies[0].Match.DSCP)
	}
	if got.Policies[0].QoS.ForwardDSCP == nil || got.Policies[0].QoS.ForwardDSCP.Value.CanonicalName != "EF" {
		t.Fatalf("unexpected DSCP rewrite: %#v", got.Policies[0].QoS)
	}
	if len(got.NAT.VirtualIPs) != 1 || len(got.NAT.IPPools) != 1 {
		t.Fatalf("unexpected NAT parse: %#v", got.NAT)
	}
	if len(got.QoS.TrafficShapers) != 1 || got.QoS.TrafficShapers[0].MaximumBandwidth.BitsPerSecond != 50_000_000 {
		t.Fatalf("unexpected QoS parse: %#v", got.QoS)
	}
	if len(got.VPNs) != 1 {
		t.Fatalf("unexpected VPN parse: %#v", got.VPNs)
	}
	if len(got.References.Unresolved) != 0 {
		t.Fatalf("expected all fixture references to resolve, got %#v", got.References.Unresolved)
	}
}

func TestYAMLParserImplementsIngestContract(t *testing.T) {
	file, err := os.Open("testdata/fortigate-sanitized.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := (YAMLParser{}).Parse(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	if result.Probe.Format != "fortios-yaml" {
		t.Fatalf("unexpected probe: %#v", result.Probe)
	}
	if _, ok := result.Parsed.(*snapshot.FirewallSnapshot); !ok {
		t.Fatalf("unexpected parsed type %T", result.Parsed)
	}
}

func TestUnresolvedReferencesBecomeFindings(t *testing.T) {
	input := `vdom:
  - root:
      firewall_policy:
        - 1:
            srcintf: [missing-interface]
            dstintf: [all]
            srcaddr: [all]
            dstaddr: [all]
            service: [ALL]
            schedule: always
            action: accept
`
	got, err := ParseFortiGateYAML(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.References.Unresolved) != 1 {
		t.Fatalf("expected one unresolved reference, got %#v", got.References.Unresolved)
	}
	if got.References.Unresolved[0].Resolution != snapshot.ReferenceUnresolved {
		t.Fatalf("expected explicit unresolved classification, got %#v", got.References.Unresolved[0])
	}
	if len(got.Findings) == 0 || got.Findings[0].Title != "Unresolved object reference" {
		t.Fatalf("expected unresolved reference finding, got %#v", got.Findings)
	}
}

func TestBoundedParserRejectsUnsupportedYAMLFeatures(t *testing.T) {
	_, err := ParseYAMLDocument(strings.NewReader("global: &anchor\n  system_global: {}\n"))
	if err == nil {
		t.Fatal("expected anchor rejection")
	}
}
