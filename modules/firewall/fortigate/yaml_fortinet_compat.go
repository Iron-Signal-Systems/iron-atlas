package fortigate

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

var (
	fortinetMappingValueLine = regexp.MustCompile(
		`^(?P<prefix>[ ]*(-[ ]+)?[A-Za-z0-9_-]+:[ ]*)(?P<value>.*)$`,
	)
	fortinetPrefixSubexpression = fortinetMappingValueLine.SubexpIndex("prefix")
	fortinetValueSubexpression  = fortinetMappingValueLine.SubexpIndex("value")
)

type fortinetCompatibilityFragment struct {
	raw    []byte
	quoted bool
}

type fortinetCompatibilityState struct {
	limits    yamlDecodeLimits
	rewrites  int
	fragments int
	quote     byte
}

// normalizeFortinetYAMLCompatibility repairs documented Fortinet export defect
// families before maintained YAML decoding: adjacent multi-value fragments
// beginning with a double-quoted value become a flow sequence, and restricted
// literal object-name keys beginning with YAML indicators are quoted. It
// preserves physical line count and does not decode or log source values.
func normalizeFortinetYAMLCompatibility(data []byte, limits yamlDecodeLimits) ([]byte, error) {
	state := fortinetCompatibilityState{limits: limits}
	var output bytes.Buffer
	output.Grow(len(data))
	for index, physicalLine := range bytes.SplitAfter(data, []byte{'\n'}) {
		line, ending := splitYAMLPhysicalLine(physicalLine)
		original := line
		if state.quote == 0 {
			var err error
			line, err = state.rewriteUnsafeMappingKey(line, index+1)
			if err != nil {
				return nil, err
			}
			line, err = state.rewriteAdjacentQuotedValue(line, index+1)
			if err != nil {
				return nil, err
			}
		}
		state.quote = updateYAMLQuotedScalarState(original, state.quote)
		if int64(output.Len())+int64(len(line))+int64(len(ending)) > limits.MaxInputBytes {
			return nil, errors.New("YAML admission rejected: MaxInputBytes limit exceeded after Fortinet compatibility rewrite")
		}
		output.Write(line)
		output.Write(ending)
	}
	return output.Bytes(), nil
}

func splitYAMLPhysicalLine(physicalLine []byte) (line []byte, ending []byte) {
	if bytes.HasSuffix(physicalLine, []byte{'\n'}) {
		line = physicalLine[:len(physicalLine)-1]
		ending = []byte{'\n'}
		if bytes.HasSuffix(line, []byte{'\r'}) {
			line = line[:len(line)-1]
			ending = []byte{'\r', '\n'}
		}
		return line, ending
	}
	return physicalLine, nil
}

func (state *fortinetCompatibilityState) rewriteAdjacentQuotedValue(line []byte, lineNumber int) ([]byte, error) {
	match := fortinetMappingValueLine.FindSubmatchIndex(line)
	if match == nil {
		return line, nil
	}
	prefixStart, prefixEnd := subexpressionBounds(match, fortinetPrefixSubexpression)
	valueStart, valueEnd := subexpressionBounds(match, fortinetValueSubexpression)
	fragments, suffix, compatible := parseFortinetAdjacentValue(line[valueStart:valueEnd])
	if !compatible {
		return line, nil
	}

	if err := state.recordCompatibilityRewrite(lineNumber, len(fragments)); err != nil {
		return nil, err
	}

	var rewritten bytes.Buffer
	rewritten.Grow(len(line) + 2*len(fragments))
	rewritten.Write(line[prefixStart:prefixEnd])
	rewritten.WriteByte('[')
	for index, fragment := range fragments {
		if index > 0 {
			rewritten.WriteString(", ")
		}
		if fragment.quoted {
			rewritten.Write(fragment.raw)
			continue
		}
		rewritten.WriteByte('"')
		rewritten.Write(fragment.raw)
		rewritten.WriteByte('"')
	}
	rewritten.WriteByte(']')
	rewritten.Write(suffix)
	return rewritten.Bytes(), nil
}

func (state *fortinetCompatibilityState) rewriteUnsafeMappingKey(line []byte, lineNumber int) ([]byte, error) {
	prefixEnd, keyStart, keyEnd, suffixStart, compatible := splitFortinetUnsafeMappingKey(line)
	if !compatible {
		return line, nil
	}
	if err := state.recordCompatibilityRewrite(lineNumber, 1); err != nil {
		return nil, err
	}

	var rewritten bytes.Buffer
	rewritten.Grow(len(line) + 2)
	rewritten.Write(line[:prefixEnd])
	rewritten.WriteByte('"')
	rewritten.Write(line[keyStart:keyEnd])
	rewritten.WriteByte('"')
	rewritten.Write(line[suffixStart:])
	return rewritten.Bytes(), nil
}

