package fortigate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Iron-Signal-Systems/atlas/internal/ingest"
	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

type YAMLParser struct{}

func (YAMLParser) ID() string { return "firewall.fortigate.yaml.v1" }

func (YAMLParser) Probe(ctx context.Context, reader io.Reader) (ingest.Probe, error) {
	doc, err := ParseYAMLDocumentContext(ctx, reader)
	if err != nil {
		return ingest.Probe{}, err
	}
	if err := validateFortiGateYAMLRoot(doc.Root); err != nil {
		return ingest.Probe{}, err
	}
	version := detectFortiOSVersion(doc)
	return ingest.Probe{Vendor: "fortinet", Platform: "fortigate", Format: "fortios-yaml", Version: version}, nil
}

func (YAMLParser) Parse(ctx context.Context, reader io.Reader) (ingest.Result, error) {
	normalized, err := parseFortiGateYAMLWithLimits(ctx, reader, defaultYAMLDecodeLimits())
	if err != nil {
		return ingest.Result{}, err
	}
	return ingest.Result{
		Probe: ingest.Probe{
			Vendor:   "fortinet",
			Platform: "fortigate",
			Format:   "fortios-yaml",
			Version:  normalized.Source.FortiOSVersion,
		},
		Parsed:   normalized,
		Warnings: findingWarnings(normalized.Findings),
	}, nil
}

func ParseFortiGateYAML(reader io.Reader) (*snapshot.FirewallSnapshot, error) {
	normalized, _, err := parseFortiGateYAMLWithLayoutAndLimits(context.Background(), reader, defaultYAMLDecodeLimits())
	return normalized, err
}

func ParseFortiGateYAMLContext(ctx context.Context, reader io.Reader) (*snapshot.FirewallSnapshot, error) {
	normalized, _, err := parseFortiGateYAMLWithLayoutAndLimits(ctx, reader, defaultYAMLDecodeLimits())
	return normalized, err
}

// ParseFortiGateYAMLWithLayout performs one bounded decode and returns both the
// normalized snapshot and an upload-safe structural diagnostic derived from
// the same admitted document.
func ParseFortiGateYAMLWithLayout(reader io.Reader) (*snapshot.FirewallSnapshot, FortiGateYAMLLayout, error) {
	return parseFortiGateYAMLWithLayoutAndLimits(context.Background(), reader, defaultYAMLDecodeLimits())
}

func parseFortiGateYAMLWithLimits(
	ctx context.Context,
	reader io.Reader,
	limits yamlDecodeLimits,
) (*snapshot.FirewallSnapshot, error) {
	normalized, _, err := parseFortiGateYAMLWithLayoutAndLimits(ctx, reader, limits)
	return normalized, err
}

func parseFortiGateYAMLWithLayoutAndLimits(
	ctx context.Context,
	reader io.Reader,
	limits yamlDecodeLimits,
) (*snapshot.FirewallSnapshot, FortiGateYAMLLayout, error) {
	doc, err := parseYAMLDocumentWithLimits(ctx, reader, limits)
	if err != nil {
		return nil, FortiGateYAMLLayout{}, err
	}
	layout := DiagnoseFortiGateYAMLLayout(doc)
	if err := validateFortiGateYAMLRoot(doc.Root); err != nil {
		return nil, layout, err
	}
	normalized, err := NormalizeFortiGateYAML(doc)
	if err != nil {
		return nil, layout, err
	}
	if err := validateNormalizedSnapshotLimits(normalized, limits); err != nil {
		return nil, layout, err
	}
	return normalized, layout, nil
}

func validateFortiGateYAMLRoot(root *YAMLNode) error {
	if root == nil || root.Kind != YAMLMapping {
		return errors.New("FortiGate YAML root must be a mapping")
	}
	global := root.Child("global")
	if global != nil && global.Kind != YAMLMapping {
		return errors.New("FortiGate YAML global section must be a mapping")
	}
	vdom := root.Child("vdom")
	if vdom != nil && vdom.Kind != YAMLSequence && vdom.Kind != YAMLMapping {
		return errors.New("FortiGate YAML vdom section must be a sequence or mapping")
	}
	if fortiGateGlobalScope(root) == nil && len(fortiGateVDOMEntries(root)) == 0 {
		return errors.New("input does not appear to be a supported FortiGate YAML configuration layout")
	}
	return nil
}

var fortiOSVersionPattern = regexp.MustCompile(`(?:^|[^0-9])(\d+\.\d+\.\d+)(?:[^0-9]|$)`)

func detectFortiOSVersion(doc *YAMLDocument) string {
	if doc == nil {
		return ""
	}
	for _, comment := range doc.Comments {
		lower := strings.ToLower(comment)
		for _, marker := range []string{"config-version=", "fortios-version=", "version="} {
			if index := strings.Index(lower, marker); index >= 0 {
				value := strings.Trim(strings.TrimSpace(comment[index+len(marker):]), "# \t")
				if match := fortiOSVersionPattern.FindStringSubmatch(value); len(match) == 2 {
					return match[1]
				}
				return value
			}
		}
	}
	if node := doc.Root.At("global", "system_global", "version"); node != nil {
		return node.Scalar()
	}
	if global := fortiGateGlobalScope(doc.Root); global != nil {
		if node := childAlias(global, "system_global"); node != nil {
			return scalarField(node, "version")
		}
	}
	return ""
}

func findingWarnings(findings []snapshot.Finding) []string {
	warnings := make([]string, 0)
	for _, finding := range findings {
		if finding.Severity == "warning" || finding.Severity == "error" {
			warnings = append(warnings, fmt.Sprintf("%s: %s", finding.Title, finding.Detail))
		}
	}
	return warnings
}
