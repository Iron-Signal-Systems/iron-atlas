package fortigate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	yaml "go.yaml.in/yaml/v4"
)

var (
	errYAMLDecoderDepth = errors.New("YAML decoder combined-depth limit exceeded")
	errYAMLDecoderAlias = errors.New("YAML aliases are not admitted")
	yamlLinePattern     = regexp.MustCompile(`(?i)line\s+([0-9]+)`)
	yamlColumnPattern   = regexp.MustCompile(`(?i)column\s+([0-9]+)`)
)

type yamlDecoderLimitPlugin struct {
	maxDepth int
}

func (plugin *yamlDecoderLimitPlugin) CheckDepth(depth int, _ *yaml.DepthContext) error {
	// The decoder counts stream/document structure in addition to collection
	// nodes. Atlas performs the exact contract check during node conversion;
	// this small allowance prevents hostile recursion before that pass without
	// making decoder-internal depth accounting part of the Atlas contract.
	if depth > plugin.maxDepth+16 {
		return errYAMLDecoderDepth
	}
	return nil
}

func (*yamlDecoderLimitPlugin) CheckAlias(aliasCount, _ int) error {
	if aliasCount > 0 {
		return errYAMLDecoderAlias
	}
	return nil
}

type yamlAdmissionState struct {
	limits               yamlDecodeLimits
	nodes                int
	totalMappingEntries  int
	totalSequenceEntries int
}

func parseYAMLDocumentWithLimits(ctx context.Context, reader io.Reader, limits yamlDecodeLimits) (*YAMLDocument, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if reader == nil {
		return nil, errors.New("reader is required")
	}
	if err := limits.validate(); err != nil {
		return nil, err
	}

	data, err := readBoundedYAML(ctx, reader, limits.MaxInputBytes)
	if err != nil {
		return nil, err
	}
	data, err = normalizeFortinetYAMLCompatibility(data, limits)
	if err != nil {
		return nil, err
	}

	var source yaml.Node
	err = yaml.Load(
		data,
		&source,
		yaml.WithV4Defaults(),
		yaml.WithUniqueKeys(true),
		yaml.WithPlugin(&yamlDecoderLimitPlugin{maxDepth: limits.MaxCombinedDepth}),
	)
	if err != nil {
		return nil, sanitizeYAMLDecodeError(err)
	}
	if source.Kind != yaml.DocumentNode || len(source.Content) != 1 || source.Content[0] == nil {
		return nil, errors.New("empty YAML document")
	}

	state := &yamlAdmissionState{limits: limits}
	root, err := state.convertNode(ctx, source.Content[0], 0, 0, 0, false)
	if err != nil {
		return nil, err
	}
	return &YAMLDocument{
		Root:     root,
		Comments: collectYAMLComments(&source),
	}, nil
}

func readBoundedYAML(ctx context.Context, reader io.Reader, maxBytes int64) ([]byte, error) {
	limited := io.LimitReader(&contextReader{ctx: ctx, reader: reader}, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, errors.New("YAML admission rejected: MaxInputBytes limit exceeded")
	}
	return data, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (reader *contextReader) Read(buffer []byte) (int, error) {
	select {
	case <-reader.ctx.Done():
		return 0, reader.ctx.Err()
	default:
		return reader.reader.Read(buffer)
	}
}

func sanitizeYAMLDecodeError(err error) error {
	if errors.Is(err, errYAMLDecoderDepth) {
		return errYAMLDecoderDepth
	}
	if errors.Is(err, errYAMLDecoderAlias) {
		return errYAMLDecoderAlias
	}
	if loadError := firstStructuredYAMLLoadError(err); loadError != nil {
		stage := admittedYAMLLoadStage(loadError.Stage)
		switch {
		case loadError.Mark.Line > 0 && loadError.Mark.Column > 0:
			return fmt.Errorf(
				"YAML syntax decoding failed in %s stage at line %d, column %d",
				stage,
				loadError.Mark.Line,
				loadError.Mark.Column,
			)
		case loadError.Mark.Line > 0:
			return fmt.Errorf(
				"YAML syntax decoding failed in %s stage at line %d",
				stage,
				loadError.Mark.Line,
			)
		default:
			return fmt.Errorf("YAML syntax decoding failed in %s stage", stage)
		}
	}

	message := err.Error()
	line := firstYAMLPosition(yamlLinePattern, message)
	column := firstYAMLPosition(yamlColumnPattern, message)
	switch {
	case line > 0 && column > 0:
		return fmt.Errorf("YAML syntax decoding failed at line %d, column %d", line, column)
	case line > 0:
		return fmt.Errorf("YAML syntax decoding failed at line %d", line)
	default:
		return errors.New("YAML syntax decoding failed")
	}
}

func firstStructuredYAMLLoadError(err error) *yaml.LoadError {
	var loadError *yaml.LoadError
	if errors.As(err, &loadError) {
		return loadError
	}
	return nil
}

func admittedYAMLLoadStage(stage yaml.Stage) string {
	switch stage {
	case yaml.ReaderStage:
		return "reader"
	case yaml.ScannerStage:
		return "scanner"
	case yaml.ParserStage:
		return "parser"
	case yaml.ComposerStage:
		return "composer"
	case yaml.ResolverStage:
		return "resolver"
	case yaml.ConstructorStage:
		return "constructor"
	default:
		return "decoder"
	}
}

func firstYAMLPosition(pattern *regexp.Regexp, message string) int {
	match := pattern.FindStringSubmatch(message)
	if len(match) != 2 {
		return 0
	}
	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0
	}
	return value
}

