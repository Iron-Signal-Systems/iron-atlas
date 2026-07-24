package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/fortigate"
	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func main() {
	var input string
	var output string
	var format string
	var compact bool
	var redact bool
	flag.StringVar(&input, "input", "", "FortiGate-generated YAML configuration file")
	flag.StringVar(&output, "output", "", "optional output file; stdout is used when omitted")
	flag.StringVar(&format, "format", "summary", "output format: summary, json, structure, or quality")
	flag.BoolVar(&compact, "compact", false, "emit compact JSON instead of indented JSON")
	flag.BoolVar(&redact, "redact", false, "omit source-derived identity, values, paths, and finding details from upload-safe output")
	flag.Parse()

	if input == "" && flag.NArg() > 0 {
		input = flag.Arg(0)
	}
	if input == "" {
		exitf("an input file is required; use -input <fortigate.yaml>")
	}
	if redact && format != "summary" && format != "structure" && format != "quality" {
		exitf("-redact requires -format summary, -format structure, or -format quality")
	}
	if format == "quality" && !redact {
		exitf("-format quality requires -redact")
	}

	file, err := os.Open(input)
	if err != nil {
		exitf("read %s: %v", inputLabel(input, redact), err)
	}
	defer file.Close()

	digest := sha256.New()
	reader := io.TeeReader(file, digest)
	if format == "structure" {
		doc, err := fortigate.ParseYAMLDocument(reader)
		if err != nil {
			exitf("decode %s for structure: %v", inputLabel(input, redact), err)
		}
		writeOutput(output, []byte(renderStructure(fortigate.DiagnoseFortiGateYAMLLayout(doc))))
		return
	}
	if format == "quality" {
		snapshotValue, layout, err := fortigate.ParseFortiGateYAMLWithLayout(reader)
		if err != nil {
			exitf("parse %s for semantic quality: %v", inputLabel(input, redact), err)
		}
		writeOutput(output, []byte(renderSemanticQuality(snapshotValue, layout)))
		return
	}

	snapshotValue, err := fortigate.ParseFortiGateYAML(reader)
	if err != nil {
		exitf("parse %s: %v", inputLabel(input, redact), err)
	}
	snapshotValue.Source.Filename = filepath.Base(input)
	snapshotValue.Source.SHA256 = hex.EncodeToString(digest.Sum(nil))

	var rendered []byte
	switch format {
	case "summary":
		rendered = []byte(renderSummary(snapshotValue, redact))
	case "json":
		if compact {
			rendered, err = json.Marshal(snapshotValue)
		} else {
			rendered, err = json.MarshalIndent(snapshotValue, "", "  ")
		}
		if err == nil {
			rendered = append(rendered, '\n')
		}
	default:
		exitf("unsupported format %q; use summary, json, structure, or quality", format)
	}
	if err != nil {
		exitf("render output: %v", err)
	}

	writeOutput(output, rendered)
}

func renderStructure(layout fortigate.FortiGateYAMLLayout) string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "Upload-safe FortiGate YAML structure")
	fmt.Fprintf(&b, "Root kind: %s\n", layout.RootKind)
	fmt.Fprintf(&b, "Root entries: %d\n", layout.RootEntries)
	fmt.Fprintf(&b, "Canonical global wrapper: %t\n", layout.CanonicalGlobal)
	fmt.Fprintf(&b, "Canonical vdom wrapper: %t\n", layout.CanonicalVDOM)
	fmt.Fprintf(&b, "Recognized direct sections: %s\n", displayList(layout.RecognizedDirectSections))
	fmt.Fprintf(&b, "Recognized nested sections: %s\n", displayList(layout.RecognizedNestedSections))
	fmt.Fprintf(&b, "Nested mappings: %d\n", layout.NestedMappingCount)
	fmt.Fprintf(&b, "Detected VDOM containers: %d\n", layout.DetectedVDOMContainerCount)
	fmt.Fprintf(&b, "Unrecognized root entries: %d\n", layout.UnrecognizedRootEntryCount)
	fmt.Fprintln(&b, "Private scalar values, unrecognized keys, and VDOM names: omitted")
	return b.String()
}

