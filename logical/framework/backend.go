package framework

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/parseutil"
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

	// PeriodicFunc is the callback, which if set, will be invoked when the
	// periodic timer of RollbackManager ticks. This can be used by
	// backends to do anything it wishes to do periodically.
	//
	// PeriodicFunc can be invoked to, say to periodically delete stale
	// entries in backend's storage, while the backend is still being used.
	// (Note the different of this action from what `Clean` does, which is
	// invoked just before the backend is unmounted).
	PeriodicFunc periodicFunc

	// WALRollback is called when a WAL entry (see wal.go) has to be rolled
	// back. It is called with the data from the entry.
	//
	// WALRollbackMinAge is the minimum age of a WAL entry before it is attempted
	// to be rolled back. This should be longer than the maximum time it takes
	// to successfully create a secret.
	WALRollback       WALRollbackFunc
	WALRollbackMinAge time.Duration

	// Clean is called on unload to clean up e.g any existing connections
	// to the backend, if required.
	Clean CleanupFunc

	// Initialize is called after a backend is created. Storage should not be
	// written to before this function is called.
	Init InitializeFunc

	// Invalidate is called when a keys is modified if required
	Invalidate InvalidateFunc

	// AuthRenew is the callback to call when a RenewRequest for an
	// authentication comes in. By default, renewal won't be allowed.
	// See the built-in AuthRenew helpers in lease.go for common callbacks.
	AuthRenew OperationFunc

	// LicenseRegistration is called to register the license for a backend.
	LicenseRegistration LicenseRegistrationFunc

	// Type is the logical.BackendType for the backend implementation
	BackendType logical.BackendType

	logger  log.Logger
	system  logical.SystemView
	once    sync.Once
	pathsRe []*regexp.Regexp
}

// periodicFunc is the callback called when the RollbackManager's timer ticks.
// This can be utilized by the backends to do anything it wants.
type periodicFunc func(*logical.Request) error

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(*logical.Request, *FieldData) (*logical.Response, error)

// WALRollbackFunc is the callback for rollbacks.
type WALRollbackFunc func(*logical.Request, string, interface{}) error

// CleanupFunc is the callback for backend unload.
type CleanupFunc func()

// InitializeFunc is the callback for backend creation.
type InitializeFunc func() error

// InvalidateFunc is the callback for backend key invalidation.
type InvalidateFunc func(string)

// LicenseRegistrationFunc is the callback for backend license registration.
type LicenseRegistrationFunc func(interface{}) error

// HandleExistenceCheck is the logical.Backend implementation.
func (b *Backend) HandleExistenceCheck(req *logical.Request) (checkFound bool, exists bool, err error) {
	b.once.Do(b.init)

	// Ensure we are only doing this when one of the correct operations is in play
	switch req.Operation {
	case logical.CreateOperation:
	case logical.UpdateOperation:
	default:
		return false, false, fmt.Errorf("incorrect operation type %v for an existence check", req.Operation)
	}

	// Find the matching route
	path, captures := b.route(req.Path)
	if path == nil {
		return false, false, logical.ErrUnsupportedPath
	}

	if path.ExistenceCheck == nil {
		return false, false, nil
	}

	checkFound = true

	// Build up the data for the route, with the URL taking priority
	// for the fields over the PUT data.
	raw := make(map[string]interface{}, len(path.Fields))
	for k, v := range req.Data {
		raw[k] = v
	}
	for k, v := range captures {
		raw[k] = v
	}

	fd := FieldData{
		Raw:    raw,
		Schema: path.Fields}

	err = fd.Validate()
	if err != nil {
		return false, false, errutil.UserError{Err: err.Error()}
	}

	// Call the callback with the request and the data
	exists, err = path.ExistenceCheck(req, &fd)
	return
}

// HandleRequest is the logical.Backend implementation.
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

// SpecialPaths is the logical.Backend implementation.
func (b *Backend) SpecialPaths() *logical.Paths {
	return b.PathsSpecial
}

// Cleanup is used to release resources and prepare to stop the backend
func (b *Backend) Cleanup() {
	if b.Clean != nil {
		b.Clean()
	}
}

// Initialize calls the backend's Init func if set.
func (b *Backend) Initialize() error {
	if b.Init != nil {
		return b.Init()
	}

	return nil
}

// InvalidateKey is used to clear caches and reset internal state on key changes
func (b *Backend) InvalidateKey(key string) {
	if b.Invalidate != nil {
		b.Invalidate(key)
	}
}

// Setup is used to initialize the backend with the initial backend configuration
func (b *Backend) Setup(config *logical.BackendConfig) error {
	b.logger = config.Logger
	b.system = config.System
	return nil
}

