package common

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
)

type CommandBundle map[string]string

type CommandStatus string

const (
	CommandComplete    CommandStatus = "complete"
	CommandIncomplete  CommandStatus = "incomplete"
	CommandUnsupported CommandStatus = "unsupported"
	CommandFailed      CommandStatus = "failed"
)

type DiagnosticSeverity string

const (
	SeverityInformational DiagnosticSeverity = "informational"
	SeverityWarning       DiagnosticSeverity = "warning"
	SeverityError         DiagnosticSeverity = "error"
)

const (
	EvidenceBundleSchemaVersion  = "iron-atlas.cisco.evidence-bundle.v1"
	CommandEvidenceSchemaVersion = "iron-atlas.cisco.command-evidence.v1"
	EvidenceParserVersion        = "cisco-offline-bundle-parser.v1"
)

const (
	DiagnosticNestedCommand      = "CISCO_BUNDLE_NESTED_COMMAND"
	DiagnosticEndWithoutCommand  = "CISCO_BUNDLE_END_WITHOUT_COMMAND"
	DiagnosticUnclosedCommand    = "CISCO_BUNDLE_UNCLOSED_COMMAND"
	DiagnosticEmptyCommandName   = "CISCO_BUNDLE_EMPTY_COMMAND_NAME"
	DiagnosticTextOutsideCommand = "CISCO_BUNDLE_TEXT_OUTSIDE_COMMAND"
	DiagnosticInputLimit         = "CISCO_BUNDLE_INPUT_LIMIT"
	DiagnosticCommandLimit       = "CISCO_BUNDLE_COMMAND_SIZE_LIMIT"
	DiagnosticLineLimit          = "CISCO_BUNDLE_LINE_SIZE_LIMIT"
	DiagnosticCommandCountLimit  = "CISCO_BUNDLE_COMMAND_COUNT_LIMIT"
	DiagnosticCancelled          = "CISCO_BUNDLE_CANCELLED"
)

var (
	ErrInvalidLimits   = errors.New("invalid Cisco evidence parser limits")
	ErrMalformedBundle = errors.New("malformed Cisco evidence bundle")
	ErrLimitExceeded   = errors.New("Cisco evidence parser limit exceeded")
)

const commandPrefix = "===== COMMAND: "
const commandSuffix = " ====="
const endMarker = "===== END COMMAND ====="

type ParserLimits struct {
	MaxInputBytes   int64
	MaxCommandBytes int64
	MaxLineBytes    int
	MaxCommands     int
}

func DefaultParserLimits() ParserLimits {
	return ParserLimits{
		MaxInputBytes:   32 * 1024 * 1024,
		MaxCommandBytes: 8 * 1024 * 1024,
		MaxLineBytes:    1024 * 1024,
		MaxCommands:     256,
	}
}

func (limits ParserLimits) validate() error {
	if limits.MaxInputBytes <= 0 || limits.MaxCommandBytes <= 0 ||
		limits.MaxLineBytes <= 0 || limits.MaxCommands <= 0 {
		return ErrInvalidLimits
	}
	if limits.MaxCommandBytes > limits.MaxInputBytes ||
		int64(limits.MaxLineBytes) > limits.MaxCommandBytes {
		return ErrInvalidLimits
	}
	return nil
}

type Diagnostic struct {
	Code            string
	Severity        DiagnosticSeverity
	Stage           string
	Line            int
	CommandSequence int
	Detail          string
}

type UploadSafeDiagnostic struct {
	Code            string             `json:"code"`
	Severity        DiagnosticSeverity `json:"severity"`
	Stage           string             `json:"stage"`
	Line            int                `json:"line,omitempty"`
	CommandSequence int                `json:"command_sequence,omitempty"`
}

func (diagnostic Diagnostic) UploadSafe() UploadSafeDiagnostic {
	return UploadSafeDiagnostic{
		Code:            diagnostic.Code,
		Severity:        diagnostic.Severity,
		Stage:           diagnostic.Stage,
		Line:            diagnostic.Line,
		CommandSequence: diagnostic.CommandSequence,
	}
}

