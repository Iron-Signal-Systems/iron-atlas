package fortigate

import (
	"errors"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func validateNormalizedSnapshotLimits(value *snapshot.FirewallSnapshot, limits yamlDecodeLimits) error {
	if value == nil {
		return errors.New("FortiGate normalization returned no snapshot")
	}

	if normalizedRecordCount(value) > limits.MaxNormalizedRecords {
		return errors.New("FortiGate normalization rejected: MaxNormalizedRecords limit exceeded")
	}

	references := len(value.References.Edges) + len(value.References.BuiltIns) + len(value.References.Unresolved)
	if references > limits.MaxReferences {
		return errors.New("FortiGate normalization rejected: MaxReferences limit exceeded")
	}
	if len(value.Findings) > limits.MaxFindings {
		return errors.New("FortiGate normalization rejected: MaxFindings limit exceeded")
	}
	return nil
}

func normalizedRecordCount(value *snapshot.FirewallSnapshot) int {
	return snapshot.TotalNormalizedRecords(value)
}
