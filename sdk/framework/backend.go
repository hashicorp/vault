// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-kms-wrapping/entropy/v2"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
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
	// PathsSpecial is the list of path patterns that denote the paths above
	// that require special privileges.
	Paths        []*Path
	PathsSpecial *logical.Paths

	// Secrets is the list of secret types that this backend can
	// return. It is used to automatically generate proper responses,
	// and ease specifying callbacks for revocation, renewal, etc.
	Secrets []*Secret

	// InitializeFunc is the callback, which if set, will be invoked via
	// Initialize() just after a plugin has been mounted.
	//
	// Note that storage writes should only occur on the active instance within a
	// primary cluster or local mount on a performance secondary. If your InitializeFunc
	// writes to storage, you can use the backend's WriteSafeReplicationState() method
	// to prevent it from attempting to write on a Vault instance with read-only storage.
	InitializeFunc InitializeFunc

	// PeriodicFunc is the callback, which if set, will be invoked when the
	// periodic timer of RollbackManager ticks. This can be used by
	// backends to do anything it wishes to do periodically.
	//
	// PeriodicFunc can be invoked to, say periodically delete stale
	// entries in backend's storage, while the backend is still being used.
	// (Note the difference between this action and `Clean`, which is
	// invoked just before the backend is unmounted).
	//
	// Note that storage writes should only occur on the active instance within a
	// primary cluster or local mount on a performance secondary. If your PeriodicFunc
	// writes to storage, you can use the backend's WriteSafeReplicationState() method
	// to prevent it from attempting to write on a Vault instance with read-only storage.
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

	// Invalidate is called when a key is modified, if required.
	Invalidate InvalidateFunc

	// AuthRenew is the callback to call when a RenewRequest for an
	// authentication comes in. By default, renewal won't be allowed.
	// See the built-in AuthRenew helpers in lease.go for common callbacks.
	AuthRenew OperationFunc

	// BackendType is the logical.BackendType for the backend implementation
	BackendType logical.BackendType

	// RunningVersion is the optional version that will be self-reported
	RunningVersion string

	logger  log.Logger
	system  logical.SystemView
	events  logical.EventSender
	once    sync.Once
	pathsRe []*regexp.Regexp
}

// periodicFunc is the callback called when the RollbackManager's timer ticks.
// This can be utilized by the backends to do anything it wants.
type periodicFunc func(context.Context, *logical.Request) error

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(context.Context, *logical.Request, *FieldData) (*logical.Response, error)

// ExistenceFunc is the callback called for an existence check on a path.
type ExistenceFunc func(context.Context, *logical.Request, *FieldData) (bool, error)

// WALRollbackFunc is the callback for rollbacks.
type WALRollbackFunc func(context.Context, *logical.Request, string, interface{}) error

// CleanupFunc is the callback for backend unload.
type CleanupFunc func(context.Context)

// InvalidateFunc is the callback for backend key invalidation.
type InvalidateFunc func(context.Context, string)

// InitializeFunc is the callback, which if set, will be invoked via
// Initialize() just after a plugin has been mounted.
type InitializeFunc func(context.Context, *logical.InitializationRequest) error

// PatchPreprocessorFunc is used by HandlePatchOperation in order to shape
// the input as defined by request handler prior to JSON marshaling
type PatchPreprocessorFunc func(map[string]interface{}) (map[string]interface{}, error)

// ErrNoEvents is returned when attempting to send an event, but when the event
// sender was not passed in during `backend.Setup()`.
var ErrNoEvents = errors.New("no event sender configured")

// Initialize is the logical.Backend implementation.
func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	if b.InitializeFunc != nil {
		return b.InitializeFunc(ctx, req)
	}
	return nil
}

// HandleExistenceCheck is the logical.Backend implementation.
func (b *Backend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (checkFound bool, exists bool, err error) {
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
		Schema: path.Fields,
	}

	err = fd.Validate()
	if err != nil {
		return false, false, errutil.UserError{Err: err.Error()}
	}

	// Call the callback with the request and the data
	exists, err = path.ExistenceCheck(ctx, req, &fd)
	return
}

