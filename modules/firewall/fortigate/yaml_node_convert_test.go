package fortigate

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	yaml "go.yaml.in/yaml/v4"
)

func TestMaintainedYAMLDecoderAcceptsMultilineQuotedScalar(t *testing.T) {
	file, err := os.Open("testdata/fortigate-multiline-scalar.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	doc, err := ParseYAMLDocument(file)
	if err != nil {
		t.Fatal(err)
	}
	banner := doc.Root.At("global", "system_global", "admin-banner")
	if banner == nil || banner.Kind != YAMLScalar {
		t.Fatalf("unexpected multiline node: %#v", banner)
	}
	if got, want := banner.Value, `Authorized operators only; escaped quote: "preserved"; continuation complete.`; got != want {
		t.Fatalf("unexpected multiline scalar:\n got: %q\nwant: %q", got, want)
	}
	if banner.Line != 5 || banner.Column <= 0 {
		t.Fatalf("unexpected source location: line=%d column=%d", banner.Line, banner.Column)
	}
}

func TestMaintainedYAMLDecoderPreservesOrderAndScalarText(t *testing.T) {
	input := `global:
  system_global:
    first: yes
    second: true
    third: 00123
    fourth: 7.4.5
    fifth: 192.0.2.1
    sixth: 255.255.255.0
`
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	global := doc.Root.Child("global")
	settings := global.Child("system_global")
	if got, want := strings.Join(settings.Order, ","), "first,second,third,fourth,fifth,sixth"; got != want {
		t.Fatalf("mapping order changed: got %q want %q", got, want)
	}
	wants := map[string]string{
		"first":  "yes",
		"second": "true",
		"third":  "00123",
		"fourth": "7.4.5",
		"fifth":  "192.0.2.1",
		"sixth":  "255.255.255.0",
	}
	for key, want := range wants {
		if got := settings.Child(key).Scalar(); got != want {
			t.Fatalf("scalar %s changed: got %q want %q", key, got, want)
		}
	}
}

func TestMaintainedYAMLDecoderRejectsUnsupportedFeatures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "duplicate key", input: "global:\n  value: one\n  value: two\n"},
		{name: "anchor", input: "global: &saved\n  value: one\n"},
		{name: "alias", input: "global: &saved\n  value: one\nvdom: *saved\n"},
		{name: "custom tag", input: "global:\n  value: !secret hidden\n"},
		{name: "multiple documents", input: "global: {}\n---\nvdom: {}\n"},
		{name: "flow mapping", input: "global: {value: one}\n"},
		{name: "literal block scalar", input: "global:\n  value: |\n    line\n"},
		{name: "folded block scalar", input: "global:\n  value: >\n    line\n"},
		{name: "empty document", input: "# comment only\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ParseYAMLDocument(strings.NewReader(test.input)); err == nil {
				t.Fatal("expected admission rejection")
			}
		})
	}
}

func TestFortiGateYAMLRejectsNonMappingRoot(t *testing.T) {
	_, err := ParseFortiGateYAML(strings.NewReader("- global\n"))
	if err == nil || !strings.Contains(err.Error(), "root must be a mapping") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestYAMLAdmissionLimits(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		configure func(*yamlDecodeLimits)
		limitName string
	}{
		{
			name:      "input bytes",
			input:     "a: value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxInputBytes = int64(len("a: value\n") - 1) },
			limitName: "MaxInputBytes",
		},
		{
			name:      "node count",
			input:     "a: value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxNodes = 2 },
			limitName: "MaxNodes",
		},
		{
			name:      "mapping depth",
			input:     "a:\n  b: value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxMappingDepth = 1 },
			limitName: "MaxMappingDepth",
		},
		{
			name:      "sequence depth",
			input:     "a:\n  - - value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxSequenceDepth = 1 },
			limitName: "MaxSequenceDepth",
		},
		{
			name:      "combined depth",
			input:     "a:\n  - b: value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxCombinedDepth = 2 },
			limitName: "MaxCombinedDepth",
		},
		{
			name:      "mapping entries",
			input:     "a: one\nb: two\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxMappingEntries = 1 },
			limitName: "MaxMappingEntries",
		},
		{
			name:      "sequence entries",
			input:     "a: [one, two]\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxSequenceEntries = 1 },
			limitName: "MaxSequenceEntries",
		},
		{
			name:      "scalar bytes",
			input:     "a: four\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxScalarBytes = 3 },
			limitName: "MaxScalarBytes",
		},
		{
			name:      "key bytes",
			input:     "four: value\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxKeyBytes = 3 },
			limitName: "MaxKeyBytes",
		},
		{
			name:      "total mapping entries",
			input:     "a:\n  b: one\nc: two\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxTotalMappingEntries = 2 },
			limitName: "MaxTotalMappingEntries",
		},
		{
			name:      "total sequence entries",
			input:     "a: [one, two]\nb: [three]\n",
			configure: func(limits *yamlDecodeLimits) { limits.MaxTotalSequenceEntries = 2 },
			limitName: "MaxTotalSequenceEntries",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limits := defaultYAMLDecodeLimits()
			test.configure(&limits)
			_, err := parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(test.input), limits)
			if err == nil || !strings.Contains(err.Error(), test.limitName) {
				t.Fatalf("expected %s rejection, got %v", test.limitName, err)
			}
		})
	}
}

