package fortigate

import (
	"strings"
	"testing"
)

func TestNativeFortiGateLayoutNormalizesDirectGlobalAndDetectedVDOM(t *testing.T) {
	const input = `# config-version=FGT60F-7.2.8-FW-build1639-240313
system_global:
  hostname: native-layout
  version: 7.2.8
root:
  system_interface:
    - port1:
        ip: [192.0.2.2, 255.255.255.0]
        status: enable
  firewall_address:
    - example:
        subnet: [192.0.2.0, 255.255.255.0]
`

	got, err := ParseFortiGateYAML(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if got.Device.Hostname != "native-layout" {
		t.Fatalf("unexpected hostname %q", got.Device.Hostname)
	}
	if got.Source.FortiOSVersion != "7.2.8" {
		t.Fatalf("unexpected version %q", got.Source.FortiOSVersion)
	}
	if len(got.Domains) != 1 || got.Domains[0].Name != "root" {
		t.Fatalf("unexpected domains: %#v", got.Domains)
	}
	if len(got.Interfaces) != 1 || len(got.Objects.Addresses) != 1 {
		t.Fatalf("unexpected normalized counts: interfaces=%d addresses=%d", len(got.Interfaces), len(got.Objects.Addresses))
	}
}

func TestNativeFortiGateLayoutNormalizesDirectVDOMSections(t *testing.T) {
	const input = `system_interface:
  - port1:
      status: enable
firewall_policy:
  - 1:
      srcintf: [port1]
      dstintf: [all]
      srcaddr: [all]
      dstaddr: [all]
      service: [ALL]
      schedule: always
      action: accept
`

	got, err := ParseFortiGateYAML(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Domains) != 1 || got.Domains[0].Name != "root" {
		t.Fatalf("unexpected domains: %#v", got.Domains)
	}
	if len(got.Interfaces) != 1 || len(got.Policies) != 1 {
		t.Fatalf("unexpected normalized counts: interfaces=%d policies=%d", len(got.Interfaces), len(got.Policies))
	}
}

func TestFortiGateLayoutDiagnosticOmitsUnknownKeysAndVDOMNames(t *testing.T) {
	const privateVDOM = "customer-private-vdom"
	const privateRootKey = "customer-private-root-key"
	input := privateRootKey + `: value
` + privateVDOM + `:
  system_interface: []
  firewall_policy: []
`
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	diagnostic := DiagnoseFortiGateYAMLLayout(doc)
	if diagnostic.RootKind != "mapping" || diagnostic.RootEntries != 2 {
		t.Fatalf("unexpected diagnostic: %#v", diagnostic)
	}
	if diagnostic.DetectedVDOMContainerCount != 1 || diagnostic.UnrecognizedRootEntryCount != 1 {
		t.Fatalf("unexpected diagnostic counts: %#v", diagnostic)
	}
	joined := strings.Join(append(diagnostic.RecognizedDirectSections, diagnostic.RecognizedNestedSections...), " ")
	if strings.Contains(joined, privateVDOM) || strings.Contains(joined, privateRootKey) {
		t.Fatalf("diagnostic disclosed a private key: %#v", diagnostic)
	}
	for _, expected := range []string{"firewall_policy", "system_interface"} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("diagnostic omitted recognized schema section %q: %#v", expected, diagnostic)
		}
	}
}

func TestFortiGateLayoutRejectsUnrelatedYAML(t *testing.T) {
	_, err := ParseFortiGateYAML(strings.NewReader("application:\n  name: unrelated\n"))
	if err == nil {
		t.Fatal("expected unrelated YAML rejection")
	}
}
