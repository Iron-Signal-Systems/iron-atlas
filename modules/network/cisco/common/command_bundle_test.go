package common

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestParseCommandBundleCompatibility(t *testing.T) {
	input := "===== COMMAND: show version =====\nCisco IOS Software\n===== END COMMAND =====\n===== COMMAND: show interfaces trunk =====\nGi1/0/48 on 802.1q trunking 999\n===== END COMMAND =====\n"
	bundle, err := ParseCommandBundle(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle) != 2 || bundle["show version"] != "Cisco IOS Software" {
		t.Fatalf("unexpected bundle: %#v", bundle)
	}
}

func TestEvidenceBundlePreservesOrderAndDuplicates(t *testing.T) {
	input := bundleSection("show version", "first") + bundleSection("show version", "second")
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), DefaultParserLimits())
	if err != nil {
		t.Fatal(err)
	}
	if !bundle.Complete || bundle.Truncated || len(bundle.Commands) != 2 {
		t.Fatalf("unexpected bundle state: %#v", bundle)
	}
	if bundle.Commands[0].Sequence != 1 || bundle.Commands[0].Output != "first" ||
		bundle.Commands[1].Sequence != 2 || bundle.Commands[1].Output != "second" {
		t.Fatalf("command order or duplicates were not preserved: %#v", bundle.Commands)
	}
	for _, command := range bundle.Commands {
		if command.SchemaVersion != CommandEvidenceSchemaVersion || command.Status != CommandComplete ||
			len(command.NormalizedOutputSHA256) != 64 {
			t.Fatalf("unexpected command metadata: %#v", command)
		}
	}
	if bundle.SchemaVersion != EvidenceBundleSchemaVersion || bundle.ParserVersion != EvidenceParserVersion ||
		len(bundle.InputSHA256) != 64 || len(bundle.BundleSHA256) != 64 {
		t.Fatalf("unexpected bundle metadata: %#v", bundle)
	}
}

func TestEvidenceBundleNestedCommandRetainsPartialEvidence(t *testing.T) {
	input := "===== COMMAND: first =====\npartial\n===== COMMAND: second =====\ncomplete\n===== END COMMAND =====\n"
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), DefaultParserLimits())
	if !errors.Is(err, ErrMalformedBundle) {
		t.Fatalf("expected malformed bundle error, got %v", err)
	}
	if len(bundle.Commands) != 2 || bundle.Commands[0].Status != CommandIncomplete ||
		bundle.Commands[0].Output != "partial" || bundle.Commands[1].Status != CommandComplete {
		t.Fatalf("partial evidence was not retained: %#v", bundle.Commands)
	}
	requireDiagnostic(t, bundle, DiagnosticNestedCommand)
}

func TestEvidenceBundleEndMarkerWithoutCommand(t *testing.T) {
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(endMarker+"\n"), DefaultParserLimits())
	if !errors.Is(err, ErrMalformedBundle) {
		t.Fatalf("expected malformed bundle error, got %v", err)
	}
	requireDiagnostic(t, bundle, DiagnosticEndWithoutCommand)
}

func TestEvidenceBundleUnclosedCommand(t *testing.T) {
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader("===== COMMAND: show version =====\npartial\n"), DefaultParserLimits())
	if !errors.Is(err, ErrMalformedBundle) {
		t.Fatalf("expected malformed bundle error, got %v", err)
	}
	if len(bundle.Commands) != 1 || bundle.Commands[0].Status != CommandIncomplete || bundle.Commands[0].Output != "partial" {
		t.Fatalf("unexpected retained command: %#v", bundle.Commands)
	}
	requireDiagnostic(t, bundle, DiagnosticUnclosedCommand)
}

func TestEvidenceBundleEmptyCommandName(t *testing.T) {
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(bundleSection("", "content")), DefaultParserLimits())
	if !errors.Is(err, ErrMalformedBundle) {
		t.Fatalf("expected malformed bundle error, got %v", err)
	}
	if len(bundle.Commands) != 1 || bundle.Commands[0].Status != CommandFailed {
		t.Fatalf("unexpected empty-name command state: %#v", bundle.Commands)
	}
	requireDiagnostic(t, bundle, DiagnosticEmptyCommandName)
}

