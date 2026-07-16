package fortigate

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type YAMLKind uint8

const (
	YAMLScalar YAMLKind = iota + 1
	YAMLMapping
	YAMLSequence
)

type YAMLNode struct {
	Kind   YAMLKind
	Value  string
	Map    map[string]*YAMLNode
	Order  []string
	Seq    []*YAMLNode
	Line   int
	Column int
}

type YAMLDocument struct {
	Root     *YAMLNode
	Comments []string
}

type yamlLine struct {
	indent int
	text   string
	line   int
}

func ParseYAMLDocument(reader io.Reader) (*YAMLDocument, error) {
	if reader == nil {
		return nil, errors.New("reader is required")
	}

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)

	var lines []yamlLine
	doc := &YAMLDocument{}
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if lineNo == 1 {
			raw = strings.TrimPrefix(raw, "\ufeff")
		}
		if strings.ContainsRune(raw, '\t') {
			return nil, fmt.Errorf("line %d: tabs are not permitted in FortiGate YAML indentation", lineNo)
		}

		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || trimmed == "---" || trimmed == "..." {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			doc.Comments = append(doc.Comments, trimmed)
			continue
		}

		content, err := stripYAMLComment(raw)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if strings.TrimSpace(content) == "" {
			continue
		}
		indent := len(content) - len(strings.TrimLeft(content, " "))
		lines = append(lines, yamlLine{indent: indent, text: strings.TrimSpace(content), line: lineNo})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("empty YAML document")
	}
	if lines[0].indent != 0 {
		return nil, fmt.Errorf("line %d: top-level YAML content must start at column 1", lines[0].line)
	}

	root, next, err := parseYAMLBlock(lines, 0, 0)
	if err != nil {
		return nil, err
	}
	if next != len(lines) {
		return nil, fmt.Errorf("line %d: unexpected trailing YAML content", lines[next].line)
	}
	doc.Root = root
	return doc, nil
}

func parseYAMLBlock(lines []yamlLine, start, indent int) (*YAMLNode, int, error) {
	if start >= len(lines) {
		return nil, start, errors.New("unexpected end of YAML document")
	}
	if lines[start].indent != indent {
		return nil, start, fmt.Errorf("line %d: expected indentation %d, got %d", lines[start].line, indent, lines[start].indent)
	}
	if strings.HasPrefix(lines[start].text, "-") {
		return parseYAMLSequence(lines, start, indent)
	}
	return parseYAMLMapping(lines, start, indent)
}

func parseYAMLMapping(lines []yamlLine, start, indent int) (*YAMLNode, int, error) {
	node := &YAMLNode{Kind: YAMLMapping, Map: make(map[string]*YAMLNode), Line: lines[start].line, Column: indent + 1}
	i := start
	for i < len(lines) {
		line := lines[i]
		if line.indent < indent {
			break
		}
		if line.indent > indent {
			return nil, i, fmt.Errorf("line %d: unexpected indentation %d; expected %d", line.line, line.indent, indent)
		}
		if strings.HasPrefix(line.text, "-") {
			break
		}

		key, rawValue, err := splitYAMLKeyValue(line.text)
		if err != nil {
			return nil, i, fmt.Errorf("line %d: %w", line.line, err)
		}
		if _, exists := node.Map[key]; exists {
			return nil, i, fmt.Errorf("line %d: duplicate mapping key %q", line.line, key)
		}

		var child *YAMLNode
		if rawValue != "" {
			child, err = parseYAMLInlineValue(rawValue, line.line, indent+len(key)+2)
			if err != nil {
				return nil, i, err
			}
			i++
		} else if i+1 < len(lines) && lines[i+1].indent > indent {
			child, i, err = parseYAMLBlock(lines, i+1, lines[i+1].indent)
			if err != nil {
				return nil, i, err
			}
		} else {
			child = &YAMLNode{Kind: YAMLMapping, Map: make(map[string]*YAMLNode), Line: line.line, Column: indent + len(key) + 2}
			i++
		}
		node.Map[key] = child
		node.Order = append(node.Order, key)
	}
	return node, i, nil
}