func splitFortinetUnsafeMappingKey(line []byte) (
	prefixEnd int,
	keyStart int,
	keyEnd int,
	suffixStart int,
	compatible bool,
) {
	cursor := 0
	for cursor < len(line) && line[cursor] == ' ' {
		cursor++
	}
	if cursor < len(line) && line[cursor] == '-' {
		cursor++
		spaceStart := cursor
		for cursor < len(line) && line[cursor] == ' ' {
			cursor++
		}
		if cursor == spaceStart {
			return 0, 0, 0, 0, false
		}
	}
	prefixEnd = cursor
	keyStart = cursor
	if keyStart >= len(line) || !isUnsafeYAMLKeyIndicator(line[keyStart]) {
		return 0, 0, 0, 0, false
	}

	colonOffset := bytes.IndexByte(line[keyStart:], ':')
	if colonOffset <= 1 {
		return 0, 0, 0, 0, false
	}
	keyEnd = keyStart + colonOffset
	for _, value := range line[keyStart:keyEnd] {
		if !isSafeFortinetLiteralKeyByte(value) {
			return 0, 0, 0, 0, false
		}
	}

	suffixStart = keyEnd
	cursor = keyEnd + 1
	for cursor < len(line) && line[cursor] == ' ' {
		cursor++
	}
	if cursor == len(line) || line[cursor] == '#' {
		return prefixEnd, keyStart, keyEnd, suffixStart, true
	}
	return 0, 0, 0, 0, false
}

func isUnsafeYAMLKeyIndicator(value byte) bool {
	switch value {
	case '*', '&', '!', '%', '@':
		return true
	default:
		return false
	}
}

func isSafeFortinetLiteralKeyByte(value byte) bool {
	if value < 0x21 || value > 0x7e {
		return false
	}
	switch value {
	case '"', '\'', '\\', '[', ']', '{', '}', ',', '#', ':':
		return false
	default:
		return true
	}
}

func (state *fortinetCompatibilityState) recordCompatibilityRewrite(lineNumber int, fragmentCount int) error {
	state.rewrites++
	if state.rewrites > state.limits.MaxCompatibilityRewrites {
		return fmt.Errorf(
			"YAML admission rejected at line %d: MaxCompatibilityRewrites limit exceeded",
			lineNumber,
		)
	}
	state.fragments += fragmentCount
	if state.fragments > state.limits.MaxCompatibilityFragments {
		return fmt.Errorf(
			"YAML admission rejected at line %d: MaxCompatibilityFragments limit exceeded",
			lineNumber,
		)
	}
	return nil
}

func parseFortinetAdjacentValue(value []byte) ([]fortinetCompatibilityFragment, []byte, bool) {
	if len(value) == 0 || value[0] != '"' {
		return nil, nil, false
	}

	fragments := make([]fortinetCompatibilityFragment, 0, 4)
	for cursor := 0; cursor < len(value); {
		fragment := fortinetCompatibilityFragment{}
		if value[cursor] == '"' {
			end, ok := scanFortinetDoubleQuotedFragment(value, cursor)
			if !ok {
				return nil, nil, false
			}
			fragment.raw = value[cursor:end]
			fragment.quoted = true
			cursor = end
		} else {
			end := cursor
			for end < len(value) && value[end] != ' ' {
				end++
			}
			fragment.raw = value[cursor:end]
			if !isSafeFortinetBareFragment(fragment.raw) {
				return nil, nil, false
			}
			cursor = end
		}
		fragments = append(fragments, fragment)

		if cursor == len(value) {
			return fragments, nil, len(fragments) >= 2
		}
		spaceStart := cursor
		for cursor < len(value) && value[cursor] == ' ' {
			cursor++
		}
		if cursor == spaceStart {
			return nil, nil, false
		}
		if cursor == len(value) || value[cursor] == '#' {
			return fragments, value[spaceStart:], len(fragments) >= 2
		}
	}
	return nil, nil, false
}

func scanFortinetDoubleQuotedFragment(value []byte, start int) (int, bool) {
	for index := start + 1; index < len(value); index++ {
		switch value[index] {
		case '\\':
			index++
		case '"':
			return index + 1, true
		}
	}
	return 0, false
}

func isSafeFortinetBareFragment(value []byte) bool {
	if len(value) == 0 || !isASCIIAlphaNumeric(value[0]) {
		return false
	}
	for _, current := range value {
		if isASCIIAlphaNumeric(current) {
			continue
		}
		switch current {
		case '-', '_', '.', '/', '+', ':', '@', '%':
			continue
		default:
			return false
		}
	}
	return true
}

func isASCIIAlphaNumeric(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9'
}

func subexpressionBounds(match []int, subexpression int) (int, int) {
	offset := subexpression * 2
	return match[offset], match[offset+1]
}

func updateYAMLQuotedScalarState(line []byte, quote byte) byte {
	for index := 0; index < len(line); index++ {
		value := line[index]
		switch quote {
		case '\'':
			if value != '\'' {
				continue
			}
			if index+1 < len(line) && line[index+1] == '\'' {
				index++
				continue
			}
			quote = 0
		case '"':
			if value == '\\' {
				index++
				continue
			}
			if value == '"' {
				quote = 0
			}
		default:
			if value == '#' && (index == 0 || line[index-1] == ' ' || line[index-1] == '\t') {
				return 0
			}
			if (value == '\'' || value == '"') && yamlQuoteCanStart(line, index) {
				quote = value
			}
		}
	}
	return quote
}

func yamlQuoteCanStart(line []byte, index int) bool {
	for index--; index >= 0; index-- {
		switch line[index] {
		case ' ', '\t':
			continue
		case ':', '-', '?', '[', '{', ',':
			return true
		default:
			return false
		}
	}
	return true
}
