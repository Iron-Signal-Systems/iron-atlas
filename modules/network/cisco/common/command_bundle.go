package common

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CommandBundle map[string]string

const commandPrefix = "===== COMMAND: "
const commandSuffix = " ====="
const endMarker = "===== END COMMAND ====="

func ParseCommandBundle(reader io.Reader) (CommandBundle, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 32*1024*1024)
	bundle := make(CommandBundle)
	var command string
	var body strings.Builder
	lineNo := 0
	flush := func() {
		if command != "" {
			bundle[command] = strings.TrimRight(body.String(), "\n")
			command = ""
			body.Reset()
		}
	}
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.HasPrefix(line, commandPrefix) && strings.HasSuffix(line, commandSuffix) {
			if command != "" {
				return nil, fmt.Errorf("line %d: nested command section", lineNo)
			}
			command = strings.TrimSuffix(strings.TrimPrefix(line, commandPrefix), commandSuffix)
			continue
		}
		if line == endMarker {
			if command == "" {
				return nil, fmt.Errorf("line %d: end marker without command", lineNo)
			}
			flush()
			continue
		}
		if command != "" {
			body.WriteString(line)
			body.WriteByte('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if command != "" {
		return nil, fmt.Errorf("command section %q is not closed", command)
	}
	return bundle, nil
}
