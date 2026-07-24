package trunk

import (
	"testing"

	"github.com/Iron-Signal-Systems/atlas/modules/network/cisco/common"
)

func TestTrunkExcludedFromEndpointAttribution(t *testing.T) {
	result := Classify(common.Interface{Name: "Gi1/0/48", AdministrativeMode: "trunk", OperationalMode: "trunk", MACAddresses: []string{"0011.2233.4455"}})
	if result.EndpointAttribution {
		t.Fatal("trunk MAC entries must not be used for local endpoint attribution")
	}
	if result.Class != common.ClassInfrastructureTrunk {
		t.Fatalf("unexpected class: %s", result.Class)
	}
}

func TestPortChannelMemberExcluded(t *testing.T) {
	result := Classify(common.Interface{Name: "Te1/1/1", PortChannelID: "10"})
	if result.EndpointAttribution || result.Class != common.ClassPortChannelMember {
		t.Fatalf("unexpected classification: %#v", result)
	}
}

func TestTrunkAnalysisRequiresDescriptionNeighborVLANAndSTP(t *testing.T) {
	findings := Analyze(common.Interface{Name: "Gi1/0/48", AdministrativeMode: "trunk"})
	if len(findings) != 4 {
		t.Fatalf("expected 4 baseline findings, got %d: %#v", len(findings), findings)
	}
}
