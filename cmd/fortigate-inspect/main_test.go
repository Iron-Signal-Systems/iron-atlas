package main

import (
	"os"
	"strings"
	"testing"

	"github.com/Iron-Signal-Systems/atlas/modules/firewall/fortigate"
	"github.com/Iron-Signal-Systems/atlas/modules/firewall/snapshot"
)

func TestRedactedSummaryOmitsIdentityVersionAndFindingDetails(t *testing.T) {
	const secret = "DO-NOT-ECHO-PRIVATE-VALUE"
	value := &snapshot.FirewallSnapshot{
		Source: snapshot.SnapshotSource{FortiOSVersion: secret},
		Device: snapshot.Device{Hostname: secret},
		Findings: []snapshot.Finding{{
			ID:       "test",
			Severity: "info",
			Title:    secret,
			Detail:   secret,
		}},
	}
	rendered := renderSummary(value, true)
	if strings.Contains(rendered, secret) {
		t.Fatalf("redacted summary disclosed private content: %q", rendered)
	}
	for _, expected := range []string{
		"Device: [redacted]",
		"FortiOS: [redacted]",
		"Findings: 1",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("redacted summary missing %q: %q", expected, rendered)
		}
	}
}

func TestInputLabelRedaction(t *testing.T) {
	const privatePath = "/private/customer/device.yaml"
	if got := inputLabel(privatePath, true); got != "input" {
		t.Fatalf("unexpected redacted input label: %q", got)
	}
	if got := inputLabel(privatePath, false); got != privatePath {
		t.Fatalf("unexpected full input label: %q", got)
	}
}

func TestStructureOutputContainsOnlySafeSchemaDiagnostics(t *testing.T) {
	const privateValue = "DO-NOT-ECHO-PRIVATE-VDOM"
	rendered := renderStructure(fortigate.FortiGateYAMLLayout{
		RootKind:                   "mapping",
		RootEntries:                3,
		RecognizedDirectSections:   []string{"system_global"},
		RecognizedNestedSections:   []string{"firewall_policy", "system_interface"},
		NestedMappingCount:         1,
		DetectedVDOMContainerCount: 1,
		UnrecognizedRootEntryCount: 1,
	})
	if strings.Contains(rendered, privateValue) {
		t.Fatalf("structure output disclosed private content: %q", rendered)
	}
	for _, expected := range []string{
		"Root kind: mapping",
		"Recognized direct sections: system_global",
		"Recognized nested sections: firewall_policy, system_interface",
		"Detected VDOM containers: 1",
		"Private scalar values, unrecognized keys, and VDOM names: omitted",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("structure output missing %q: %q", expected, rendered)
		}
	}
}

