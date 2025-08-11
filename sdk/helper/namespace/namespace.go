package namespace

import (
	"context"
	"errors"
	"fmt"
)

type contextValues struct{}

type Namespace struct {
	ID             string            `json:"id" mapstructure:"id"`
	Path           string            `json:"path" mapstructure:"path"`
	CustomMetadata map[string]string `json:"custom_metadata" mapstructure:"custom_metadata"`
}

func (n *Namespace) String() string {
	return fmt.Sprintf("ID: %s. Path: %s", n.ID, n.Path)
}

const (
	RootNamespaceID = "root"
)

var (
	contextNamespace contextValues = struct{}{}
	ErrNoNamespace   error         = errors.New("no namespace")
	RootNamespace    *Namespace    = &Namespace{
		ID:             RootNamespaceID,
		Path:           "",
		CustomMetadata: make(map[string]string),
	}
)

// FromContext retrieves the namespace from a context, or an error
// if there is no namespace in the context.
func FromContext(ctx context.Context) (*Namespace, error) {
	if ctx == nil {
		return nil, errors.New("context was nil")
	}

	nsRaw := ctx.Value(contextNamespace)
	if nsRaw == nil {
		return nil, ErrNoNamespace
	}

	ns := nsRaw.(*Namespace)
	if ns == nil {
		return nil, ErrNoNamespace
	}

	return ns, nil
}

func ContextWithNamespace(ctx context.Context, ns *Namespace) context.Context {
	return context.WithValue(ctx, contextNamespace, ns)
}
