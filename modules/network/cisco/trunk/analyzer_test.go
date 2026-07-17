package trunk

import (
	"slices"
	"testing"

	"github.com/Iron-Signal-Systems/iron-atlas/modules/network/cisco/common"
)

func TestUnknownInterfaceDeniedEndpointAttribution(t *testing.T) {
	result := Classify(common.Interface{Name: "Gi1/0/1"})
	if result.EndpointAttribution || result.Class != common.ClassUnknown {
		t.Fatalf("unknown interface must fail closed: %#v", result)
	}
	if !slices.Contains(result.Reasons, "insufficient_access_evidence") {
		t.Fatalf("missing insufficient-evidence reason: %#v", result)
	}
}

func TestPositiveAccessEvidenceAllowsEndpointAttribution(t *testing.T) {
	tests := []struct {
		name  string
		iface common.Interface
	}{
		{name: "configured access", iface: common.Interface{Name: "Gi1/0/1", AdministrativeMode: "static access"}},
		{name: "operational access", iface: common.Interface{Name: "Gi1/0/2", OperationalMode: "static access"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Classify(test.iface)
			if !result.EndpointAttribution || result.Class != common.ClassAccessEndpoint {
				t.Fatalf("positive access evidence was not accepted: %#v", result)
			}
		})
	}
}

func TestTrunkEvidenceOverridesAccessEvidence(t *testing.T) {
	result := Classify(common.Interface{
		Name:               "Gi1/0/48",
		AdministrativeMode: "static access",
		OperationalMode:    "trunk",
	})
	if result.EndpointAttribution || result.Class != common.ClassInfrastructureTrunk {
		t.Fatalf("trunk evidence must deny endpoint attribution: %#v", result)
	}
}

func TestInfrastructureInterfacesDeniedEndpointAttribution(t *testing.T) {
	tests := []struct {
		name          string
		iface         common.Interface
		expectedClass common.InterfaceClass
	}{
		{name: "trunk", iface: common.Interface{AdministrativeMode: "trunk"}, expectedClass: common.ClassInfrastructureTrunk},
		{name: "routed", iface: common.Interface{Routed: true}, expectedClass: common.ClassRouted},
		{name: "port-channel member", iface: common.Interface{PortChannelID: "10"}, expectedClass: common.ClassPortChannelMember},
		{name: "stack link", iface: common.Interface{StackInterface: true}, expectedClass: common.ClassStack},
		{name: "fabric link", iface: common.Interface{FabricInterconnect: true}, expectedClass: common.ClassFabric},
		{name: "explicit exclusion", iface: common.Interface{ExplicitlyExcluded: true}, expectedClass: common.ClassExcluded},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Classify(test.iface)
			if result.EndpointAttribution || result.Class != test.expectedClass {
				t.Fatalf("infrastructure interface was not denied: %#v", result)
			}
		})
	}
}

func TestUnknownInterfaceDoesNotProduceTrunkFindings(t *testing.T) {
	if findings := Analyze(common.Interface{Name: "Gi1/0/1"}); len(findings) != 0 {
		t.Fatalf("unknown interface produced trunk findings: %#v", findings)
	}
}

func TestTrunkAnalysisRequiresDescriptionNeighborVLANAndSTP(t *testing.T) {
	findings := Analyze(common.Interface{Name: "Gi1/0/48", AdministrativeMode: "trunk"})
	if len(findings) != 4 {
		t.Fatalf("expected 4 baseline findings, got %d: %#v", len(findings), findings)
	}
}