type CommandEvidence struct {
	SchemaVersion          string
	Sequence               int
	Command                string
	Output                 string
	NormalizedOutputSHA256 string
	Status                 CommandStatus
	Truncated              bool
	ByteCount              int64
	LineCount              int
}

type EvidenceBundle struct {
	SchemaVersion string
	ParserVersion string
	InputSHA256   string
	BundleSHA256  string
	Complete      bool
	Truncated     bool
	Commands      []CommandEvidence
	Diagnostics   []Diagnostic
}

type UploadSafeCommandEvidence struct {
	SchemaVersion          string        `json:"schema_version"`
	Sequence               int           `json:"sequence"`
	NormalizedOutputSHA256 string        `json:"normalized_output_sha256"`
	Status                 CommandStatus `json:"status"`
	Truncated              bool          `json:"truncated"`
	ByteCount              int64         `json:"byte_count"`
	LineCount              int           `json:"line_count"`
}

type UploadSafeEvidenceBundle struct {
	SchemaVersion string                      `json:"schema_version"`
	ParserVersion string                      `json:"parser_version"`
	InputSHA256   string                      `json:"input_sha256"`
	BundleSHA256  string                      `json:"bundle_sha256"`
	Complete      bool                        `json:"complete"`
	Truncated     bool                        `json:"truncated"`
	Commands      []UploadSafeCommandEvidence `json:"commands"`
	Diagnostics   []UploadSafeDiagnostic      `json:"diagnostics"`
}