func TestYAMLAdmissionAcceptsExactBoundaries(t *testing.T) {
	input := "a: value\n"
	limits := defaultYAMLDecodeLimits()
	limits.MaxInputBytes = int64(len(input))
	limits.MaxNodes = 3
	limits.MaxMappingDepth = 1
	limits.MaxCombinedDepth = 1
	limits.MaxMappingEntries = 1
	limits.MaxTotalMappingEntries = 1
	limits.MaxScalarBytes = len("value")
	limits.MaxKeyBytes = len("a")
	if _, err := parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(input), limits); err != nil {
		t.Fatalf("exact boundary should be accepted: %v", err)
	}

	sequenceInput := "a: [one, two]\n"
	sequenceLimits := defaultYAMLDecodeLimits()
	sequenceLimits.MaxSequenceEntries = 2
	sequenceLimits.MaxTotalSequenceEntries = 2
	if _, err := parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(sequenceInput), sequenceLimits); err != nil {
		t.Fatalf("exact sequence boundary should be accepted: %v", err)
	}
}

func TestYAMLAdmissionErrorsDoNotEchoScalarValues(t *testing.T) {
	const secret = "DO-NOT-ECHO-THIS-SECRET"
	inputs := []string{
		"global:\n  value: \"" + secret + "\n",
		"global:\n  value: " + secret + "\n  value: duplicate\n",
	}
	for _, input := range inputs {
		_, err := ParseYAMLDocument(strings.NewReader(input))
		if err == nil {
			t.Fatal("expected rejection")
		}
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("error disclosed scalar content: %q", err)
		}
	}
}

func TestYAMLStructuredLoadErrorPreservesOnlySafePosition(t *testing.T) {
	const secret = "DO-NOT-ECHO-THIS-DECODER-MESSAGE"
	source := yaml.NewLoadError(
		yaml.ScannerStage,
		secret,
		yaml.Mark{Line: 11444, Column: 17},
		errors.New(secret),
	)
	err := sanitizeYAMLDecodeError(source)
	if got, want := err.Error(), "YAML syntax decoding failed in scanner stage at line 11444, column 17"; got != want {
		t.Fatalf("unexpected structured error: got %q want %q", got, want)
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error disclosed decoder content: %q", err)
	}
}

func TestYAMLAdmissionRejectsLargeAliasInput(t *testing.T) {
	var input strings.Builder
	input.WriteString("global: &base\n  value: safe\nvdom:\n")
	for index := 0; index < 10_000; index++ {
		input.WriteString("  - *base\n")
	}
	if _, err := ParseYAMLDocument(strings.NewReader(input.String())); err == nil {
		t.Fatal("expected alias rejection")
	}
}

func TestYAMLContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := ParseYAMLDocumentContext(ctx, strings.NewReader("global: {}\n"))
	if err == nil || !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestNormalizedSnapshotLimits(t *testing.T) {
	data, err := os.ReadFile("testdata/fortigate-sanitized.yaml")
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := ParseFortiGateYAML(strings.NewReader(string(data)))
	if err != nil {
		t.Fatal(err)
	}

	limits := defaultYAMLDecodeLimits()
	limits.MaxNormalizedRecords = normalizedRecordCount(baseline)
	limits.MaxReferences = len(baseline.References.Edges) + len(baseline.References.BuiltIns) + len(baseline.References.Unresolved)
	if _, err := parseFortiGateYAMLWithLimits(context.Background(), strings.NewReader(string(data)), limits); err != nil {
		t.Fatalf("exact normalized boundaries should be accepted: %v", err)
	}

	limits.MaxNormalizedRecords--
	_, err = parseFortiGateYAMLWithLimits(context.Background(), strings.NewReader(string(data)), limits)
	if err == nil || !strings.Contains(err.Error(), "MaxNormalizedRecords") {
		t.Fatalf("expected normalized-record rejection, got %v", err)
	}

	limits = defaultYAMLDecodeLimits()
	limits.MaxReferences = len(baseline.References.Edges) + len(baseline.References.BuiltIns) + len(baseline.References.Unresolved) - 1
	_, err = parseFortiGateYAMLWithLimits(context.Background(), strings.NewReader(string(data)), limits)
	if err == nil || !strings.Contains(err.Error(), "MaxReferences") {
		t.Fatalf("expected reference rejection, got %v", err)
	}
}

func TestNormalizedFindingLimit(t *testing.T) {
	input := `vdom:
  - root:
      firewall_policy:
        - 1:
            srcintf: [missing-one]
            dstintf: [all]
            srcaddr: [all]
            dstaddr: [all]
            service: [ALL]
            schedule: always
            action: accept
        - 2:
            srcintf: [missing-two]
            dstintf: [all]
            srcaddr: [all]
            dstaddr: [all]
            service: [ALL]
            schedule: always
            action: accept
`
	limits := defaultYAMLDecodeLimits()
	limits.MaxFindings = 1
	_, err := parseFortiGateYAMLWithLimits(context.Background(), strings.NewReader(input), limits)
	if err == nil || !strings.Contains(err.Error(), "MaxFindings") {
		t.Fatalf("expected finding rejection, got %v", err)
	}
}
