package fortigate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/ingest"
	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/snapshot"
)

type YAMLParser struct{}

func (YAMLParser) ID() string { return "firewall.fortigate.yaml.v1" }

func (YAMLParser) Probe(_ context.Context, reader io.Reader) (ingest.Probe, error) {
	data, err := io.ReadAll(io.LimitReader(reader, 2*1024*1024))
	if err != nil {
		return ingest.Probe{}, err
	}
	doc, err := ParseYAMLDocument(bytes.NewReader(data))
	if err != nil {
		return ingest.Probe{}, err
	}
	if err := validateFortiGateYAMLRoot(doc.Root); err != nil {
		return ingest.Probe{}, err
	}
	version := detectFortiOSVersion(doc)
	return ingest.Probe{Vendor: "fortinet", Platform: "fortigate", Format: "fortios-yaml", Version: version}, nil
}

func (YAMLParser) Parse(_ context.Context, reader io.Reader) (ingest.Result, error) {
	doc, err := ParseYAMLDocument(reader)
	if err != nil {
		return ingest.Result{}, err
	}
	if err := validateFortiGateYAMLRoot(doc.Root); err != nil {
		return ingest.Result{}, err
	}
	normalized, err := NormalizeFortiGateYAML(doc)
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
	doc, err := ParseYAMLDocument(reader)
	if err != nil {
		return nil, err
	}
	if err := validateFortiGateYAMLRoot(doc.Root); err != nil {
		return nil, err
	}
	return NormalizeFortiGateYAML(doc)
}

func validateFortiGateYAMLRoot(root *YAMLNode) error {
	if root == nil || root.Kind != YAMLMapping {
		return errors.New("FortiGate YAML root must be a mapping")
	}
	if root.Child("global") == nil && root.Child("vdom") == nil {
		return errors.New("input does not appear to be a FortiGate YAML configuration: expected global or vdom")
	}
	global := root.Child("global")
	if global != nil && global.Kind != YAMLMapping {
		return errors.New("FortiGate YAML global section must be a mapping")
	}
	vdom := root.Child("vdom")
	if vdom != nil && vdom.Kind != YAMLSequence && vdom.Kind != YAMLMapping {
		return errors.New("FortiGate YAML vdom section must be a sequence or mapping")
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
