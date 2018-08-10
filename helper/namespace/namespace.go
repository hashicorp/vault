package namespace

import (
	"context"
	"errors"
	"strings"
)

type nsContext struct {
	context.Context
	// Note: this is currently not locked because we think all uses will take
	// place within a single goroutine. If that isn't the case, this should be
	// protected by an atomic.Value.
	cachedNS *Namespace
}

type contextValues struct{}

const (
	RootNamespaceID = "root"
)

var (
	contextNamespace contextValues = struct{}{}
	ErrNoNamespace   error         = errors.New("no namespace")
)

type Namespace struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

func New(id, path string) *Namespace {
	return &Namespace{
		ID:   id,
		Path: path,
	}
}

func (n *Namespace) HasParent(possibleParent *Namespace) bool {
	switch {
	case n.Path == "":
		return false
	case possibleParent.Path == "":
		return true
	default:
		return strings.HasPrefix(n.Path, possibleParent.Path)
	}
}

func (n *Namespace) TrimmedPath(path string) string {
	return strings.TrimPrefix(path, n.Path)
}

func ContextWithNamespace(ctx context.Context, ns *Namespace) context.Context {
	nsCtx := context.WithValue(ctx, contextNamespace, ns)
	return &nsContext{
		Context:  nsCtx,
		cachedNS: ns,
	}
}

func FromContext(ctx context.Context) (*Namespace, error) {
	if ctx == nil {
		return nil, errors.New("context was nil")
	}

	nsCtx, ok := ctx.(*nsContext)
	if ok {
		if nsCtx.cachedNS != nil {
			return nsCtx.cachedNS, nil
		}
	}

	ns := ctx.Value(contextNamespace)
	if ns == nil {
		return nil, ErrNoNamespace
	}

	if ok {
		nsCtx.cachedNS = ns.(*Namespace)
	}

	return ns.(*Namespace), nil
}

func TestContext() context.Context {
	return ContextWithNamespace(context.Background(), New(RootNamespaceID, ""))
}

// Canonicalize trims any prefix '/' and adds a trailing '/' to the
// provided string
func Canonicalize(nsPath string) string {
	if nsPath == "" {
		return ""
	}

	// Canonicalize the path to not have a '/' prefix
	nsPath = strings.TrimPrefix(nsPath, "/")

	// Canonicalize the path to always having a '/' suffix
	if !strings.HasSuffix(nsPath, "/") {
		nsPath += "/"
	}

	return nsPath
}
