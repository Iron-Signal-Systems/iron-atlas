package pfsense

import (
	"context"
	"strings"
	"testing"
)

func TestParseXML(t *testing.T) {
	result, err := (Parser{}).Parse(context.Background(), strings.NewReader(`<pfsense><version>24.03</version><system><hostname>fw2</hostname><domain>example.test</domain></system></pfsense>`))
	if err != nil {
		t.Fatal(err)
	}
	cfg := result.Parsed.(Config)
	if cfg.System.Hostname != "fw2" || result.Probe.Version != "24.03" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
