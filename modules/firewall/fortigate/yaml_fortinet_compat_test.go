package fortigate

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFortinetCompatibilityRepairsAdjacentQuotedValues(t *testing.T) {
	file, err := os.Open("testdata/fortigate-adjacent-quoted-values.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	doc, err := ParseYAMLDocument(file)
	if err != nil {
		t.Fatal(err)
	}
	settings := doc.Root.At("global", "system_global")
	value := settings.Child("multi-value")
	if value == nil || value.Kind != YAMLSequence || value.Line != 3 {
		t.Fatalf("unexpected repaired node: %#v", value)
	}
	got := valuesField(settings, "multi-value")
	want := []string{"first value", "second", "third-value"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected repaired values: got %#v want %#v", got, want)
	}
}

func TestFortinetCompatibilityRepairsMultipleSeparatingSpaces(t *testing.T) {
	input := "global:\n  values: \"one\"   two\n"
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	got := valuesField(doc.Root.Child("global"), "values")
	want := []string{"one", "two"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected repaired values: got %#v want %#v", got, want)
	}
}

func TestFortinetCompatibilityRepairsQuotedOnlyAdjacentValues(t *testing.T) {
	input := "global:\n  values: \"one value\" \"two\"\n"
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	got := valuesField(doc.Root.Child("global"), "values")
	want := []string{"one value", "two"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected repaired values: got %#v want %#v", got, want)
	}
}

func TestFortinetCompatibilityPreservesOrdinaryValidScalars(t *testing.T) {
	input := "global:\n  description: ordinary plain words\n  quoted: \"one value\" # retained comment\n"
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	global := doc.Root.Child("global")
	if got := scalarField(global, "description"); got != "ordinary plain words" {
		t.Fatalf("plain scalar changed: %q", got)
	}
	if got := scalarField(global, "quoted"); got != "one value" {
		t.Fatalf("quoted scalar changed: %q", got)
	}
}

func TestFortinetCompatibilityQuotesUnsafeLiteralMappingKeys(t *testing.T) {
	names := []string{
		"*literal-name",
		"*literal?",
		"*literal$name",
		"&literal-name",
		"!literal-name",
		"%literal-name",
		"@literal-name",
	}
	for _, name := range names {
		input := fmt.Sprintf(
			"global:\n  objects:\n    - %s:\n        name: %q\n",
			name,
			name,
		)
		doc, err := ParseYAMLDocument(strings.NewReader(input))
		if err != nil {
			t.Fatalf("repair %q: %v", name, err)
		}
		objects := doc.Root.At("global", "objects")
		if objects == nil || objects.Kind != YAMLSequence || len(objects.Seq) != 1 {
			t.Fatalf("unexpected object sequence for %q: %#v", name, objects)
		}
		object := objects.Seq[0].Child(name)
		if object == nil || object.Kind != YAMLMapping || scalarField(object, "name") != name {
			t.Fatalf("unsafe literal key was not preserved for %q: %#v", name, object)
		}
	}
}

func TestFortinetCompatibilityDoesNotAdmitAliasValues(t *testing.T) {
	input := "global:\n  objects:\n    - value: *literal-name\n"
	if _, err := ParseYAMLDocument(strings.NewReader(input)); err == nil {
		t.Fatal("expected alias value to remain rejected")
	}
}

func TestFortinetCompatibilityDoesNotRewriteUnsafeKeysWithInlineValues(t *testing.T) {
	input := "global:\n  objects:\n    - *literal-name: inline-value\n"
	if _, err := ParseYAMLDocument(strings.NewReader(input)); err == nil {
		t.Fatal("expected unsafe key with inline value to remain rejected")
	}
}

func TestFortinetCompatibilityDoesNotRewriteQuotedScalarContinuation(t *testing.T) {
	input := "global:\n  system_global:\n    note: 'first\n      fake-key: \"not\" \"a repair\"\n      last'\n"
	doc, err := ParseYAMLDocument(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	value := doc.Root.At("global", "system_global", "note")
	if value == nil || value.Kind != YAMLScalar {
		t.Fatalf("quoted continuation was rewritten: %#v", value)
	}
}

func TestFortinetCompatibilityLimits(t *testing.T) {
	input := "global:\n  one: \"a\" \"b\"\n  two: \"c\" \"d\" \"e\"\n"

	rewriteLimits := defaultYAMLDecodeLimits()
	rewriteLimits.MaxCompatibilityRewrites = 1
	_, err := parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(input), rewriteLimits)
	if err == nil || !strings.Contains(err.Error(), "MaxCompatibilityRewrites") {
		t.Fatalf("expected rewrite-limit rejection, got %v", err)
	}

	fragmentLimits := defaultYAMLDecodeLimits()
	fragmentLimits.MaxCompatibilityFragments = 4
	_, err = parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(input), fragmentLimits)
	if err == nil || !strings.Contains(err.Error(), "MaxCompatibilityFragments") {
		t.Fatalf("expected fragment-limit rejection, got %v", err)
	}

	byteLimits := defaultYAMLDecodeLimits()
	byteLimits.MaxInputBytes = int64(len(input))
	_, err = parseYAMLDocumentWithLimits(context.Background(), strings.NewReader(input), byteLimits)
	if err == nil || !strings.Contains(err.Error(), "MaxInputBytes") {
		t.Fatalf("expected post-rewrite byte-limit rejection, got %v", err)
	}
}

func TestFortinetCompatibilityRejectsBroaderInvalidForms(t *testing.T) {
	inputs := []string{
		"global:\n  value: 'one' 'two'\n",
		"global:\n  value: \"one\" [\"two\"]\n",
		"global:\n  value: \"one\" two,three\n",
		"global:\n  value: \"one\" *alias\n",
	}
	for _, input := range inputs {
		if _, err := ParseYAMLDocument(strings.NewReader(input)); err == nil {
			t.Fatalf("expected unsupported invalid YAML rejection: %q", input)
		}
	}
}
