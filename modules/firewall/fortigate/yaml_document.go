package fortigate

import (
	"context"
	"io"
)

type YAMLKind uint8

const (
	YAMLScalar YAMLKind = iota + 1
	YAMLMapping
	YAMLSequence
)

// YAMLNode is the stable Iron Atlas representation consumed by the FortiGate
// normalizer. The maintained decoder is deliberately kept behind this type.
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

// ParseYAMLDocument preserves the original adapter contract and applies the
// governed default admission limits.
func ParseYAMLDocument(reader io.Reader) (*YAMLDocument, error) {
	return parseYAMLDocumentWithLimits(context.Background(), reader, defaultYAMLDecodeLimits())
}

// ParseYAMLDocumentContext adds cancellation at the bounded-read and Atlas
// node-admission boundaries without exposing decoder implementation details.
func ParseYAMLDocumentContext(ctx context.Context, reader io.Reader) (*YAMLDocument, error) {
	return parseYAMLDocumentWithLimits(ctx, reader, defaultYAMLDecodeLimits())
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