// HandleRequest is the logical.Backend implementation.
func (b *Backend) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	b.once.Do(b.init)

	// Check for special cased global operations. These don't route
	// to a specific Path.
	switch req.Operation {
	case logical.RenewOperation:
		fallthrough
	case logical.RevokeOperation:
		return b.handleRevokeRenew(ctx, req)
	case logical.RollbackOperation:
		return b.handleRollback(ctx, req)
	}

	// If the path is empty and it is a help operation, handle that.
	if req.Path == "" && req.Operation == logical.HelpOperation {
		return b.handleRootHelp(req)
	}

	// Find the matching route
	path, captures := b.route(req.Path)
	if path == nil {
		return nil, logical.ErrUnsupportedPath
	}

	// Check if a feature is required and if the license has that feature
	if path.FeatureRequired != license.FeatureNone {
		hasFeature := b.system.HasFeature(path.FeatureRequired)
		if !hasFeature {
			return nil, logical.CodedError(401, "Feature Not Enabled")
		}
	}

	// Build up the data for the route, with the URL taking priority
	// for the fields over the PUT data.
	raw := make(map[string]interface{}, len(path.Fields))
	var ignored []string
	for k, v := range req.Data {
		raw[k] = v
		if !path.TakesArbitraryInput && path.Fields[k] == nil {
			ignored = append(ignored, k)
		}
	}

	var replaced []string
	for k, v := range captures {
		if raw[k] != nil {
			replaced = append(replaced, k)
		}
		raw[k] = v
	}

	// Look up the callback for this operation, preferring the
	// path.Operations definition if present.
	var callback OperationFunc

	if path.Operations != nil {
		if op, ok := path.Operations[req.Operation]; ok {

			// Check whether this operation should be forwarded
			if sysView := b.System(); sysView != nil {
				replState := sysView.ReplicationState()
				props := op.Properties()

				if props.ForwardPerformanceStandby && replState.HasState(consts.ReplicationPerformanceStandby) {
					return nil, logical.ErrReadOnly
				}

				if props.ForwardPerformanceSecondary && !sysView.LocalMount() && replState.HasState(consts.ReplicationPerformanceSecondary) {
					return nil, logical.ErrReadOnly
				}
			}

			callback = op.Handler()
		}
	} else {
		callback = path.Callbacks[req.Operation]
	}
	ok := callback != nil

	if !ok {
		if req.Operation == logical.HelpOperation {
			callback = path.helpCallback(b)
			ok = true
		}
	}
	if !ok {
		return nil, logical.ErrUnsupportedOperation
	}

	fd := FieldData{
		Raw:    raw,
		Schema: path.Fields,
	}

	if req.Operation != logical.HelpOperation {
		err := fd.Validate()
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Field validation failed: %s", err.Error())), nil
		}
	}

	resp, err := callback(ctx, req, &fd)
	if err != nil {
		return resp, err
	}

	switch resp {
	case nil:
	default:
		// If fields supplied in the request are not present in the field schema
		// of the path, add a warning to the response indicating that those
		// parameters will be ignored.
		sort.Strings(ignored)

		if len(ignored) != 0 {
			resp.AddWarning(fmt.Sprintf("Endpoint ignored these unrecognized parameters: %v", ignored))
		}
		// If fields supplied in the request is being overwritten by the values
		// supplied in the API request path, add a warning to the response
		// indicating that those parameters will be replaced.
		if len(replaced) != 0 {
			resp.AddWarning(fmt.Sprintf("Endpoint replaced the value of these parameters with the values captured from the endpoint's path: %v", replaced))
		}
	}

	return resp, nil
}

// HandlePatchOperation acts as an abstraction for performing JSON merge patch
// operations (see https://datatracker.ietf.org/doc/html/rfc7396) for HTTP
// PATCH requests. It is responsible for properly processing and marshalling
// the input and existing resource prior to performing the JSON merge operation
// using the MergePatch function from the json-patch library. The preprocessor
// is an arbitrary func that can be provided to further process the input. The
// MergePatch function accepts and returns byte arrays. Null values will unset
// fields defined within the input's FieldData (as if they were never specified)
// and remove user-specified keys that exist within a map field.
func HandlePatchOperation(input *FieldData, resource map[string]interface{}, preprocessor PatchPreprocessorFunc) ([]byte, error) {
	var err error

	if resource == nil {
		return nil, fmt.Errorf("resource does not exist")
	}

	inputMap := map[string]interface{}{}

	for key := range input.Raw {
		if _, ok := input.Schema[key]; !ok {
			// Only accept fields in the schema
			continue
		}

		// Ensure data types are handled properly according to the FieldSchema
		val, ok, err := input.GetOkErr(key)
		if err != nil {
			return nil, err
		}

		if ok {
			inputMap[key] = val
		}
	}

	if preprocessor != nil {
		inputMap, err = preprocessor(inputMap)
		if err != nil {
			return nil, err
		}
	}

	marshaledResource, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	marshaledInput, err := json.Marshal(inputMap)
	if err != nil {
		return nil, err
	}

	modified, err := jsonpatch.MergePatch(marshaledResource, marshaledInput)
	if err != nil {
		return nil, err
	}

	return modified, nil
}

// SpecialPaths is the logical.Backend implementation.
func (b *Backend) SpecialPaths() *logical.Paths {
	return b.PathsSpecial
}

// Cleanup is used to release resources and prepare to stop the backend
func (b *Backend) Cleanup(ctx context.Context) {
	if b.Clean != nil {
		b.Clean(ctx)
	}
}