func (state *yamlAdmissionState) convertNode(
	ctx context.Context,
	source *yaml.Node,
	mappingDepth int,
	sequenceDepth int,
	combinedDepth int,
	implicitEmptyMapping bool,
) (*YAMLNode, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("YAML admission rejected: nil node")
	}
	if err := state.addNode(source); err != nil {
		return nil, err
	}
	if source.Anchor != "" {
		return nil, yamlAdmissionError(source, "anchors are not admitted")
	}
	if source.Kind == yaml.AliasNode || source.Alias != nil {
		return nil, yamlAdmissionError(source, "aliases are not admitted")
	}
	if !allowedYAMLTag(source.Kind, source.Tag) {
		return nil, yamlAdmissionError(source, "custom tags are not admitted")
	}

	switch source.Kind {
	case yaml.MappingNode:
		return state.convertMapping(ctx, source, mappingDepth+1, sequenceDepth, combinedDepth+1)
	case yaml.SequenceNode:
		return state.convertSequence(ctx, source, mappingDepth, sequenceDepth+1, combinedDepth+1)
	case yaml.ScalarNode:
		if source.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
			return nil, yamlAdmissionError(source, "block scalars are not admitted")
		}
		if len(source.Value) > state.limits.MaxScalarBytes {
			return nil, yamlLimitError(source, "MaxScalarBytes")
		}
		if implicitEmptyMapping && isYAMLNullTag(source.Tag) && source.Value == "" {
			return &YAMLNode{
				Kind:   YAMLMapping,
				Map:    make(map[string]*YAMLNode),
				Line:   source.Line,
				Column: source.Column,
			}, nil
		}
		return &YAMLNode{
			Kind:   YAMLScalar,
			Value:  source.Value,
			Line:   source.Line,
			Column: source.Column,
		}, nil
	default:
		return nil, yamlAdmissionError(source, "unsupported node kind")
	}
}

func (state *yamlAdmissionState) convertMapping(
	ctx context.Context,
	source *yaml.Node,
	mappingDepth int,
	sequenceDepth int,
	combinedDepth int,
) (*YAMLNode, error) {
	if source.Style&yaml.FlowStyle != 0 {
		return nil, yamlAdmissionError(source, "flow mappings are not admitted")
	}
	if mappingDepth > state.limits.MaxMappingDepth {
		return nil, yamlLimitError(source, "MaxMappingDepth")
	}
	if combinedDepth > state.limits.MaxCombinedDepth {
		return nil, yamlLimitError(source, "MaxCombinedDepth")
	}
	if len(source.Content)%2 != 0 {
		return nil, yamlAdmissionError(source, "mapping node has an invalid shape")
	}
	entries := len(source.Content) / 2
	if entries > state.limits.MaxMappingEntries {
		return nil, yamlLimitError(source, "MaxMappingEntries")
	}
	state.totalMappingEntries += entries
	if state.totalMappingEntries > state.limits.MaxTotalMappingEntries {
		return nil, yamlLimitError(source, "MaxTotalMappingEntries")
	}

	target := &YAMLNode{
		Kind:   YAMLMapping,
		Map:    make(map[string]*YAMLNode, entries),
		Order:  make([]string, 0, entries),
		Line:   source.Line,
		Column: source.Column,
	}
	for index := 0; index < len(source.Content); index += 2 {
		keyNode := source.Content[index]
		valueNode := source.Content[index+1]
		key, err := state.convertMappingKey(ctx, keyNode)
		if err != nil {
			return nil, err
		}
		if _, duplicate := target.Map[key]; duplicate {
			return nil, yamlAdmissionError(keyNode, "duplicate mapping key")
		}
		value, err := state.convertNode(
			ctx,
			valueNode,
			mappingDepth,
			sequenceDepth,
			combinedDepth,
			true,
		)
		if err != nil {
			return nil, err
		}
		target.Map[key] = value
		target.Order = append(target.Order, key)
	}
	return target, nil
}