func parseYAMLSequence(lines []yamlLine, start, indent int) (*YAMLNode, int, error) {
	node := &YAMLNode{Kind: YAMLSequence, Line: lines[start].line, Column: indent + 1}
	i := start
	for i < len(lines) {
		line := lines[i]
		if line.indent < indent {
			break
		}
		if line.indent != indent || !strings.HasPrefix(line.text, "-") {
			break
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line.text, "-"))
		if rest == "" {
			if i+1 >= len(lines) || lines[i+1].indent <= indent {
				return nil, i, fmt.Errorf("line %d: sequence item requires a value", line.line)
			}
			child, next, err := parseYAMLBlock(lines, i+1, lines[i+1].indent)
			if err != nil {
				return nil, i, err
			}
			node.Seq = append(node.Seq, child)
			i = next
			continue
		}

		if key, rawValue, ok := trySplitYAMLKeyValue(rest); ok {
			item := &YAMLNode{Kind: YAMLMapping, Map: make(map[string]*YAMLNode), Line: line.line, Column: indent + 3}
			var child *YAMLNode
			var err error
			if rawValue != "" {
				child, err = parseYAMLInlineValue(rawValue, line.line, indent+3+len(key)+2)
				if err != nil {
					return nil, i, err
				}
				i++
			} else if i+1 < len(lines) && lines[i+1].indent > indent {
				child, i, err = parseYAMLBlock(lines, i+1, lines[i+1].indent)
				if err != nil {
					return nil, i, err
				}
			} else {
				child = &YAMLNode{Kind: YAMLMapping, Map: make(map[string]*YAMLNode), Line: line.line, Column: indent + 3 + len(key) + 2}
				i++
			}
			item.Map[key] = child
			item.Order = append(item.Order, key)

			// A sequence item may continue with sibling mapping keys indented beneath the dash.
			if i < len(lines) && lines[i].indent > indent && !strings.HasPrefix(lines[i].text, "-") {
				siblings, next, err := parseYAMLMapping(lines, i, lines[i].indent)
				if err != nil {
					return nil, i, err
				}
				for _, siblingKey := range siblings.Order {
					if _, duplicate := item.Map[siblingKey]; duplicate {
						return nil, i, fmt.Errorf("line %d: duplicate sequence-item key %q", siblings.Map[siblingKey].Line, siblingKey)
					}
					item.Map[siblingKey] = siblings.Map[siblingKey]
					item.Order = append(item.Order, siblingKey)
				}
				i = next
			}
			node.Seq = append(node.Seq, item)
			continue
		}

		child, err := parseYAMLInlineValue(rest, line.line, indent+3)
		if err != nil {
			return nil, i, err
		}
		node.Seq = append(node.Seq, child)
		i++
	}
	return node, i, nil
}

func parseYAMLInlineValue(raw string, line, column int) (*YAMLNode, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return &YAMLNode{Kind: YAMLScalar, Line: line, Column: column}, nil
	}
	if strings.HasPrefix(raw, "|") || strings.HasPrefix(raw, ">") {
		return nil, fmt.Errorf("line %d: block scalars are not supported by the bounded FortiGate YAML parser", line)
	}
	if strings.Contains(raw, "&") || strings.HasPrefix(raw, "*") || strings.HasPrefix(raw, "!") {
		return nil, fmt.Errorf("line %d: YAML anchors, aliases, and tags are not supported", line)
	}
	if strings.HasPrefix(raw, "[") {
		if !strings.HasSuffix(raw, "]") {
			return nil, fmt.Errorf("line %d: unterminated flow sequence", line)
		}
		parts, err := splitYAMLFlow(strings.TrimSpace(raw[1 : len(raw)-1]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		node := &YAMLNode{Kind: YAMLSequence, Line: line, Column: column}
		for _, part := range parts {
			value, err := decodeYAMLScalar(part)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", line, err)
			}
			node.Seq = append(node.Seq, &YAMLNode{Kind: YAMLScalar, Value: value, Line: line, Column: column})
		}
		return node, nil
	}
	if strings.HasPrefix(raw, "{") {
		return nil, fmt.Errorf("line %d: flow mappings are not supported", line)
	}
	value, err := decodeYAMLScalar(raw)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", line, err)
	}
	return &YAMLNode{Kind: YAMLScalar, Value: value, Line: line, Column: column}, nil
}

