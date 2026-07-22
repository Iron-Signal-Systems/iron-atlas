package trunk

import (
	"strings"

	"github.com/Iron-Signal-Systems/atlas/modules/network/cisco/common"
)

func Classify(iface common.Interface) common.Classification {
	result := common.Classification{Class: common.ClassUnknown, EndpointAttribution: true}
	modeAdmin := strings.ToLower(iface.AdministrativeMode)
	modeOperational := strings.ToLower(iface.OperationalMode)

	switch {
	case iface.ExplicitlyExcluded:
		result.Class = common.ClassExcluded
		result.EndpointAttribution = false
		result.Reasons = append(result.Reasons, "explicitly_excluded")
	case iface.StackInterface:
		result.Class = common.ClassStack
		result.EndpointAttribution = false
		result.Reasons = append(result.Reasons, "stack_interface")
	case iface.FabricInterconnect:
		result.Class = common.ClassFabric
		result.EndpointAttribution = false
		result.Reasons = append(result.Reasons, "fabric_interconnect")
	case iface.PortChannelID != "":
		result.Class = common.ClassPortChannelMember
		result.EndpointAttribution = false
		result.Reasons = append(result.Reasons, "port_channel_member")
	case iface.Routed:
		result.Class = common.ClassRouted
		result.EndpointAttribution = false
		result.Reasons = append(result.Reasons, "routed_interface")
	case strings.Contains(modeAdmin, "trunk") || strings.Contains(modeOperational, "trunk"):
		result.Class = common.ClassInfrastructureTrunk
		result.EndpointAttribution = false
		if strings.Contains(modeAdmin, "trunk") {
			result.Reasons = append(result.Reasons, "configured_trunk")
		}
		if strings.Contains(modeOperational, "trunk") {
			result.Reasons = append(result.Reasons, "operational_trunk")
		}
	default:
		result.Class = common.ClassAccessEndpoint
	}
	return result
}

type Finding struct{ Code, Severity, Message string }

func Analyze(iface common.Interface) []Finding {
	classification := Classify(iface)
	if classification.EndpointAttribution {
		return nil
	}
	findings := make([]Finding, 0)
	if strings.TrimSpace(iface.Description) == "" {
		findings = append(findings, Finding{Code: "TRUNK_DESCRIPTION_MISSING", Severity: "low", Message: "infrastructure interface has no description"})
	}
	if iface.CDPNeighbor == nil && iface.LLDPNeighbor == nil {
		findings = append(findings, Finding{Code: "TRUNK_NEIGHBOR_NOT_OBSERVED", Severity: "informational", Message: "no CDP or LLDP neighbor was observed"})
	}
	if classification.Class == common.ClassInfrastructureTrunk && len(iface.AllowedVLANs) == 0 {
		findings = append(findings, Finding{Code: "TRUNK_ALLOWED_VLANS_UNKNOWN", Severity: "moderate", Message: "effective allowed VLAN set was not established"})
	}
	if len(iface.SpanningTree) == 0 && classification.Class == common.ClassInfrastructureTrunk {
		findings = append(findings, Finding{Code: "TRUNK_STP_STATE_UNKNOWN", Severity: "moderate", Message: "spanning-tree state was not collected"})
	}
	return findings
}
