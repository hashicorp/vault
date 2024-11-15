// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	// ErrVariableNotFound was used as the content of an error string.
	//
	// Deprecated: use ErrVariablePathNotFound instead.
	ErrVariableNotFound = "variable not found"
)

var (
	// ErrVariablePathNotFound is returned when trying to read a variable that
	// does not exist.
	ErrVariablePathNotFound = errors.New("variable not found")
)

// Variables is used to access variables.
type Variables struct {
	client *Client
}

// Variables returns a new handle on the variables.
func (c *Client) Variables() *Variables {
	return &Variables{client: c}
}

// Create is used to create a variable.
func (vars *Variables) Create(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out Variable
	wm, err := vars.client.put("/v1/var/"+v.Path, v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// CheckedCreate is used to create a variable if it doesn't exist
// already. If it does, it will return a ErrCASConflict that can be unwrapped
// for more details.
func (vars *Variables) CheckedCreate(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out Variable
	wm, err := vars.writeChecked("/v1/var/"+v.Path+"?cas=0", v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// Read is used to query a single variable by path. This will error
// if the variable is not found.
func (vars *Variables) Read(path string, qo *QueryOptions) (*Variable, *QueryMeta, error) {
	path = cleanPathString(path)
	var v = new(Variable)
	qm, err := vars.readInternal("/v1/var/"+path, &v, qo)
	if err != nil {
		return nil, nil, err
	}
	if v == nil {
		return nil, qm, ErrVariablePathNotFound
	}
	return v, qm, nil
}

// Peek is used to query a single variable by path, but does not error
// when the variable is not found
func (vars *Variables) Peek(path string, qo *QueryOptions) (*Variable, *QueryMeta, error) {
	path = cleanPathString(path)
	var v = new(Variable)
	qm, err := vars.readInternal("/v1/var/"+path, &v, qo)
	if err != nil {
		return nil, nil, err
	}
	return v, qm, nil
}

// Update is used to update a variable.
func (vars *Variables) Update(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out Variable

	wm, err := vars.client.put("/v1/var/"+v.Path, v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// CheckedUpdate is used to updated a variable if the modify index
// matches the one on the server.  If it does not, it will return an
// ErrCASConflict that can be unwrapped for more details.
func (vars *Variables) CheckedUpdate(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out Variable

	wm, err := vars.writeChecked("/v1/var/"+v.Path+"?cas="+fmt.Sprint(v.ModifyIndex), v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// Delete is used to delete a variable
func (vars *Variables) Delete(path string, qo *WriteOptions) (*WriteMeta, error) {
	path = cleanPathString(path)
	wm, err := vars.deleteInternal(path, qo)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// CheckedDelete is used to conditionally delete a variable. If the
// existing variable does not match the provided checkIndex, it will return an
// ErrCASConflict that can be unwrapped for more details.
func (vars *Variables) CheckedDelete(path string, checkIndex uint64, qo *WriteOptions) (*WriteMeta, error) {
	path = cleanPathString(path)
	wm, err := vars.deleteChecked(path, checkIndex, qo)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// List is used to dump all of the variables, can be used to pass prefix
// via QueryOptions rather than as a parameter
func (vars *Variables) List(qo *QueryOptions) ([]*VariableMetadata, *QueryMeta, error) {
	var resp []*VariableMetadata
	qm, err := vars.client.query("/v1/vars", &resp, qo)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// PrefixList is used to do a prefix List search over variables.
func (vars *Variables) PrefixList(prefix string, qo *QueryOptions) ([]*VariableMetadata, *QueryMeta, error) {
	if qo == nil {
		qo = &QueryOptions{Prefix: prefix}
	} else {
		qo.Prefix = prefix
	}
	return vars.List(qo)
}

// GetItems returns the inner Items collection from a variable at a given path.
//
// Deprecated: Use GetVariableItems instead.
func (vars *Variables) GetItems(path string, qo *QueryOptions) (*VariableItems, *QueryMeta, error) {
	vi, qm, err := vars.GetVariableItems(path, qo)
	if err != nil {
		return nil, nil, err
	}
	return &vi, qm, nil
}

// GetVariableItems returns the inner Items collection from a variable at a given path.
func (vars *Variables) GetVariableItems(path string, qo *QueryOptions) (VariableItems, *QueryMeta, error) {
	path = cleanPathString(path)
	v := new(Variable)

	qm, err := vars.readInternal("/v1/var/"+path, &v, qo)
	if err != nil {
		return nil, nil, err
	}

	// note: readInternal will in fact turn our v into a nil if not found
	if v == nil {
		return nil, nil, ErrVariablePathNotFound
	}

	return v.Items, qm, nil
}

// RenewLock renews the lease for the lock on the given variable. It has to be called
// before the lock's TTL expires or the lock will be automatically released after the
// delay period.
func (vars *Variables) RenewLock(v *Variable, qo *WriteOptions) (*VariableMetadata, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out VariableMetadata

	wm, err := vars.client.put("/v1/var/"+v.Path+"?lock-renew", v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// ReleaseLock removes the lock on the given variable.
func (vars *Variables) ReleaseLock(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	return vars.lockOperation(v, qo, "lock-release")
}

// AcquireLock adds a lock on the given variable and starts a lease on it. In order
// to make any update on the locked variable, the lock ID has to be included in the
// request. In order to maintain ownership of the lock, the lease needs to be
// periodically renewed before the lock's TTL expires.
func (vars *Variables) AcquireLock(v *Variable, qo *WriteOptions) (*Variable, *WriteMeta, error) {
	return vars.lockOperation(v, qo, "lock-acquire")
}

func (vars *Variables) lockOperation(v *Variable, qo *WriteOptions, operation string) (*Variable, *WriteMeta, error) {
	v.Path = cleanPathString(v.Path)
	var out Variable

	wm, err := vars.client.put("/v1/var/"+v.Path+"?"+operation, v, &out, qo)
	if err != nil {
		return nil, wm, err
	}
	return &out, wm, nil
}

// readInternal exists because the API's higher-level read method requires
// the status code to be 200 (OK). For Peek(), we do not consider 403 (Permission
// Denied or 404 (Not Found) an error, this function just returns a nil in those
// cases.
func (vars *Variables) readInternal(endpoint string, out **Variable, q *QueryOptions) (*QueryMeta, error) {
	// todo(shoenig): seems like this could just return a *Variable instead of taking
	// in a **Variable and modifying it?

	r, err := vars.client.newRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)

	checkFn := requireStatusIn(http.StatusOK, http.StatusNotFound, http.StatusForbidden) //nolint:bodyclose
	rtt, resp, err := checkFn(vars.client.doRequest(r))                                  //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	qm := &QueryMeta{}
	_ = parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if resp.StatusCode == http.StatusNotFound {
		*out = nil
		_ = resp.Body.Close()
		return qm, nil
	}

	if resp.StatusCode == http.StatusForbidden {
		*out = nil
		_ = resp.Body.Close()
		// On a 403, there is no QueryMeta to parse, but consul-template--the
		// main consumer of the Peek() func that calls this method needs the
		// value to be non-zero; so set them to a reasonable but artificial
		// value. Index 1 doesn't say anything about the cluster, and there
		// has to be a KnownLeader to get a 403.
		qm.LastIndex = 1
		qm.KnownLeader = true
		return qm, nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	if err = decodeBody(resp, out); err != nil {
		return nil, err
	}

	return qm, nil
}

// deleteInternal exists because the API's higher-level delete method requires
// the status code to be 200 (OK). The SV HTTP API returns a 204 (No Content)
// on success.
func (vars *Variables) deleteInternal(path string, q *WriteOptions) (*WriteMeta, error) {
	r, err := vars.client.newRequest("DELETE", fmt.Sprintf("/v1/var/%s", path))
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)

	checkFn := requireStatusIn(http.StatusOK, http.StatusNoContent) //nolint:bodyclose
	rtt, resp, err := checkFn(vars.client.doRequest(r))             //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{RequestTime: rtt}
	_ = parseWriteMeta(resp, wm)
	return wm, nil
}

// deleteChecked exists because the API's higher-level delete method requires
// the status code to be OK. The SV HTTP API returns a 204 (No Content) on
// success and a 409 (Conflict) on a CAS error.
func (vars *Variables) deleteChecked(path string, checkIndex uint64, q *WriteOptions) (*WriteMeta, error) {
	r, err := vars.client.newRequest("DELETE", fmt.Sprintf("/v1/var/%s?cas=%v", path, checkIndex))
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	checkFn := requireStatusIn(http.StatusOK, http.StatusNoContent, http.StatusConflict) //nolint:bodyclose
	rtt, resp, err := checkFn(vars.client.doRequest(r))                                  //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{RequestTime: rtt}
	_ = parseWriteMeta(resp, wm)

	// The only reason we should decode the response body is if
	// it is a conflict response. Otherwise, there won't be one.
	if resp.StatusCode == http.StatusConflict {

		conflict := new(Variable)
		if err = decodeBody(resp, &conflict); err != nil {
			return nil, err
		}
		return nil, ErrCASConflict{
			Conflict:   conflict,
			CheckIndex: checkIndex,
		}
	}
	return wm, nil
}

// writeChecked exists because the API's higher-level write method requires
// the status code to be OK. The SV HTTP API returns a 200 (OK) on
// success and a 409 (Conflict) on a CAS error.
func (vars *Variables) writeChecked(endpoint string, in *Variable, out *Variable, q *WriteOptions) (*WriteMeta, error) {
	r, err := vars.client.newRequest("PUT", endpoint)
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	r.obj = in

	checkFn := requireStatusIn(http.StatusOK, http.StatusNoContent, http.StatusConflict) //nolint:bodyclose
	rtt, resp, err := checkFn(vars.client.doRequest(r))                                  //nolint:bodyclose

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	wm := &WriteMeta{RequestTime: rtt}
	_ = parseWriteMeta(resp, wm)

	if resp.StatusCode == http.StatusConflict {

		conflict := new(Variable)
		if err = decodeBody(resp, &conflict); err != nil {
			return nil, err
		}
		return nil, ErrCASConflict{
			Conflict:   conflict,
			CheckIndex: in.ModifyIndex,
		}
	}
	if out != nil {
		if err = decodeBody(resp, &out); err != nil {
			return nil, err
		}
	}
	return wm, nil
}

// Variable specifies the metadata and contents to be stored in the
// encrypted Nomad backend.
type Variable struct {
	// Namespace is the Nomad namespace associated with the variable
	Namespace string `hcl:"namespace"`

	// Path is the path to the variable
	Path string `hcl:"path"`

	// CreateIndex tracks the index of creation time
	CreateIndex uint64 `hcl:"create_index"`

	// ModifyTime is the unix nano of the last modified time
	ModifyIndex uint64 `hcl:"modify_index"`

	// CreateTime is the unix nano of the creation time
	CreateTime int64 `hcl:"create_time"`

	// ModifyTime is the unix nano of the last modified time
	ModifyTime int64 `hcl:"modify_time"`

	// Items contains the k/v variable component
	Items VariableItems `hcl:"items"`

	// Lock holds the information about the variable lock if its being used.
	Lock *VariableLock `hcl:",lock,optional" json:",omitempty"`
}

// VariableMetadata specifies the metadata for a variable and
// is used as the list object
type VariableMetadata struct {
	// Namespace is the Nomad namespace associated with the variable
	Namespace string `hcl:"namespace"`

	// Path is the path to the variable
	Path string `hcl:"path"`

	// CreateIndex tracks the index of creation time
	CreateIndex uint64 `hcl:"create_index"`

	// ModifyTime is the unix nano of the last modified time
	ModifyIndex uint64 `hcl:"modify_index"`

	// CreateTime is the unix nano of the creation time
	CreateTime int64 `hcl:"create_time"`

	// ModifyTime is the unix nano of the last modified time
	ModifyTime int64 `hcl:"modify_time"`

	// Lock holds the information about the variable lock if its being used.
	Lock *VariableLock `hcl:",lock,optional" json:",omitempty"`
}

type VariableLock struct {
	// ID is generated by Nomad to provide a unique caller ID which can be used
	// for renewals and unlocking.
	ID string

	// TTL describes the time-to-live of the current lock holder.
	// This is a string version of a time.Duration like "2m".
	TTL string

	// LockDelay describes a grace period that exists after a lock is lost,
	// before another client may acquire the lock. This helps protect against
	// split-brains. This is a string version of a time.Duration like "2m".
	LockDelay string
}

// VariableItems are the key/value pairs of a Variable.
type VariableItems map[string]string

// NewVariable is a convenience method to more easily create a
// ready-to-use variable
func NewVariable(path string) *Variable {
	return &Variable{
		Path:  path,
		Items: make(VariableItems),
	}
}

// Copy returns a new deep copy of this Variable
func (v *Variable) Copy() *Variable {
	var out = *v
	out.Items = make(VariableItems)
	for key, value := range v.Items {
		out.Items[key] = value
	}
	return &out
}

// Metadata returns the VariableMetadata component of
// a Variable. This can be useful for comparing against
// a List result.
func (v *Variable) Metadata() *VariableMetadata {
	return &VariableMetadata{
		Namespace:   v.Namespace,
		Path:        v.Path,
		CreateIndex: v.CreateIndex,
		ModifyIndex: v.ModifyIndex,
		CreateTime:  v.CreateTime,
		ModifyTime:  v.ModifyTime,
	}
}

// IsZeroValue can be used to test if a Variable has been changed
// from the default values it gets at creation
func (v *Variable) IsZeroValue() bool {
	return *v.Metadata() == VariableMetadata{} && v.Items == nil
}

// cleanPathString removes leading and trailing slashes since they
// would trigger go's path cleaning/redirection behavior in the
// standard HTTP router
func cleanPathString(path string) string {
	return strings.Trim(path, " /")
}

// AsJSON returns the Variable as a JSON-formatted string
func (v *Variable) AsJSON() string {
	var b []byte
	b, _ = json.Marshal(v)
	return string(b)
}

// AsPrettyJSON returns the Variable as a JSON-formatted string with
// indentation
func (v *Variable) AsPrettyJSON() string {
	var b []byte
	b, _ = json.MarshalIndent(v, "", "  ")
	return string(b)
}

// LockID returns the ID of the lock. In the event this is not held, or the
// variable is not a lock, this string will be empty.
func (v *Variable) LockID() string {
	if v.Lock == nil {
		return ""
	}

	return v.Lock.ID
}

type ErrCASConflict struct {
	CheckIndex uint64
	Conflict   *Variable
}

func (e ErrCASConflict) Error() string {
	return fmt.Sprintf("cas conflict: expected ModifyIndex %v; found %v", e.CheckIndex, e.Conflict.ModifyIndex)
}