func TestSemanticQualityReportOmitsEverySourceDerivedField(t *testing.T) {
	const secret = "DO-NOT-ECHO-PRIVATE-SEMANTIC-VALUE"
	value := &snapshot.FirewallSnapshot{
		Source: snapshot.SnapshotSource{
			FortiOSVersion: secret,
			Filename:       secret,
			SHA256:         secret,
		},
		Device: snapshot.Device{
			Hostname:     secret,
			Alias:        secret,
			SerialNumber: secret,
		},
		Domains: []snapshot.RoutingDomain{{
			ID:    snapshot.ObjectID(secret),
			Name:  secret,
			Scope: snapshot.ObjectScope{VDOM: secret},
		}},
		References: snapshot.ReferenceGraph{
			Edges: []snapshot.ReferenceEdge{
				{From: snapshot.ObjectID(secret), To: snapshot.ObjectID(secret), Kind: snapshot.ObjectKindAddress, Role: "policy-source-address", Source: snapshot.SourceLocation{YAMLPath: secret}},
				{Kind: snapshot.ObjectKind(secret), Role: secret},
			},
			BuiltIns: []snapshot.BuiltInReference{
				{From: snapshot.ObjectID(secret), Kind: snapshot.ObjectKindSchedule, Role: "policy-schedule", Source: snapshot.SourceLocation{YAMLPath: secret}},
				{Kind: snapshot.ObjectKind(secret), Role: secret},
			},
			Unresolved: []snapshot.UnresolvedReference{
				{From: snapshot.ObjectID(secret), Kind: snapshot.ObjectKindInterfaceOrZone, VendorName: secret, Role: "policy-source-interface", Scope: snapshot.ObjectScope{VDOM: secret}, Resolution: snapshot.ReferenceUnresolved, Source: snapshot.SourceLocation{YAMLPath: secret}},
				{From: snapshot.ObjectID(secret), Kind: snapshot.ObjectKindService, VendorName: secret, Role: "policy-service", Resolution: snapshot.ReferenceAmbiguous},
				{Kind: snapshot.ObjectKind(secret), VendorName: secret, Role: secret, Resolution: snapshot.ReferenceUnresolved},
			},
		},
		Findings: []snapshot.Finding{
			{ID: secret, Severity: "warning", Category: "reference", Title: "Unresolved object reference", Detail: secret, ObjectID: snapshot.ObjectID(secret), Source: snapshot.SourceLocation{YAMLPath: secret}},
			{ID: secret, Severity: secret, Category: secret, Title: secret, Detail: secret},
		},
		Extensions: []snapshot.Extension{{Scope: snapshot.ObjectScope{VDOM: secret}, Path: secret, Name: secret, Reason: secret}},
	}
	layout := fortigate.FortiGateYAMLLayout{
		RootEntries:                12,
		RecognizedDirectSections:   []string{"firewall_policy", "system_interface"},
		UnrecognizedRootEntryCount: 10,
	}

	rendered := renderSemanticQuality(value, layout)
	if strings.Contains(rendered, secret) {
		t.Fatalf("semantic-quality report disclosed private content: %q", rendered)
	}
	for _, expected := range []string{
		"Privacy mode: aggregate-only",
		"- address: 0",
		"- resolved: 2",
		"- built-in: 2",
		"- unresolved: 2",
		"- ambiguous: 1",
		"- policy-source-address | address: 1",
		"- schedule: 1",
		"- policy-source-interface | interface-or-zone: 1",
		"- policy-service | service: 1",
		"- warning: 1",
		"- reference: 1",
		"- Unresolved object reference: 1",
		"- unclassified-safe-label-buckets: 8",
		"Private vendor names, object identifiers, source paths, scalar values, finding details, VDOM names, input paths, and normalized JSON: omitted",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("semantic-quality report missing %q: %q", expected, rendered)
		}
	}
}

func TestSemanticQualityReportUsesStableSortedAggregateRows(t *testing.T) {
	value := &snapshot.FirewallSnapshot{
		References: snapshot.ReferenceGraph{
			Edges: []snapshot.ReferenceEdge{
				{Kind: snapshot.ObjectKindService, Role: "policy-service"},
				{Kind: snapshot.ObjectKindAddress, Role: "policy-source-address"},
			},
		},
	}
	rendered := renderSemanticQuality(value, fortigate.FortiGateYAMLLayout{})
	addressIndex := strings.Index(rendered, "policy-source-address | address")
	serviceIndex := strings.Index(rendered, "policy-service | service")
	if addressIndex < 0 || serviceIndex < 0 || serviceIndex > addressIndex {
		t.Fatalf("expected role/kind rows in stable lexical order: %q", rendered)
	}
}

func TestSanitizedFixtureSemanticQualityUsesOnlyClassifiedLabels(t *testing.T) {
	file, err := os.Open("../../modules/firewall/fortigate/testdata/fortigate-sanitized.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	value, layout, err := fortigate.ParseFortiGateYAMLWithLayout(file)
	if err != nil {
		t.Fatal(err)
	}
	report := buildSemanticQualityReport(value, layout)
	if report.SafeLabelFallbacks != 0 {
		t.Fatalf("sanitized fixture produced unclassified safe labels: %#v", report)
	}
	if got, want := sumRoleKindCounts(report.Resolved), 38; got != want {
		t.Fatalf("unexpected resolved reference total: got %d want %d", got, want)
	}
	if got, want := sumStringCounts(report.BuiltIns), 6; got != want {
		t.Fatalf("unexpected built-in reference total: got %d want %d", got, want)
	}
	if got := sumRoleKindCounts(report.Unresolved) + sumRoleKindCounts(report.Ambiguous); got != 0 {
		t.Fatalf("expected no unresolved fixture references, got %d", got)
	}
}