func TestEvidenceBundleInputLimit(t *testing.T) {
	limits := testLimits()
	limits.MaxInputBytes = 80
	limits.MaxCommandBytes = 64
	limits.MaxLineBytes = 64
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(bundleSection("show version", strings.Repeat("x", 64))), limits)
	if !errors.Is(err, ErrLimitExceeded) || !bundle.Truncated {
		t.Fatalf("expected truncated input-limit failure, got bundle=%#v err=%v", bundle, err)
	}
	requireDiagnostic(t, bundle, DiagnosticInputLimit)
}

func TestEvidenceBundleCommandSizeLimit(t *testing.T) {
	limits := testLimits()
	limits.MaxCommandBytes = 64
	limits.MaxLineBytes = 64
	output := strings.Repeat("1", 40) + "\n" + strings.Repeat("2", 40)
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(bundleSection("show version", output)), limits)
	if !errors.Is(err, ErrLimitExceeded) || len(bundle.Commands) != 1 || !bundle.Commands[0].Truncated {
		t.Fatalf("expected command-size failure, got bundle=%#v err=%v", bundle, err)
	}
	requireDiagnostic(t, bundle, DiagnosticCommandLimit)
}

func TestEvidenceBundleLineSizeLimit(t *testing.T) {
	limits := testLimits()
	limits.MaxLineBytes = 40
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(bundleSection("show version", strings.Repeat("x", 41))), limits)
	if !errors.Is(err, ErrLimitExceeded) || !bundle.Truncated {
		t.Fatalf("expected line-size failure, got bundle=%#v err=%v", bundle, err)
	}
	requireDiagnostic(t, bundle, DiagnosticLineLimit)
}

func TestEvidenceBundleCommandCountLimit(t *testing.T) {
	limits := testLimits()
	limits.MaxCommands = 1
	input := bundleSection("first", "one") + bundleSection("second", "two")
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), limits)
	if !errors.Is(err, ErrLimitExceeded) || len(bundle.Commands) != 1 || !bundle.Truncated {
		t.Fatalf("expected command-count failure, got bundle=%#v err=%v", bundle, err)
	}
	requireDiagnostic(t, bundle, DiagnosticCommandCountLimit)
}

func TestEvidenceBundleCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	bundle, err := ParseEvidenceBundle(ctx, strings.NewReader(bundleSection("show version", "content")), DefaultParserLimits())
	if !errors.Is(err, context.Canceled) || bundle.Complete {
		t.Fatalf("expected cancellation, got bundle=%#v err=%v", bundle, err)
	}
	requireDiagnostic(t, bundle, DiagnosticCancelled)
}

func TestEvidenceBundleDigestsAreDeterministic(t *testing.T) {
	input := bundleSection("show version", "Cisco IOS Software")
	first, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), DefaultParserLimits())
	if err != nil {
		t.Fatal(err)
	}
	second, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), DefaultParserLimits())
	if err != nil {
		t.Fatal(err)
	}
	if first.InputSHA256 != second.InputSHA256 || first.BundleSHA256 != second.BundleSHA256 ||
		first.Commands[0].NormalizedOutputSHA256 != second.Commands[0].NormalizedOutputSHA256 {
		t.Fatalf("digests are not deterministic: first=%#v second=%#v", first, second)
	}
}

func TestUploadSafeProjectionDoesNotLeakRawEvidence(t *testing.T) {
	protectedText := "protected-device-name secret-address.example"
	protectedCommand := "show protected-device-name"
	input := protectedText + "\n" + bundleSection(protectedCommand, protectedText)
	bundle, err := ParseEvidenceBundle(context.Background(), strings.NewReader(input), DefaultParserLimits())
	if err != nil {
		t.Fatal(err)
	}
	requireDiagnostic(t, bundle, DiagnosticTextOutsideCommand)
	encoded, err := json.Marshal(bundle.UploadSafe())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), protectedText) || strings.Contains(string(encoded), protectedCommand) {
		t.Fatalf("upload-safe projection leaked raw evidence: %s", encoded)
	}
}

func bundleSection(command, output string) string {
	return "===== COMMAND: " + command + " =====\n" + output + "\n===== END COMMAND =====\n"
}

func testLimits() ParserLimits {
	return ParserLimits{MaxInputBytes: 4096, MaxCommandBytes: 1024, MaxLineBytes: 256, MaxCommands: 8}
}

func requireDiagnostic(t *testing.T, bundle EvidenceBundle, code string) {
	t.Helper()
	for _, diagnostic := range bundle.Diagnostics {
		if diagnostic.Code == code {
			return
		}
	}
	t.Fatalf("diagnostic %s not found in %#v", code, bundle.Diagnostics)
}
