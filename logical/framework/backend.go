package framework

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/logical"
)

// Backend is an implementation of logical.Backend that allows
// the implementer to code a backend using a much more programmer-friendly
// framework that handles a lot of the routing and validation for you.
//
// This is recommended over implementing logical.Backend directly.
type Backend struct {
	// Help is the help text that is shown when a help request is made
	// on the root of this resource. The root help is special since we
	// show all the paths that can be requested.
	Help string

	// Paths are the various routes that the backend responds to.
	// This cannot be modified after construction (i.e. dynamically changing
	// paths, including adding or removing, is not allowed once the
	// backend is in use).
	//
	// PathsSpecial is the list of path patterns that denote the
	// paths above that require special privileges. These can't be
	// regular expressions, it is either exact match or prefix match.
	// For prefix match, append '*' as a suffix.
	Paths        []*Path
	PathsSpecial *logical.Paths

	// Secrets is the list of secret types that this backend can
	// return. It is used to automatically generate proper responses,
	// and ease specifying callbacks for revocation, renewal, etc.
	Secrets []*Secret

	// Rollback is called when a WAL entry (see wal.go) has to be rolled
	// back. It is called with the data from the entry.
	//
	// RollbackMinAge is the minimum age of a WAL entry before it is attempted
	// to be rolled back. This should be longer than the maximum time it takes
	// to successfully create a secret.
	Rollback       RollbackFunc
	RollbackMinAge time.Duration

	// AuthRenew is the callback to call when a RenewRequest for an
	// authentication comes in. By default, renewal won't be allowed.
	// See the built-in AuthRenew helpers in lease.go for common callbacks.
	AuthRenew OperationFunc

	logger  *log.Logger
	system  logical.SystemView
	once    sync.Once
	pathsRe []*regexp.Regexp
}

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(*logical.Request, *FieldData) (*logical.Response, error)

// RollbackFunc is the callback for rollbacks.
type RollbackFunc func(*logical.Request, string, interface{}) error

