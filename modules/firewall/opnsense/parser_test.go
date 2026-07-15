package opnsense

import (
	"context"
	"strings"
	"testing"
)

func TestParseXML(t *testing.T) {
	result, err := (Parser{}).Parse(context.Background(), strings.NewReader(`<opnsense><version>24.7</version><system><hostname>fw1</hostname><domain>example.test</domain></system></opnsense>`))
	if err != nil {
		t.Fatal(err)
	}
	cfg := result.Parsed.(Config)
	if cfg.System.Hostname != "fw1" || result.Probe.Version != "24.7" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