func displayList(values []string) string {
	if len(values) == 0 {
		return "none"
	}
	return strings.Join(values, ", ")
}

func writeOutput(output string, rendered []byte) {
	if output == "" {
		if _, err := os.Stdout.Write(rendered); err != nil {
			exitf("write stdout: %v", err)
		}
		return
	}
	if err := os.WriteFile(output, rendered, 0o600); err != nil {
		exitf("write %s: %v", output, err)
	}
}

func renderSummary(value *snapshot.FirewallSnapshot, redact bool) string {
	var b bytes.Buffer
	if redact {
		fmt.Fprintln(&b, "Device: [redacted]")
		fmt.Fprintln(&b, "FortiOS: [redacted]")
	} else {
		fmt.Fprintf(&b, "Device: %s\n", display(value.Device.Hostname, "unknown"))
		fmt.Fprintf(&b, "FortiOS: %s\n", display(value.Source.FortiOSVersion, "unknown"))
	}
	fmt.Fprintf(&b, "VDOMs: %d\n", len(value.Domains))
	fmt.Fprintf(&b, "Interfaces: %d\n", len(value.Interfaces))
	fmt.Fprintf(&b, "VLANs: %d\n", len(value.VLANs))
	fmt.Fprintf(&b, "Subnets: %d\n", len(value.Subnets))
	fmt.Fprintf(&b, "DHCP servers: %d\n", len(value.DHCP.Servers))
	fmt.Fprintf(&b, "Address objects: %d\n", len(value.Objects.Addresses))
	fmt.Fprintf(&b, "Address groups: %d\n", len(value.Objects.AddressGroups))
	fmt.Fprintf(&b, "Service objects: %d\n", len(value.Objects.Services))
	fmt.Fprintf(&b, "Service groups: %d\n", len(value.Objects.ServiceGroups))
	fmt.Fprintf(&b, "Firewall policies: %d\n", len(value.Policies))
	fmt.Fprintf(&b, "Static routes: %d\n", len(value.Routing.Routes))
	fmt.Fprintf(&b, "Policy routes: %d\n", len(value.Routing.PolicyRoutes))
	fmt.Fprintf(&b, "SD-WAN zones: %d\n", len(value.SDWAN.Zones))
	fmt.Fprintf(&b, "SD-WAN members: %d\n", len(value.SDWAN.Members))
	fmt.Fprintf(&b, "SD-WAN health checks: %d\n", len(value.SDWAN.HealthChecks))
	fmt.Fprintf(&b, "SD-WAN rules: %d\n", len(value.SDWAN.Rules))
	fmt.Fprintf(&b, "Virtual IPs: %d\n", len(value.NAT.VirtualIPs))
	fmt.Fprintf(&b, "IP pools: %d\n", len(value.NAT.IPPools))
	fmt.Fprintf(&b, "VPNs: %d\n", len(value.VPNs))
	fmt.Fprintf(&b, "Traffic shapers: %d\n", len(value.QoS.TrafficShapers))
	fmt.Fprintf(&b, "Reference edges: %d\n", len(value.References.Edges))
	fmt.Fprintf(&b, "Built-in references: %d\n", len(value.References.BuiltIns))
	fmt.Fprintf(&b, "Unresolved references: %d\n", len(value.References.Unresolved))
	fmt.Fprintf(&b, "Findings: %d\n", len(value.Findings))

	if len(value.Findings) > 0 && !redact {
		fmt.Fprintln(&b, "\nFindings:")
		findings := append([]snapshot.Finding(nil), value.Findings...)
		sort.SliceStable(findings, func(i, j int) bool {
			if findings[i].Severity != findings[j].Severity {
				return findings[i].Severity < findings[j].Severity
			}
			return findings[i].ID < findings[j].ID
		})
		for _, finding := range findings {
			fmt.Fprintf(&b, "- [%s] %s: %s\n", finding.Severity, finding.Title, finding.Detail)
		}
	}
	return b.String()
}

func inputLabel(input string, redact bool) string {
	if redact {
		return "input"
	}
	return input
}

func display(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fortigate-inspect: "+format+"\n", args...)
	os.Exit(1)
}
