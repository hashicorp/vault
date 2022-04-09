package namespace

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

type contextValues struct{}

type Namespace struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

const (
	RootNamespaceID = "root"
)

var (
	contextNamespace contextValues = struct{}{}
	ErrNoNamespace   error         = errors.New("no namespace")
	RootNamespace    *Namespace    = &Namespace{
		ID:   RootNamespaceID,
		Path: "",
	}
)

func (n *Namespace) HasParent(possibleParent *Namespace) bool {
	switch {
	case possibleParent.Path == "":
		return true
	case n.Path == "":
		return false
	default:
		return strings.HasPrefix(n.Path, possibleParent.Path)
	}
}

func (n *Namespace) TrimmedPath(path string) string {
	return strings.TrimPrefix(path, n.Path)
}

func ContextWithNamespace(ctx context.Context, ns *Namespace) context.Context {
	return context.WithValue(ctx, contextNamespace, ns)
}

func RootContext(ctx context.Context) context.Context {
	if ctx == nil {
		return ContextWithNamespace(context.Background(), RootNamespace)
	}
	return ContextWithNamespace(ctx, RootNamespace)
}

// This function caches the ns to avoid doing a .Value lookup over and over,
// because it's called a *lot* in the request critical path. .Value is
// concurrency-safe so uses some kind of locking/atomicity, but it should never
// be read before first write, plus we don't believe this will be called from
// different goroutines, so it should be safe.
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

func SplitIDFromString(input string) (string, string) {
	prefix := ""
	slashIdx := strings.LastIndex(input, "/")

	switch {
	case strings.HasPrefix(input, consts.LegacyBatchTokenPrefix):
		prefix = consts.LegacyBatchTokenPrefix
		input = input[2:]

	case strings.HasPrefix(input, consts.LegacyServiceTokenPrefix):
		prefix = consts.LegacyServiceTokenPrefix
		input = input[2:]
	case strings.HasPrefix(input, consts.BatchTokenPrefix):
		prefix = consts.BatchTokenPrefix
		input = input[4:]
	case strings.HasPrefix(input, consts.ServiceTokenPrefix):
		prefix = consts.ServiceTokenPrefix
		input = input[4:]

	case slashIdx > 0:
		// Leases will never have a b./s. to start
		if slashIdx == len(input)-1 {
			return input, ""
		}
		prefix = input[:slashIdx+1]
		input = input[slashIdx+1:]
	}

	idx := strings.LastIndex(input, ".")
	if idx == -1 {
		return prefix + input, ""
	}
	if idx == len(input)-1 {
		return prefix + input, ""
	}

	return prefix + input[:idx], input[idx+1:]
}

// MountPathDetails contains the details of a mount's location,
// consisting of the namespace of the mount and the path of the
// mount within the namespace
type MountPathDetails struct {
	Namespace *Namespace
	MountPath string
}

func (mpd *MountPathDetails) GetRelativePath(currNs *Namespace) string {
	subNsPath := strings.TrimPrefix(mpd.Namespace.Path, currNs.Path)
	return subNsPath + mpd.MountPath
}

func (mpd *MountPathDetails) GetFullPath() string {
	return mpd.Namespace.Path + mpd.MountPath
}
