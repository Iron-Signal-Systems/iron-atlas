package pfsense

import (
	"context"
	"encoding/xml"
	"errors"
	"io"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/ingest"
)

type Parser struct{}

type Config struct {
	XMLName xml.Name `xml:"pfsense"`
	Version string   `xml:"version"`
	System  struct {
		Hostname string `xml:"hostname"`
		Domain   string `xml:"domain"`
	} `xml:"system"`
}

func (Parser) ID() string { return "firewall.pfsense.xml.v1" }
func (Parser) Probe(_ context.Context, reader io.Reader) (ingest.Probe, error) {
	var root struct{ XMLName xml.Name }
	if err := xml.NewDecoder(reader).Decode(&root); err != nil {
		return ingest.Probe{}, err
	}
	if root.XMLName.Local != "pfsense" {
		return ingest.Probe{}, errors.New("root element is not pfsense")
	}
	return ingest.Probe{Vendor: "netgate", Platform: "pfsense", Format: "xml"}, nil
}
func (Parser) Parse(_ context.Context, reader io.Reader) (ingest.Result, error) {
	var cfg Config
	if err := xml.NewDecoder(reader).Decode(&cfg); err != nil {
		return ingest.Result{}, err
	}
	return ingest.Result{Probe: ingest.Probe{Vendor: "netgate", Platform: "pfsense", Format: "xml", Version: cfg.Version}, Parsed: cfg}, nil
}
