package backend

import (
	"regexp"
	"sync"

	"github.com/hashicorp/vault/vault"
)

// Backend is an implementation of vault.LogicalBackend that allows
// the implementer to code a backend using a much more programmer-friendly
// framework that handles a lot of the routing and validation for you.
//
// This is recommended over implementing vault.LogicalBackend directly.
type Backend struct {
	// Paths are the various routes that the backend responds to.
	// This cannot be modified after construction (i.e. dynamically changing
	// paths, including adding or removing, is not allowed once the
	// backend is in use).
	Paths []*Path

	once    sync.Once
	pathsRe []*regexp.Regexp
}

// Path is a single path that the backend responds to.
type Path struct {
	// Pattern is the pattern of the URL that matches this path.
	//
	// This should be a valid regular expression. Named captures will be
	// exposed as fields that should map to a schema in Fields. If a named
	// capture is not a field in the Fields map, then it will be ignored.
	Pattern string

	// Fields is the mapping of data fields to a schema describing that
	// field. Named captures in the Pattern also map to fields. If a named
	// capture name matches a PUT body name, the named capture takes
	// priority.
	//
	// Note that only named capture fields are available in every operation,
	// whereas all fields are avaiable in the Write operation.
	Fields map[string]*FieldSchema

	// Root if not blank, denotes that this path requires root
	// privileges and the path pattern that is the root path. This can't
	// be a regular expression and must be an exact path. It may have a
	// trailing '*' to denote that it is a prefix, and not an exact match.
	Root string

	// Callback is what is called when this path is requested with
	// a valid set of data.
	Callback func(*vault.Request, *FieldData) (*vault.Response, error)
}

// vault.LogicalBackend impl.
func (b *Backend) HandleRequest(req *vault.Request) (*vault.Response, error) {
	// Find the matching route
	path, captures := b.route(req.Path)
	if path == nil {
		return nil, vault.ErrUnsupportedPath
	}

	// Build up the data for the route, with the URL taking priority
	// for the fields over the PUT data.
	raw := make(map[string]interface{}, len(path.Fields))
	for k, v := range req.Data {
		raw[k] = v
	}
	for k, v := range captures {
		raw[k] = v
	}

	// Call the callback with the request and the data
	return path.Callback(req, &FieldData{
		Raw:    raw,
		Schema: path.Fields,
	})
}

// vault.LogicalBackend impl.
func (b *Backend) RootPaths() []string {
	// TODO
	return nil
}

// Route looks up the path that would be used for a given path string.
func (b *Backend) Route(path string) *Path {
	result, _ := b.route(path)
	return result
}

func (b *Backend) init() {
	b.pathsRe = make([]*regexp.Regexp, len(b.Paths))
	for i, p := range b.Paths {
		b.pathsRe[i] = regexp.MustCompile(p.Pattern)
	}
}

func (b *Backend) route(path string) (*Path, map[string]string) {
	b.once.Do(b.init)

	for i, re := range b.pathsRe {
		matches := re.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		// We have a match, determine the mapping of the captures and
		// store that for returning.
		var captures map[string]string
		path := b.Paths[i]
		if captureNames := re.SubexpNames(); len(captureNames) > 1 {
			captures = make(map[string]string, len(captureNames))
			for i, name := range captureNames {
				if name != "" {
					captures[name] = matches[i]
				}
			}
		}

		return path, captures
	}

	return nil, nil
}

// FieldSchema is a basic schema to describe the format of a path field.
type FieldSchema struct {
	Type    FieldType
	Default interface{}
}

// DefaultOrZero returns the default value if it is set, or otherwise
// the zero value of the type.
func (s *FieldSchema) DefaultOrZero() interface{} {
	if s.Default != nil {
		return s.Default
	}

	return s.Type.Zero()
}

func (t FieldType) Zero() interface{} {
	switch t {
	case TypeString:
		return ""
	case TypeInt:
		return 0
	case TypeBool:
		return false
	default:
		panic("unknown type: " + t.String())
	}
}