// InvalidateKey is used to clear caches and reset internal state on key changes
func (b *Backend) InvalidateKey(ctx context.Context, key string) {
	if b.Invalidate != nil {
		b.Invalidate(ctx, key)
	}
}

// Setup is used to initialize the backend with the initial backend configuration
func (b *Backend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	b.logger = config.Logger
	b.system = config.System
	b.events = config.EventsSender
	return nil
}

// GetRandomReader returns an io.Reader to use for generating key material in
// backends. If the backend has access to an external entropy source it will
// return that, otherwise it returns crypto/rand.Reader.
func (b *Backend) GetRandomReader() io.Reader {
	if sourcer, ok := b.System().(entropy.Sourcer); ok {
		return entropy.NewReader(sourcer)
	}

	return rand.Reader
}

// Logger can be used to get the logger. If no logger has been set,
// the logs will be discarded.
func (b *Backend) Logger() log.Logger {
	if b.logger != nil {
		return b.logger
	}

	return logging.NewVaultLoggerWithWriter(ioutil.Discard, log.NoLevel)
}

// System returns the backend's system view.
func (b *Backend) System() logical.SystemView {
	return b.system
}

// Type returns the backend type
func (b *Backend) Type() logical.BackendType {
	return b.BackendType
}

// Version returns the plugin version information
func (b *Backend) PluginVersion() logical.PluginVersion {
	return logical.PluginVersion{
		Version: b.RunningVersion,
	}
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

// WriteSafeReplicationState returns true if this backend instance is capable of writing
// to storage without receiving an ErrReadOnly error. The active instance in a primary
// cluster or a local mount on a performance secondary is capable of writing to storage.
func (b *Backend) WriteSafeReplicationState() bool {
	replicationState := b.System().ReplicationState()
	return (b.System().LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby)
}

// init runs as a sync.Once function from any plugin entry point which needs to route requests by paths.
// It may panic if a coding error in the plugin is detected.
// For builtin plugins, this is unit tested in helper/builtinplugins/builtinplugins_test.go.
// For other plugins, any unit test that attempts to perform any request to the plugin will exercise these checks.
func (b *Backend) init() {
	b.pathsRe = make([]*regexp.Regexp, len(b.Paths))
	for i, p := range b.Paths {
		// Detect the coding error of failing to initialise Pattern
		if len(p.Pattern) == 0 {
			panic(fmt.Sprintf("Routing pattern cannot be blank"))
		}

		// Detect the coding error of attempting to define a CreateOperation without defining an ExistenceCheck
		if p.ExistenceCheck == nil {
			if _, ok := p.Operations[logical.CreateOperation]; ok {
				panic(fmt.Sprintf("Pattern %v defines a CreateOperation but no ExistenceCheck", p.Pattern))
			}
			if _, ok := p.Callbacks[logical.CreateOperation]; ok {
				panic(fmt.Sprintf("Pattern %v defines a CreateOperation but no ExistenceCheck", p.Pattern))
			}
		}

		// Automatically anchor the pattern
		if p.Pattern[0] != '^' {
			p.Pattern = "^" + p.Pattern
		}
		if p.Pattern[len(p.Pattern)-1] != '$' {
			p.Pattern = p.Pattern + "$"
		}

		// Detect the coding error of an invalid Pattern
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

func (b *Backend) handleRootHelp(req *logical.Request) (*logical.Response, error) {
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

	// Plugins currently don't have a direct knowledge of their own "type"
	// (e.g. "kv", "cubbyhole"). It defaults to the name of the executable but
	// can be overridden when the plugin is mounted. Since we need this type to
	// form the request & response full names, we are passing it as an optional
	// request parameter to the plugin's root help endpoint. If specified in
	// the request, the type will be used as part of the request/response body
	// names in the OAS document.
	requestResponsePrefix := req.GetString("requestResponsePrefix")

	// Build OpenAPI response for the entire backend
	vaultVersion := "unknown"
	if b.System() != nil {
		env, err := b.System().PluginEnv(context.Background())
		if err != nil {
			return nil, err
		}
		vaultVersion = env.VaultVersion
	}

	doc := NewOASDocument(vaultVersion)
	if err := documentPaths(b, requestResponsePrefix, doc); err != nil {
		b.Logger().Warn("error generating OpenAPI", "error", err)
	}

	return logical.HelpResponse(help, nil, doc), nil
}

func (b *Backend) handleRevokeRenew(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// Special case renewal of authentication for credential backends
	if req.Operation == logical.RenewOperation && req.Auth != nil {
		return b.handleAuthRenew(ctx, req)
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
		return secret.HandleRenew(ctx, req)
	case logical.RevokeOperation:
		return secret.HandleRevoke(ctx, req)
	default:
		return nil, fmt.Errorf("invalid operation for revoke/renew: %q", req.Operation)
	}
}

// handleRollback invokes the PeriodicFunc set on the backend. It also does a
// WAL rollback operation.
func (b *Backend) handleRollback(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// Response is not expected from the periodic operation.
	var resp *logical.Response

	merr := new(multierror.Error)
	if b.PeriodicFunc != nil {
		if err := b.PeriodicFunc(ctx, req); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	if b.WALRollback != nil {
		var err error
		resp, err = b.handleWALRollback(ctx, req)
		if err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return resp, merr.ErrorOrNil()
}

func (b *Backend) handleAuthRenew(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.AuthRenew == nil {
		return logical.ErrorResponse("this auth type doesn't support renew"), nil
	}

	return b.AuthRenew(ctx, req, nil)
}

func (b *Backend) handleWALRollback(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.WALRollback == nil {
		return nil, logical.ErrUnsupportedOperation
	}

	var merr error
	keys, err := ListWAL(ctx, req.Storage)
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
		entry, err := GetWAL(ctx, req.Storage, k)
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
		err = b.WALRollback(ctx, req, entry.Kind, entry.Data)
		if err != nil {
			err = errwrap.Wrapf(fmt.Sprintf("error rolling back %q entry: {{err}}", entry.Kind), err)
		}
		if err == nil {
			err = DeleteWAL(ctx, req.Storage, k)
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

// SendEvent is used to send events through the underlying EventSender.
// It returns ErrNoEvents if the events system has not been configured or enabled.
func (b *Backend) SendEvent(ctx context.Context, eventType logical.EventType, event *logical.EventData) error {
	if b.events == nil {
		return ErrNoEvents
	}
	return b.events.SendEvent(ctx, eventType, event)
}

// FieldSchema is a basic schema to describe the format of a path field.
type FieldSchema struct {
	Type        FieldType
	Default     interface{}
	Description string

	// The Required and Deprecated members are only used by openapi, and are not actually
	// used by the framework.
	Required   bool
	Deprecated bool

	// Query indicates this field will be expected as a query parameter as part
	// of ReadOperation, ListOperation or DeleteOperation requests:
	//
	//   /v1/foo/bar?some_param=some_value
	//
	// The field will still be expected as a request body parameter for
	// CreateOperation or UpdateOperation requests!
	//
	// To put that another way, you should set Query for any non-path parameter
	// you want to use in a read/list/delete operation.  While setting the Query
	// field to `true` is not required in such cases (Vault will expose the
	// query parameters to you via req.Data regardless), it is highly
	// recommended to do so in order to improve the quality of the generated
	// OpenAPI documentation (as well as any code generation based on it), which
	// will otherwise incorrectly omit the parameter.
	//
	// The reason for this design is historical: back at the start of 2018,
	// query parameters were not mapped to fields at all, and it was implicit
	// that all non-path fields were exclusively for the use of create/update
	// operations.  Since then, support for query parameters has gradually been
	// extended to read, delete and list operations - and now this declarative
	// metadata is needed, so that the OpenAPI generator can know which
	// parameters are actually referred to, from within the code of
	// read/delete/list operation handler functions.
	Query bool

	// AllowedValues is an optional list of permitted values for this field.
	// This constraint is not (yet) enforced by the framework, but the list is
	// output as part of OpenAPI generation and may affect documentation and
	// dynamic UI generation.
	AllowedValues []interface{}

	// DisplayAttrs provides hints for UI and documentation generators. They
	// will be included in OpenAPI output if set.
	DisplayAttrs *DisplayAttributes
}

// DefaultOrZero returns the default value if it is set, or otherwise
// the zero value of the type.
func (s *FieldSchema) DefaultOrZero() interface{} {
	if s.Default != nil {
		switch s.Type {
		case TypeDurationSecond, TypeSignedDurationSecond:
			resultDur, err := parseutil.ParseDurationSecond(s.Default)
			if err != nil {
				return s.Type.Zero()
			}
			return int(resultDur.Seconds())

		default:
			return s.Default
		}
	}

	return s.Type.Zero()
}

// Zero returns the correct zero-value for a specific FieldType
func (t FieldType) Zero() interface{} {
	switch t {
	case TypeString, TypeNameString, TypeLowerCaseString:
		return ""
	case TypeInt:
		return 0
	case TypeInt64:
		return int64(0)
	case TypeBool:
		return false
	case TypeMap:
		return map[string]interface{}{}
	case TypeKVPairs:
		return map[string]string{}
	case TypeDurationSecond, TypeSignedDurationSecond:
		return 0
	case TypeSlice:
		return []interface{}{}
	case TypeStringSlice, TypeCommaStringSlice:
		return []string{}
	case TypeCommaIntSlice:
		return []int{}
	case TypeHeader:
		return http.Header{}
	case TypeFloat:
		return 0.0
	case TypeTime:
		return time.Time{}
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
