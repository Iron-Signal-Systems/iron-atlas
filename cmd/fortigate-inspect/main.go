package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/fortigate"
	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/snapshot"
)

func main() {
	var input string
	var output string
	var format string
	var compact bool
	flag.StringVar(&input, "input", "", "FortiGate-generated YAML configuration file")
	flag.StringVar(&output, "output", "", "optional output file; stdout is used when omitted")
	flag.StringVar(&format, "format", "summary", "output format: summary or json")
	flag.BoolVar(&compact, "compact", false, "emit compact JSON instead of indented JSON")
	flag.Parse()

	if input == "" && flag.NArg() > 0 {
		input = flag.Arg(0)
	}
	if input == "" {
		exitf("an input file is required; use -input <fortigate.yaml>")
	}

	data, err := os.ReadFile(input)
	if err != nil {
		exitf("read %s: %v", input, err)
	}
	snapshotValue, err := fortigate.ParseFortiGateYAML(bytes.NewReader(data))
	if err != nil {
		exitf("parse %s: %v", input, err)
	}
	digest := sha256.Sum256(data)
	snapshotValue.Source.Filename = filepath.Base(input)
	snapshotValue.Source.SHA256 = hex.EncodeToString(digest[:])

	var rendered []byte
	switch format {
	case "summary":
		rendered = []byte(renderSummary(snapshotValue))
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
		exitf("unsupported format %q; use summary or json", format)
	}
	if err != nil {
		exitf("render output: %v", err)
	}

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

func renderSummary(value *snapshot.FirewallSnapshot) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Device: %s\n", display(value.Device.Hostname, "unknown"))
	fmt.Fprintf(&b, "FortiOS: %s\n", display(value.Source.FortiOSVersion, "unknown"))
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
	fmt.Fprintf(&b, "Unresolved references: %d\n", len(value.References.Unresolved))
	fmt.Fprintf(&b, "Findings: %d\n", len(value.Findings))

	if len(value.Findings) > 0 {
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
