package framework

import (
	"fmt"
	"regexp"
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
	// Paths are the various routes that the backend responds to.
	// This cannot be modified after construction (i.e. dynamically changing
	// paths, including adding or removing, is not allowed once the
	// backend is in use).
	//
	// PathsRoot is the list of path patterns that denote the
	// paths above that require root-level privileges. These can't be
	// regular expressions, it is either exact match or prefix match.
	// For prefix match, append '*' as a suffix.
	Paths     []*Path
	PathsRoot []string

	// Secrets is the list of secret types that this backend can
	// return. It is used to automatically generate proper responses,
	// and ease specifying callbacks for revocation, renewal, etc.
	Secrets []*Secret

	// Rollback is called when a WAL entry (see wal.go) has to be rolled
	// back. It is called with the data from the entry. Boolean true should
	// be returned on success. Errors should just be logged.
	Rollback       func(kind string, data interface{}) bool
	RollbackMinAge time.Duration

	once    sync.Once
	pathsRe []*regexp.Regexp
}

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(*Request) (*logical.Response, error)

// logical.Backend impl.
func (b *Backend) HandleRequest(req *logical.Request) (*logical.Response, error) {
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
		if req.Operation == logical.HelpOperation && path.HelpSynopsis != "" {
			callback = path.helpCallback
			ok = true
		}
	}
	if !ok {
		return nil, logical.ErrUnsupportedOperation
	}

	// Call the callback with the request and the data
	return callback(&Request{
		Backend:        b,
		LogicalRequest: req,
		Data: &FieldData{
			Raw:    raw,
			Schema: path.Fields,
		},
	})
}

// logical.Backend impl.
func (b *Backend) RootPaths() []string {
	return b.PathsRoot
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

func (b *Backend) handleRevokeRenew(
	req *logical.Request) (*logical.Response, error) {
	leaseRaw, ok := req.Data["previous_lease"]
	if !ok {
		return nil, fmt.Errorf("no previous lease for revoke")
	}

	lease, ok := leaseRaw.(*logical.Lease)
	if !ok {
		return nil, fmt.Errorf("previous_lease is now *logical.Lease")
	}

	secretType, secretId := SecretType(lease.VaultID)
	if secretType == "" {
		return nil, fmt.Errorf("secret is unsupported by this backend")
	}

	secret := b.Secret(secretType)
	if secret == nil {
		return nil, fmt.Errorf("secret is unsupported by this backend")
	}

	var fn OperationFunc
	switch req.Operation {
	case logical.RenewOperation:
		fn = secret.Renew
	case logical.RevokeOperation:
		fn = secret.Revoke
	default:
		return nil, fmt.Errorf(
			"invalid operation for revoke/renew: %s", req.Operation)
	}

	if fn == nil {
		return nil, logical.ErrUnsupportedOperation
	}

	var data map[string]interface{}
	if raw, ok := req.Data["previous_data"]; ok {
		data = raw.(map[string]interface{})
	}

	return fn(&Request{
		Backend:        b,
		LogicalRequest: req,

		Data: &FieldData{
			Raw:    data,
			Schema: secret.Fields,
		},

		SecretType: secretType,
		SecretId:   secretId,
	})
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
		if b.Rollback(entry.Kind, entry.Data) {
			if err := DeleteWAL(req.Storage, k); err != nil {
				merr = multierror.Append(merr, err)
			}
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
	default:
		panic("unknown type: " + t.String())
	}
}