func splitYAMLKeyValue(text string) (string, string, error) {
	key, value, ok := trySplitYAMLKeyValue(text)
	if !ok {
		return "", "", errors.New("mapping entry must contain ':'")
	}
	if key == "" {
		return "", "", errors.New("mapping key is empty")
	}
	return key, value, nil
}

func trySplitYAMLKeyValue(text string) (string, string, bool) {
	quoted := rune(0)
	escaped := false
	for i, r := range text {
		if escaped {
			escaped = false
			continue
		}
		if quoted == '"' && r == '\\' {
			escaped = true
			continue
		}
		if r == '\'' || r == '"' {
			if quoted == 0 {
				quoted = r
			} else if quoted == r {
				quoted = 0
			}
			continue
		}
		if r == ':' && quoted == 0 {
			runesBefore := []rune(text[:i])
			keyRaw := strings.TrimSpace(string(runesBefore))
			key, err := decodeYAMLScalar(keyRaw)
			if err != nil {
				return "", "", false
			}
			return key, strings.TrimSpace(text[i+1:]), true
		}
	}
	return "", "", false
}

func decodeYAMLScalar(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" || raw == "~" {
		return "", nil
	}
	if strings.HasPrefix(raw, "\"") {
		if !strings.HasSuffix(raw, "\"") || len(raw) < 2 {
			return "", errors.New("unterminated double-quoted scalar")
		}
		value, err := strconv.Unquote(raw)
		if err != nil {
			return "", fmt.Errorf("invalid double-quoted scalar: %w", err)
		}
		return value, nil
	}
	if strings.HasPrefix(raw, "'") {
		if !strings.HasSuffix(raw, "'") || len(raw) < 2 {
			return "", errors.New("unterminated single-quoted scalar")
		}
		return strings.ReplaceAll(raw[1:len(raw)-1], "''", "'"), nil
	}
	return raw, nil
}

func stripYAMLComment(raw string) (string, error) {
	quoted := rune(0)
	escaped := false
	for i, r := range raw {
		if escaped {
			escaped = false
			continue
		}
		if quoted == '"' && r == '\\' {
			escaped = true
			continue
		}
		if r == '\'' || r == '"' {
			if quoted == 0 {
				quoted = r
			} else if quoted == r {
				quoted = 0
			}
			continue
		}
		if r == '#' && quoted == 0 {
			if i == 0 || raw[i-1] == ' ' {
				return strings.TrimRight(raw[:i], " "), nil
			}
		}
	}
	if quoted != 0 {
		return "", errors.New("unterminated quoted scalar")
	}
	return raw, nil
}

func splitYAMLFlow(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var parts []string
	var current strings.Builder
	quoted := rune(0)
	escaped := false
	for _, r := range raw {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if quoted == '"' && r == '\\' {
			current.WriteRune(r)
			escaped = true
			continue
		}
		if r == '\'' || r == '"' {
			current.WriteRune(r)
			if quoted == 0 {
				quoted = r
			} else if quoted == r {
				quoted = 0
			}
			continue
		}
		if r == ',' && quoted == 0 {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			continue
		}
		current.WriteRune(r)
	}
	if quoted != 0 {
		return nil, errors.New("unterminated quoted flow-sequence value")
	}
	parts = append(parts, strings.TrimSpace(current.String()))
	return parts, nil
}

func (n *YAMLNode) Child(key string) *YAMLNode {
	if n == nil || n.Kind != YAMLMapping {
		return nil
	}
	return n.Map[key]
}

func (n *YAMLNode) At(path ...string) *YAMLNode {
	current := n
	for _, key := range path {
		current = current.Child(key)
		if current == nil {
			return nil
		}
	}
	return current
}

func (n *YAMLNode) Scalar() string {
	if n == nil || n.Kind != YAMLScalar {
		return ""
	}
	return n.Value
}

func (n *YAMLNode) Scalars() []string {
	if n == nil {
		return nil
	}
	if n.Kind == YAMLScalar {
		if n.Value == "" {
			return nil
		}
		return []string{n.Value}
	}
	if n.Kind != YAMLSequence {
		return nil
	}
	values := make([]string, 0, len(n.Seq))
	for _, child := range n.Seq {
		if child.Kind == YAMLScalar && child.Value != "" {
			values = append(values, child.Value)
		}
	}
	return values
}
