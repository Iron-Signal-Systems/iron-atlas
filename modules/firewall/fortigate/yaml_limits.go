package fortigate

import (
	"errors"
	"math"
)

// yamlDecodeLimits is intentionally internal. It defines the complete bounded
// admission policy used before untrusted FortiGate YAML reaches normalization.
type yamlDecodeLimits struct {
	MaxInputBytes             int64
	MaxNodes                  int
	MaxMappingDepth           int
	MaxSequenceDepth          int
	MaxCombinedDepth          int
	MaxMappingEntries         int
	MaxSequenceEntries        int
	MaxScalarBytes            int
	MaxKeyBytes               int
	MaxCompatibilityRewrites  int
	MaxCompatibilityFragments int
	MaxTotalMappingEntries    int
	MaxTotalSequenceEntries   int
	MaxNormalizedRecords      int
	MaxReferences             int
	MaxFindings               int
}

func defaultYAMLDecodeLimits() yamlDecodeLimits {
	return yamlDecodeLimits{
		MaxInputBytes:             64 * 1024 * 1024,
		MaxNodes:                  4_000_000,
		MaxMappingDepth:           128,
		MaxSequenceDepth:          128,
		MaxCombinedDepth:          192,
		MaxMappingEntries:         500_000,
		MaxSequenceEntries:        500_000,
		MaxScalarBytes:            8 * 1024 * 1024,
		MaxKeyBytes:               64 * 1024,
		MaxCompatibilityRewrites:  100_000,
		MaxCompatibilityFragments: 1_000_000,
		MaxTotalMappingEntries:    2_000_000,
		MaxTotalSequenceEntries:   2_000_000,
		MaxNormalizedRecords:      1_000_000,
		MaxReferences:             2_000_000,
		MaxFindings:               500_000,
	}
}

func (limits yamlDecodeLimits) validate() error {
	if limits.MaxInputBytes <= 0 || limits.MaxInputBytes == math.MaxInt64 {
		return errors.New("YAML MaxInputBytes must be between 1 and math.MaxInt64-1")
	}
	values := []int{
		limits.MaxNodes,
		limits.MaxMappingDepth,
		limits.MaxSequenceDepth,
		limits.MaxCombinedDepth,
		limits.MaxMappingEntries,
		limits.MaxSequenceEntries,
		limits.MaxScalarBytes,
		limits.MaxKeyBytes,
		limits.MaxCompatibilityRewrites,
		limits.MaxCompatibilityFragments,
		limits.MaxTotalMappingEntries,
		limits.MaxTotalSequenceEntries,
		limits.MaxNormalizedRecords,
		limits.MaxReferences,
		limits.MaxFindings,
	}
	for _, value := range values {
		if value <= 0 {
			return errors.New("all YAML admission limits must be positive")
		}
	}
	return nil
}