func (state *yamlAdmissionState) convertMappingKey(ctx context.Context, source *yaml.Node) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if source == nil {
		return "", errors.New("YAML admission rejected: nil mapping key")
	}
	if err := state.addNode(source); err != nil {
		return "", err
	}
	if source.Kind != yaml.ScalarNode {
		return "", yamlAdmissionError(source, "mapping keys must be scalars")
	}
	if source.Anchor != "" || source.Alias != nil {
		return "", yamlAdmissionError(source, "anchors and aliases are not admitted")
	}
	if !allowedYAMLTag(source.Kind, source.Tag) {
		return "", yamlAdmissionError(source, "custom tags are not admitted")
	}
	if source.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
		return "", yamlAdmissionError(source, "block scalar keys are not admitted")
	}
	if source.Value == "" {
		return "", yamlAdmissionError(source, "mapping keys must not be empty")
	}
	if len(source.Value) > state.limits.MaxKeyBytes {
		return "", yamlLimitError(source, "MaxKeyBytes")
	}
	return source.Value, nil
}

func (state *yamlAdmissionState) convertSequence(
	ctx context.Context,
	source *yaml.Node,
	mappingDepth int,
	sequenceDepth int,
	combinedDepth int,
) (*YAMLNode, error) {
	if sequenceDepth > state.limits.MaxSequenceDepth {
		return nil, yamlLimitError(source, "MaxSequenceDepth")
	}
	if combinedDepth > state.limits.MaxCombinedDepth {
		return nil, yamlLimitError(source, "MaxCombinedDepth")
	}
	entries := len(source.Content)
	if entries > state.limits.MaxSequenceEntries {
		return nil, yamlLimitError(source, "MaxSequenceEntries")
	}
	state.totalSequenceEntries += entries
	if state.totalSequenceEntries > state.limits.MaxTotalSequenceEntries {
		return nil, yamlLimitError(source, "MaxTotalSequenceEntries")
	}

	target := &YAMLNode{
		Kind:   YAMLSequence,
		Seq:    make([]*YAMLNode, 0, entries),
		Line:   source.Line,
		Column: source.Column,
	}
	for _, child := range source.Content {
		value, err := state.convertNode(
			ctx,
			child,
			mappingDepth,
			sequenceDepth,
			combinedDepth,
			false,
		)
		if err != nil {
			return nil, err
		}
		target.Seq = append(target.Seq, value)
	}
	return target, nil
}

func (state *yamlAdmissionState) addNode(source *yaml.Node) error {
	state.nodes++
	if state.nodes > state.limits.MaxNodes {
		return yamlLimitError(source, "MaxNodes")
	}
	return nil
}

func allowedYAMLTag(kind yaml.Kind, tag string) bool {
	if tag == "" {
		return true
	}
	switch kind {
	case yaml.MappingNode:
		return tag == "!!map" || tag == "tag:yaml.org,2002:map"
	case yaml.SequenceNode:
		return tag == "!!seq" || tag == "tag:yaml.org,2002:seq"
	case yaml.ScalarNode:
		switch tag {
		case "!!binary", "!!bool", "!!float", "!!int", "!!null", "!!str", "!!timestamp",
			"tag:yaml.org,2002:binary", "tag:yaml.org,2002:bool", "tag:yaml.org,2002:float",
			"tag:yaml.org,2002:int", "tag:yaml.org,2002:null", "tag:yaml.org,2002:str",
			"tag:yaml.org,2002:timestamp":
			return true
		}
	}
	return false
}

func isYAMLNullTag(tag string) bool {
	return tag == "!!null" || tag == "tag:yaml.org,2002:null"
}

func yamlAdmissionError(node *yaml.Node, reason string) error {
	return fmt.Errorf(
		"YAML admission rejected at line %d, column %d: %s",
		node.Line,
		node.Column,
		reason,
	)
}

func yamlLimitError(node *yaml.Node, limitName string) error {
	return fmt.Errorf(
		"YAML admission rejected at line %d, column %d: %s limit exceeded",
		node.Line,
		node.Column,
		limitName,
	)
}

func collectYAMLComments(root *yaml.Node) []string {
	if root == nil {
		return nil
	}
	comments := make([]string, 0)
	var visit func(*yaml.Node)
	visit = func(node *yaml.Node) {
		if node == nil {
			return
		}
		appendYAMLComment(&comments, node.HeadComment)
		appendYAMLComment(&comments, node.LineComment)
		for _, child := range node.Content {
			visit(child)
		}
		appendYAMLComment(&comments, node.FootComment)
	}
	visit(root)
	return comments
}

func appendYAMLComment(comments *[]string, value string) {
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			line = "# " + line
		}
		*comments = append(*comments, line)
	}
}