func (bundle EvidenceBundle) UploadSafe() UploadSafeEvidenceBundle {
	safe := UploadSafeEvidenceBundle{
		SchemaVersion: bundle.SchemaVersion,
		ParserVersion: bundle.ParserVersion,
		InputSHA256:   bundle.InputSHA256,
		BundleSHA256:  bundle.BundleSHA256,
		Complete:      bundle.Complete,
		Truncated:     bundle.Truncated,
		Commands:      make([]UploadSafeCommandEvidence, 0, len(bundle.Commands)),
		Diagnostics:   make([]UploadSafeDiagnostic, 0, len(bundle.Diagnostics)),
	}
	for _, command := range bundle.Commands {
		safe.Commands = append(safe.Commands, UploadSafeCommandEvidence{
			SchemaVersion:          command.SchemaVersion,
			Sequence:               command.Sequence,
			NormalizedOutputSHA256: command.NormalizedOutputSHA256,
			Status:                 command.Status,
			Truncated:              command.Truncated,
			ByteCount:              command.ByteCount,
			LineCount:              command.LineCount,
		})
	}
	for _, diagnostic := range bundle.Diagnostics {
		safe.Diagnostics = append(safe.Diagnostics, diagnostic.UploadSafe())
	}
	return safe
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (reader contextReader) Read(buffer []byte) (int, error) {
	if err := reader.ctx.Err(); err != nil {
		return 0, err
	}
	count, err := reader.reader.Read(buffer)
	if count == 0 {
		if contextErr := reader.ctx.Err(); contextErr != nil {
			return 0, contextErr
		}
	}
	return count, err
}

type commandBuilder struct {
	evidence CommandEvidence
	body     strings.Builder
	invalid  bool
	overSize bool
}

func ParseEvidenceBundle(ctx context.Context, reader io.Reader, limits ParserLimits) (EvidenceBundle, error) {
	bundle := EvidenceBundle{
		SchemaVersion: EvidenceBundleSchemaVersion,
		ParserVersion: EvidenceParserVersion,
		Complete:      true,
		Commands:      make([]CommandEvidence, 0),
		Diagnostics:   make([]Diagnostic, 0),
	}
	if ctx == nil {
		return bundle, errors.New("Cisco evidence parser context is nil")
	}
	if reader == nil {
		return bundle, errors.New("Cisco evidence reader is nil")
	}
	if err := limits.validate(); err != nil {
		return bundle, err
	}
	if err := ctx.Err(); err != nil {
		bundle.Complete = false
		bundle.Diagnostics = append(bundle.Diagnostics, Diagnostic{
			Code: DiagnosticCancelled, Severity: SeverityError, Stage: "read",
		})
		return bundle, err
	}

	inputHasher := sha256.New()
	limited := &io.LimitedReader{
		R: contextReader{ctx: ctx, reader: reader},
		N: limits.MaxInputBytes + 1,
	}
	scanner := bufio.NewScanner(io.TeeReader(limited, inputHasher))
	initialBuffer := 64 * 1024
	if limits.MaxLineBytes+1 < initialBuffer {
		initialBuffer = limits.MaxLineBytes + 1
	}
	scanner.Buffer(make([]byte, initialBuffer), limits.MaxLineBytes+1)

	var active *commandBuilder
	lineNumber := 0
	malformed := false
	limitExceeded := false
	outsideReported := false
	commandCountReported := false
	ignoreSection := false

	addDiagnostic := func(code string, severity DiagnosticSeverity, stage string, line, sequence int, detail string) {
		bundle.Diagnostics = append(bundle.Diagnostics, Diagnostic{
			Code: code, Severity: severity, Stage: stage, Line: line,
			CommandSequence: sequence, Detail: detail,
		})
	}
	finalize := func(status CommandStatus) {
		if active == nil {
			return
		}
		if active.invalid {
			status = CommandFailed
		}
		if active.overSize {
			status = CommandIncomplete
		}
		active.evidence.Output = strings.TrimRight(active.body.String(), "\n")
		active.evidence.ByteCount = int64(len(active.evidence.Output))
		active.evidence.Status = status
		active.evidence.NormalizedOutputSHA256 = sha256String(active.evidence.Output)
		bundle.Commands = append(bundle.Commands, active.evidence)
		active = nil
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if strings.HasPrefix(line, commandPrefix) && strings.HasSuffix(line, commandSuffix) {
			if active != nil {
				addDiagnostic(DiagnosticNestedCommand, SeverityError, "structure", lineNumber, active.evidence.Sequence, "")
				malformed = true
				bundle.Complete = false
				finalize(CommandIncomplete)
			}
			ignoreSection = false
			if len(bundle.Commands) >= limits.MaxCommands {
				if !commandCountReported {
					addDiagnostic(DiagnosticCommandCountLimit, SeverityError, "limits", lineNumber, 0, "")
					commandCountReported = true
				}
				limitExceeded = true
				bundle.Complete = false
				bundle.Truncated = true
				ignoreSection = true
				continue
			}
			commandName := strings.TrimSuffix(strings.TrimPrefix(line, commandPrefix), commandSuffix)
			active = &commandBuilder{evidence: CommandEvidence{
				SchemaVersion: CommandEvidenceSchemaVersion,
				Sequence:      len(bundle.Commands) + 1,
				Command:       commandName,
				Status:        CommandIncomplete,
			}}
			if strings.TrimSpace(commandName) == "" {
				addDiagnostic(DiagnosticEmptyCommandName, SeverityError, "structure", lineNumber, active.evidence.Sequence, "")
				active.invalid = true
				malformed = true
				bundle.Complete = false
			}
			continue
		}

		if line == endMarker {
			if ignoreSection {
				ignoreSection = false
				continue
			}
			if active == nil {
				addDiagnostic(DiagnosticEndWithoutCommand, SeverityError, "structure", lineNumber, 0, "")
				malformed = true
				bundle.Complete = false
				continue
			}
			finalize(CommandComplete)
			continue
		}

		if ignoreSection {
			continue
		}
		if active == nil {
			if strings.TrimSpace(line) != "" && !outsideReported {
				addDiagnostic(DiagnosticTextOutsideCommand, SeverityWarning, "structure", lineNumber, 0, "")
				outsideReported = true
			}
			continue
		}

		active.evidence.LineCount++
		if active.overSize {
			continue
		}
		prospectiveBytes := int64(active.body.Len() + len(line) + 1)
		if prospectiveBytes > limits.MaxCommandBytes+1 {
			active.overSize = true
			active.evidence.Truncated = true
			limitExceeded = true
			bundle.Complete = false
			bundle.Truncated = true
			addDiagnostic(DiagnosticCommandLimit, SeverityError, "limits", lineNumber, active.evidence.Sequence, "")
			continue
		}
		active.body.WriteString(line)
		active.body.WriteByte('\n')
	}

	scanErr := scanner.Err()
	bundle.InputSHA256 = fmt.Sprintf("%x", inputHasher.Sum(nil))
	if active != nil {
		active.evidence.Truncated = active.evidence.Truncated || scanErr != nil || limited.N == 0
		if scanErr != nil || limited.N == 0 {
			bundle.Truncated = true
		}
		finalize(CommandIncomplete)
	}

	if scanErr != nil {
		bundle.Complete = false
		switch {
		case ctx.Err() != nil:
			addDiagnostic(DiagnosticCancelled, SeverityError, "read", lineNumber+1, 0, "")
			bundle.BundleSHA256 = bundleDigest(bundle)
			return bundle, ctx.Err()
		case limited.N == 0:
			addDiagnostic(DiagnosticInputLimit, SeverityError, "limits", lineNumber+1, 0, "")
			bundle.Truncated = true
			bundle.BundleSHA256 = bundleDigest(bundle)
			return bundle, ErrLimitExceeded
		default:
			addDiagnostic(DiagnosticLineLimit, SeverityError, "limits", lineNumber+1, 0, "")
			bundle.Truncated = true
			bundle.BundleSHA256 = bundleDigest(bundle)
			return bundle, ErrLimitExceeded
		}
	}

	if limited.N == 0 {
		bundle.Complete = false
		bundle.Truncated = true
		addDiagnostic(DiagnosticInputLimit, SeverityError, "limits", lineNumber+1, 0, "")
		bundle.BundleSHA256 = bundleDigest(bundle)
		return bundle, ErrLimitExceeded
	}
	if active != nil {
		panic("Cisco evidence parser retained an active command after finalization")
	}
	if len(bundle.Commands) > 0 {
		last := bundle.Commands[len(bundle.Commands)-1]
		if last.Status == CommandIncomplete && !last.Truncated {
			addDiagnostic(DiagnosticUnclosedCommand, SeverityError, "structure", lineNumber, last.Sequence, "")
			malformed = true
			bundle.Complete = false
		}
	}
	bundle.BundleSHA256 = bundleDigest(bundle)
	if limitExceeded {
		return bundle, ErrLimitExceeded
	}
	if malformed {
		return bundle, ErrMalformedBundle
	}
	return bundle, nil
}

func ParseCommandBundle(reader io.Reader) (CommandBundle, error) {
	evidence, err := ParseEvidenceBundle(context.Background(), reader, DefaultParserLimits())
	if err != nil {
		return nil, err
	}
	bundle := make(CommandBundle, len(evidence.Commands))
	for _, command := range evidence.Commands {
		bundle[command.Command] = command.Output
	}
	return bundle, nil
}

func sha256String(value string) string {
	digest := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", digest)
}

func bundleDigest(bundle EvidenceBundle) string {
	hasher := sha256.New()
	writeHashField(hasher, bundle.SchemaVersion)
	writeHashField(hasher, bundle.ParserVersion)
	writeHashField(hasher, bundle.InputSHA256)
	writeHashField(hasher, fmt.Sprintf("%t", bundle.Complete))
	writeHashField(hasher, fmt.Sprintf("%t", bundle.Truncated))
	for _, command := range bundle.Commands {
		writeHashField(hasher, command.SchemaVersion)
		writeHashField(hasher, fmt.Sprintf("%d", command.Sequence))
		writeHashField(hasher, command.Command)
		writeHashField(hasher, command.NormalizedOutputSHA256)
		writeHashField(hasher, string(command.Status))
		writeHashField(hasher, fmt.Sprintf("%t", command.Truncated))
	}
	for _, diagnostic := range bundle.Diagnostics {
		writeHashField(hasher, diagnostic.Code)
		writeHashField(hasher, string(diagnostic.Severity))
		writeHashField(hasher, diagnostic.Stage)
		writeHashField(hasher, fmt.Sprintf("%d", diagnostic.Line))
		writeHashField(hasher, fmt.Sprintf("%d", diagnostic.CommandSequence))
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func writeHashField(hasher hash.Hash, value string) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = hasher.Write(length[:])
	_, _ = io.WriteString(hasher, value)
}
