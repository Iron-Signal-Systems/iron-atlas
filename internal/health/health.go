package health

import "context"

// Checker reports whether a required application dependency is ready.
type Checker interface {
	Name() string
	Check(context.Context) error
}

// Static is a permanently ready dependency used by the in-memory development store.
type Static struct {
	DependencyName string
}

func (s Static) Name() string {
	if s.DependencyName == "" {
		return "memory"
	}
	return s.DependencyName
}

func (Static) Check(context.Context) error { return nil }