// Logger can be used to get the logger. If no logger has been set,
// the logs will be discarded.
func (b *Backend) Logger() log.Logger {
	if b.logger != nil {
		return b.logger
	}

	return logformat.NewVaultLoggerWithWriter(ioutil.Discard, log.LevelOff)
}

// System returns the backend's system view.
func (b *Backend) System() logical.SystemView {
	return b.system
}

// Type returns the backend type
func (b *Backend) Type() logical.BackendType {
	return b.BackendType
}

// RegisterLicense performs backend license registration.
func (b *Backend) RegisterLicense(license interface{}) error {
	if b.LicenseRegistration == nil {
		return nil
	}
	return b.LicenseRegistration(license)
}

// SanitizeTTLStr takes in the TTL and MaxTTL values provided by the user,
// compares those with the SystemView values. If they are empty a value of 0 is
// set, which will cause initial secret or LeaseExtend operations to use the
// mount/system defaults.  If they are set, their boundaries are validated.
func (b *Backend) SanitizeTTLStr(ttlStr, maxTTLStr string) (ttl, maxTTL time.Duration, err error) {
	if len(ttlStr) == 0 || ttlStr == "0" {
		ttl = 0
	} else {
		ttl, err = time.ParseDuration(ttlStr)
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid ttl: %s", err)
		}
	}

	if len(maxTTLStr) == 0 || maxTTLStr == "0" {
		maxTTL = 0
	} else {
		maxTTL, err = time.ParseDuration(maxTTLStr)
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid max_ttl: %s", err)
		}
	}

	ttl, maxTTL, err = b.SanitizeTTL(ttl, maxTTL)

	return
}

// SanitizeTTL caps the boundaries of ttl and max_ttl values to the
// backend mount's max_ttl value.
func (b *Backend) SanitizeTTL(ttl, maxTTL time.Duration) (time.Duration, time.Duration, error) {
	sysMaxTTL := b.System().MaxLeaseTTL()
	if ttl > sysMaxTTL {
		return 0, 0, fmt.Errorf("\"ttl\" value must be less than allowed max lease TTL value '%s'", sysMaxTTL.String())
	}
	if maxTTL > sysMaxTTL {
		return 0, 0, fmt.Errorf("\"max_ttl\" value must be less than allowed max lease TTL value '%s'", sysMaxTTL.String())
	}
	if ttl > maxTTL && maxTTL != 0 {
		ttl = maxTTL
	}
	return ttl, maxTTL, nil
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

// handleRollback invokes the PeriodicFunc set on the backend. It also does a WAL rollback operation.
func (b *Backend) handleRollback(
	req *logical.Request) (*logical.Response, error) {
	// Response is not expected from the periodic operation.
	if b.PeriodicFunc != nil {
		if err := b.PeriodicFunc(req); err != nil {
			return nil, err
		}
	}

	return b.handleWALRollback(req)
}

func (b *Backend) handleAuthRenew(req *logical.Request) (*logical.Response, error) {
	if b.AuthRenew == nil {
		return logical.ErrorResponse("this auth type doesn't support renew"), nil
	}

	return b.AuthRenew(req, nil)
}

func (b *Backend) handleWALRollback(
	req *logical.Request) (*logical.Response, error) {
	if b.WALRollback == nil {
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
	age := b.WALRollbackMinAge
	if age == 0 {
		age = 10 * time.Minute
	}
	minAge := time.Now().Add(-1 * age)
	if _, ok := req.Data["immediate"]; ok {
		minAge = time.Now().Add(1000 * time.Hour)
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

		// Attempt a WAL rollback
		err = b.WALRollback(req, entry.Kind, entry.Data)
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
		switch s.Type {
		case TypeDurationSecond:
			var result int
			switch inp := s.Default.(type) {
			case nil:
				return s.Type.Zero()
			case int:
				result = inp
			case int64:
				result = int(inp)
			case float32:
				result = int(inp)
			case float64:
				result = int(inp)
			case string:
				dur, err := parseutil.ParseDurationSecond(inp)
				if err != nil {
					return s.Type.Zero()
				}
				result = int(dur.Seconds())
			case json.Number:
				valInt64, err := inp.Int64()
				if err != nil {
					return s.Type.Zero()
				}
				result = int(valInt64)
			default:
				return s.Type.Zero()
			}
			return result

		default:
			return s.Default
		}
	}

	return s.Type.Zero()
}

// Zero returns the correct zero-value for a specific FieldType
func (t FieldType) Zero() interface{} {
	switch t {
	case TypeNameString:
		return ""
	case TypeString:
		return ""
	case TypeInt:
		return 0
	case TypeBool:
		return false
	case TypeMap:
		return map[string]interface{}{}
	case TypeKVPairs:
		return map[string]string{}
	case TypeDurationSecond:
		return 0
	case TypeSlice:
		return []interface{}{}
	case TypeStringSlice, TypeCommaStringSlice:
		return []string{}
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
