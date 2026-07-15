package fortigate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/ingest"
)

type Node struct {
	Kind     string
	Name     string
	Values   map[string][]string
	Children []*Node
}

type Document struct {
	Root     *Node
	Comments []string
}

type Parser struct{}

func (Parser) ID() string { return "firewall.fortigate.native.v1" }

func (Parser) Probe(_ context.Context, reader io.Reader) (ingest.Probe, error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "config ") || strings.HasPrefix(line, "#config-version=") {
			return ingest.Probe{Vendor: "fortinet", Platform: "fortigate", Format: "fortios-native"}, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return ingest.Probe{}, err
	}
	return ingest.Probe{}, errors.New("input does not appear to be a FortiOS native configuration")
}

func (Parser) Parse(_ context.Context, reader io.Reader) (ingest.Result, error) {
	doc, err := ParseNative(reader)
	if err != nil {
		return ingest.Result{}, err
	}
	return ingest.Result{Probe: ingest.Probe{Vendor: "fortinet", Platform: "fortigate", Format: "fortios-native"}, Parsed: doc}, nil
}

func ParseNative(reader io.Reader) (*Document, error) {
	root := &Node{Kind: "root", Name: "root", Values: make(map[string][]string)}
	stack := []*Node{root}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	lineNo := 0
	doc := &Document{Root: root}

	for scanner.Scan() {
		lineNo++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}
		if strings.HasPrefix(raw, "#") {
			doc.Comments = append(doc.Comments, raw)
			continue
		}
		fields := splitFields(raw)
		if len(fields) == 0 {
			continue
		}
		current := stack[len(stack)-1]
		switch fields[0] {
		case "config":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: config requires a name", lineNo)
			}
			node := &Node{Kind: "config", Name: strings.Join(fields[1:], " "), Values: make(map[string][]string)}
			current.Children = append(current.Children, node)
			stack = append(stack, node)
		case "edit":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: edit requires a name", lineNo)
			}
			node := &Node{Kind: "edit", Name: strings.Join(fields[1:], " "), Values: make(map[string][]string)}
			current.Children = append(current.Children, node)
			stack = append(stack, node)
		case "set", "append", "unset":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: %s requires a key", lineNo, fields[0])
			}
			key := fields[0] + " " + fields[1]
			current.Values[key] = append(current.Values[key], fields[2:]...)
		case "next":
			if len(stack) < 2 || current.Kind != "edit" {
				return nil, fmt.Errorf("line %d: unexpected next", lineNo)
			}
			stack = stack[:len(stack)-1]
		case "end":
			if len(stack) < 2 {
				return nil, fmt.Errorf("line %d: unexpected end", lineNo)
			}
			if current.Kind == "edit" {
				return nil, fmt.Errorf("line %d: edit must close with next before end", lineNo)
			}
			stack = stack[:len(stack)-1]
		default:
			current.Values["raw"] = append(current.Values["raw"], raw)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("configuration ended with %d unclosed blocks", len(stack)-1)
	}
	return doc, nil
}

func splitFields(line string) []string {
	var fields []string
	var b strings.Builder
	quoted := false
	escaped := false
	flush := func() {
		if b.Len() > 0 {
			fields = append(fields, b.String())
			b.Reset()
		}
	}
	for _, r := range line {
		if escaped {
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && quoted {
			escaped = true
			continue
		}
		if r == '"' {
			quoted = !quoted
			continue
		}
		if (r == ' ' || r == '\t') && !quoted {
			flush()
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return fields
}
