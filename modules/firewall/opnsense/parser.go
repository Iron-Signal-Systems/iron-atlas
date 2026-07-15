package opnsense

import (
	"context"
	"encoding/xml"
	"errors"
	"io"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/ingest"
)

type Parser struct{}

type Config struct {
	XMLName xml.Name `xml:"opnsense"`
	Version string   `xml:"version"`
	System  struct {
		Hostname string `xml:"hostname"`
		Domain   string `xml:"domain"`
	} `xml:"system"`
}

func (Parser) ID() string { return "firewall.opnsense.xml.v1" }
func (Parser) Probe(_ context.Context, reader io.Reader) (ingest.Probe, error) {
	var root struct{ XMLName xml.Name }
	if err := xml.NewDecoder(reader).Decode(&root); err != nil {
		return ingest.Probe{}, err
	}
	if root.XMLName.Local != "opnsense" {
		return ingest.Probe{}, errors.New("root element is not opnsense")
	}
	return ingest.Probe{Vendor: "deciso", Platform: "opnsense", Format: "xml"}, nil
}
func (Parser) Parse(_ context.Context, reader io.Reader) (ingest.Result, error) {
	var cfg Config
	if err := xml.NewDecoder(reader).Decode(&cfg); err != nil {
		return ingest.Result{}, err
	}
	return ingest.Result{Probe: ingest.Probe{Vendor: "deciso", Platform: "opnsense", Format: "xml", Version: cfg.Version}, Parsed: cfg}, nil
}
