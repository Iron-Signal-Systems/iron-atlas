package ingest

import (
	"context"
	"errors"
	"io"
	"sort"
)

type Probe struct {
	Vendor   string
	Platform string
	Format   string
	Version  string
}

type Result struct {
	Probe        Probe
	Parsed       any
	Warnings     []string
	Unrecognized []string
}

type Parser interface {
	ID() string
	Probe(context.Context, io.Reader) (Probe, error)
	Parse(context.Context, io.Reader) (Result, error)
}

type Registry struct{ parsers map[string]Parser }

func NewRegistry() *Registry { return &Registry{parsers: make(map[string]Parser)} }

func (r *Registry) Register(parser Parser) error {
	if parser == nil || parser.ID() == "" {
		return errors.New("parser with non-empty ID is required")
	}
	if _, exists := r.parsers[parser.ID()]; exists {
		return errors.New("parser ID is already registered")
	}
	r.parsers[parser.ID()] = parser
	return nil
}

func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.parsers))
	for id := range r.parsers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