// logical.Backend impl.
func (b *Backend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	b.once.Do(b.init)

	// Check for special cased global operations. These don't route
	// to a specific Path.
	switch req.Operation {
	case logical.RenewOperation:
		fallthrough
	case logical.RevokeOperation:
		return b.handleRevokeRenew(req)
	case logical.RollbackOperation:
		return b.handleRollback(req)
	}

	// If the path is empty and it is a help operation, handle that.
	if req.Path == "" && req.Operation == logical.HelpOperation {
		return b.handleRootHelp()
	}

	// Find the matching route
	path, captures := b.route(req.Path)
	if path == nil {
		return nil, logical.ErrUnsupportedPath
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

	// Look up the callback for this operation
	var callback OperationFunc
	var ok bool
	if path.Callbacks != nil {
		callback, ok = path.Callbacks[req.Operation]
	}
	if !ok {
		if req.Operation == logical.HelpOperation {
			callback = path.helpCallback
			ok = true
		}
	}
	if !ok {
		return nil, logical.ErrUnsupportedOperation
	}

	fd := FieldData{
		Raw:    raw,
		Schema: path.Fields}

	if req.Operation != logical.HelpOperation {
		err := fd.Validate()
		if err != nil {
			return nil, err
		}
	}

	// Call the callback with the request and the data
	return callback(req, &fd)
}

// logical.Backend impl.
func (b *Backend) SpecialPaths() *logical.Paths {
	return b.PathsSpecial
}

// Setup is used to initialize the backend with the initial backend configuration
func (b *Backend) Setup(config *logical.BackendConfig) (logical.Backend, error) {
	b.logger = config.Logger
	b.system = config.System
	return b, nil
}

// Logger can be used to get the logger. If no logger has been set,
// the logs will be discarded.
func (b *Backend) Logger() *log.Logger {
	if b.logger != nil {
		return b.logger
	}

	return log.New(ioutil.Discard, "", 0)
}

func (b *Backend) System() logical.SystemView {
	return b.system
}

// Route looks up the path that would be used for a given path string.
func (b *Backend) Route(path string) *Path {
	result, _ := b.route(path)
	return result
}

// Secret is used to look up the secret with the given type.
func (b *Backend) Secret(k string) *Secret {
	for _, s := range b.Secrets {
		if s.Type == k {
			return s
		}
	}

	return nil
}

func (b *Backend) init() {
	b.pathsRe = make([]*regexp.Regexp, len(b.Paths))
	for i, p := range b.Paths {
		if len(p.Pattern) == 0 {
			panic(fmt.Sprintf("Routing pattern cannot be blank"))
		}
		// Automatically anchor the pattern
		if p.Pattern[0] != '^' {
			p.Pattern = "^" + p.Pattern
		}
		if p.Pattern[len(p.Pattern)-1] != '$' {
			p.Pattern = p.Pattern + "$"
		}
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

func (b *Backend) handleRootHelp() (*logical.Response, error) {
	// Build a mapping of the paths and get the paths alphabetized to
	// make the output prettier.
	pathsMap := make(map[string]*Path)
	paths := make([]string, 0, len(b.Paths))
	for i, p := range b.pathsRe {
		paths = append(paths, p.String())
		pathsMap[p.String()] = b.Paths[i]
	}
	sort.Strings(paths)

	// Build the path data
	pathData := make([]rootHelpTemplatePath, 0, len(paths))
	for _, route := range paths {
		p := pathsMap[route]
		pathData = append(pathData, rootHelpTemplatePath{
			Path: route,
			Help: strings.TrimSpace(p.HelpSynopsis),
		})
	}

	help, err := executeTemplate(rootHelpTemplate, &rootHelpTemplateData{
		Help:  strings.TrimSpace(b.Help),
		Paths: pathData,
	})
	if err != nil {
		return nil, err
	}

	return logical.HelpResponse(help, nil), nil
}

func (b *Backend) handleRevokeRenew(
	req *logical.Request) (*logical.Response, error) {
	// Special case renewal of authentication for credential backends
	if req.Operation == logical.RenewOperation && req.Auth != nil {
		return b.handleAuthRenew(req)
	}

	if req.Secret == nil {
		return nil, fmt.Errorf("request has no secret")
	}

	rawSecretType, ok := req.Secret.InternalData["secret_type"]
	if !ok {
		return nil, fmt.Errorf("secret is unsupported by this backend")
	}
	secretType, ok := rawSecretType.(string)
	if !ok {
		return nil, fmt.Errorf("secret is unsupported by this backend")
	}

	secret := b.Secret(secretType)
	if secret == nil {
		return nil, fmt.Errorf("secret is unsupported by this backend")
	}

	switch req.Operation {
	case logical.RenewOperation:
		return secret.HandleRenew(req)
	case logical.RevokeOperation:
		return secret.HandleRevoke(req)
	default:
		return nil, fmt.Errorf(
			"invalid operation for revoke/renew: %s", req.Operation)
	}
}

func (b *Backend) handleAuthRenew(req *logical.Request) (*logical.Response, error) {
	if b.AuthRenew == nil {
		return logical.ErrorResponse("this auth type doesn't support renew"), nil
	}

	return b.AuthRenew(req, nil)
}

func (b *Backend) handleRollback(
	req *logical.Request) (*logical.Response, error) {
	if b.Rollback == nil {
		return nil, logical.ErrUnsupportedOperation
	}

	var merr error
	keys, err := ListWAL(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if len(keys) == 0 {
		return nil, nil
	}

	// Calculate the minimum time that the WAL entries could be
	// created in order to be rolled back.
	age := b.RollbackMinAge
	if age == 0 {
		age = 10 * time.Minute
	}
	minAge := time.Now().UTC().Add(-1 * age)
	if _, ok := req.Data["immediate"]; ok {
		minAge = time.Now().UTC().Add(1000 * time.Hour)
	}

	for _, k := range keys {
		entry, err := GetWAL(req.Storage, k)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}
		if entry == nil {
			continue
		}

		// If the entry isn't old enough, then don't roll it back
		if !time.Unix(entry.CreatedAt, 0).Before(minAge) {
			continue
		}

		// Attempt a rollback
		err = b.Rollback(req, entry.Kind, entry.Data)
		if err != nil {
			err = fmt.Errorf(
				"Error rolling back '%s' entry: %s", entry.Kind, err)
		}
		if err == nil {
			err = DeleteWAL(req.Storage, k)
		}
		if err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	if merr == nil {
		return nil, nil
	}

	return logical.ErrorResponse(merr.Error()), nil
}

// FieldSchema is a basic schema to describe the format of a path field.
type FieldSchema struct {
	Type        FieldType
	Default     interface{}
	Description string
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
	case TypeMap:
		return map[string]interface{}{}
	case TypeDurationSecond:
		return 0
	default:
		panic("unknown type: " + t.String())
	}
}

type rootHelpTemplateData struct {
	Help  string
	Paths []rootHelpTemplatePath
}

type rootHelpTemplatePath struct {
	Path string
	Help string
}

const rootHelpTemplate = `
## DESCRIPTION

{{.Help}}

## PATHS

The following paths are supported by this backend. To view help for
any of the paths below, use the help command with any route matching
the path pattern. Note that depending on the policy of your auth token,
you may or may not be able to access certain paths.

{{range .Paths}}{{indent 4 .Path}}
{{indent 8 .Help}}

{{end}}

`
