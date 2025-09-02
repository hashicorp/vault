// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
)

const (
	wrappedResponseLocation = "cubbyhole/response"
)

var (
	// The default TTL that will be used with `sys/wrapping/wrap`, can be
	// changed
	DefaultWrappingTTL = "5m"

	// The default function used if no other function is set. It honors the env
	// var to set the wrap TTL. The default wrap TTL will apply when when writing
	// to `sys/wrapping/wrap` when the env var is not set.
	DefaultWrappingLookupFunc = func(operation, path string) string {
		if os.Getenv(EnvVaultWrapTTL) != "" {
			return os.Getenv(EnvVaultWrapTTL)
		}

		if (operation == http.MethodPut || operation == http.MethodPost) && path == "sys/wrapping/wrap" {
			return DefaultWrappingTTL
		}

		return ""
	}
)

// Logical is used to perform logical backend operations on Vault.
type Logical struct {
	c *Client
}

// Logical is used to return the client for logical-backend API calls.
func (c *Client) Logical() *Logical {
	return &Logical{c: c}
}

func (c *Logical) Read(path string) (*Secret, error) {
	return c.ReadWithDataWithContext(context.Background(), path, nil)
}

func (c *Logical) ReadWithContext(ctx context.Context, path string) (*Secret, error) {
	return c.ReadWithDataWithContext(ctx, path, nil)
}

func (c *Logical) ReadWithData(path string, data map[string][]string) (*Secret, error) {
	return c.ReadWithDataWithContext(context.Background(), path, data)
}

// ReadFromSnapshot reads the data at the given Vault path from a previously
// loaded snapshot. The snapshotID parameter is the ID of the loaded snapshot
func (c *Logical) ReadFromSnapshot(path string, snapshotID string) (*Secret, error) {
	return c.ReadWithData(path, map[string][]string{"read_snapshot_id": {snapshotID}})
}

func (c *Logical) ReadWithDataWithContext(ctx context.Context, path string, data map[string][]string) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	resp, err := c.readRawWithDataWithContext(ctx, path, data, nil)
	return c.ParseRawResponseAndCloseBody(resp, err)
}

// ReadWithRequest returns a Secret for the given LogicalReadRequest. This is a
// more flexible version of ReadWithContext, which allows for passing extra
// headers to the Vault server.
func (c *Logical) ReadWithRequest(ctx context.Context, req LogicalReadRequest) (*Secret, error) {
	resp, err := c.readRawWithDataWithContext(ctx, req.Path(), req.Values(), req.Headers())
	return c.ParseRawResponseAndCloseBody(resp, err)
}

// ReadRaw attempts to read the value stored at the given Vault path
// (without '/v1/' prefix) and returns a raw *http.Response.
//
// Note: the raw-response functions do not respect the client-configured
// request timeout; if a timeout is desired, please use ReadRawWithContext
// instead and set the timeout through context.WithTimeout or context.WithDeadline.
func (c *Logical) ReadRaw(path string) (*Response, error) {
	return c.ReadRawWithDataWithContext(context.Background(), path, nil)
}

// ReadRawWithContext attempts to read the value stored at the give Vault path
// (without '/v1/' prefix) and returns a raw *http.Response.
//
// Note: the raw-response functions do not respect the client-configured
// request timeout; if a timeout is desired, please set it through
// context.WithTimeout or context.WithDeadline.
func (c *Logical) ReadRawWithContext(ctx context.Context, path string) (*Response, error) {
	return c.ReadRawWithDataWithContext(ctx, path, nil)
}

// ReadRawWithData attempts to read the value stored at the given Vault
// path (without '/v1/' prefix) and returns a raw *http.Response. The 'data' map
// is added as query parameters to the request.
//
// Note: the raw-response functions do not respect the client-configured
// request timeout; if a timeout is desired, please use
// ReadRawWithDataWithContext instead and set the timeout through
// context.WithTimeout or context.WithDeadline.
func (c *Logical) ReadRawWithData(path string, data map[string][]string) (*Response, error) {
	return c.ReadRawWithDataWithContext(context.Background(), path, data)
}

func (c *Logical) ReadRawFromSnapshot(path string, snapshotID string) (*Response, error) {
	return c.ReadRawWithDataWithContext(context.Background(), path, map[string][]string{"read_snapshot_id": {snapshotID}})
}

// ReadRawWithDataWithContext attempts to read the value stored at the given
// Vault path (without '/v1/' prefix) and returns a raw *http.Response. The 'data'
// map is added as query parameters to the request.
//
// Note: the raw-response functions do not respect the client-configured
// request timeout; if a timeout is desired, please set it through
// context.WithTimeout or context.WithDeadline.
func (c *Logical) ReadRawWithDataWithContext(ctx context.Context, path string, data map[string][]string) (*Response, error) {
	return c.readRawWithDataWithContext(ctx, path, data, nil)
}

func (c *Logical) ParseRawResponseAndCloseBody(resp *Response, err error) (*Secret, error) {
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, parseErr
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) readRawWithDataWithContext(ctx context.Context, path string, values url.Values, extraHeaders http.Header) (*Response, error) {
	r := c.c.NewRequest(http.MethodGet, "/v1/"+path)

	if err := c.addExtraHeaders(r, extraHeaders); err != nil {
		return nil, err
	}

	if values != nil {
		r.Params = maps.Clone(values)
	}

	return c.c.RawRequestWithContext(ctx, r)
}

// ListFromSnapshot lists from the Vault path using a previously loaded
// snapshot. The snapshotID parameter is the ID of the loaded snapshot
func (c *Logical) ListFromSnapshot(path string, snapshotID string) (*Secret, error) {
	r := c.c.NewRequest("LIST", "/v1/"+path)
	r.Params.Set("read_snapshot_id", snapshotID)
	return c.list(context.Background(), r)
}

func (c *Logical) List(path string) (*Secret, error) {
	return c.ListWithContext(context.Background(), path)
}

func (c *Logical) ListWithContext(ctx context.Context, path string) (*Secret, error) {
	return c.list(ctx, c.c.NewRequest("LIST", "/v1/"+path))
}

func (c *Logical) list(ctx context.Context, r *Request) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	// Set this for broader compatibility, but we use LIST above to be able to
	// handle the wrapping lookup function
	r.Method = http.MethodGet
	r.Params.Set("list", "true")

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, parseErr
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) Write(path string, data map[string]interface{}) (*Secret, error) {
	return c.WriteWithContext(context.Background(), path, data)
}

func (c *Logical) WriteWithContext(ctx context.Context, path string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest(http.MethodPut, "/v1/"+path)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	return c.write(ctx, path, r)
}

func (c *Logical) WriteRaw(path string, data []byte) (*Response, error) {
	return c.WriteRawWithContext(context.Background(), path, data)
}

func (c *Logical) WriteRawWithContext(ctx context.Context, path string, data []byte) (*Response, error) {
	r := c.c.NewRequest(http.MethodPut, "/v1/"+path)
	r.BodyBytes = data

	return c.writeRaw(ctx, r)
}

// WriteWithRequest returns a Secret for the given LogicalRequest. This is a
// more flexible version of WriteWithContext, which allows for passing extra
// headers to the Vault server.
func (c *Logical) WriteWithRequest(ctx context.Context, req LogicalWriteRequest) (*Secret, error) {
	r := c.c.NewRequest(http.MethodPut, "/v1/"+req.Path())
	if err := c.addExtraHeaders(r, req.Headers()); err != nil {
		return nil, err
	}

	if err := r.SetJSONBody(req.Data()); err != nil {
		return nil, err
	}

	return c.write(ctx, req.Path(), r)
}

// addExtraHeaders adds the given headers to the request, but only if they are
// not already set in the request. If a header is already set in the request, it
// returns an error for each header that is already set.
func (c *Logical) addExtraHeaders(r *Request, headers http.Header) error {
	var errs error

	if r == nil {
		return fmt.Errorf("cannot add extra headers to nil request")
	}

	if len(headers) == 0 {
		return nil
	}

	if r.Headers == nil {
		r.Headers = headers
		return nil
	}

	curHeaders := r.Headers.Clone()
	for k, v := range headers {
		ck := http.CanonicalHeaderKey(k)
		if curVal := curHeaders.Get(ck); curVal != "" {
			errs = errors.Join(errs, fmt.Errorf("cannot set extra header %q, it is reserved", ck))
			continue
		}
		curHeaders[ck] = v
	}
	if errs != nil {
		return fmt.Errorf("cannot add extra headers: %w", errs)
	}

	r.Headers = curHeaders
	return nil
}

// Recover recovers the data at the given Vault path from a loaded snapshot.
// The snapshotID parameter is the ID of the loaded snapshot
func (c *Logical) Recover(ctx context.Context, path string, snapshotID string) (*Secret, error) {
	return c.RecoverFromPath(ctx, path, snapshotID, "")
}

func (c *Logical) RecoverFromPath(ctx context.Context, newPath string, snapshotID string, originalPath string) (*Secret, error) {
	r := c.c.NewRequest(http.MethodPut, "/v1/"+newPath)
	r.Params.Set("recover_snapshot_id", snapshotID)
	r.Headers.Set(SnapshotHeaderName, snapshotID)
	if originalPath != "" && originalPath != newPath {
		r.Headers.Set(RecoverSourcePathHeaderName, originalPath)
	}
	return c.write(ctx, originalPath, r)
}

func (c *Logical) JSONMergePatch(ctx context.Context, path string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest(http.MethodPatch, "/v1/"+path)
	r.Headers.Set("Content-Type", "application/merge-patch+json")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	return c.write(ctx, path, r)
}

func (c *Logical) WriteBytes(path string, data []byte) (*Secret, error) {
	return c.WriteBytesWithContext(context.Background(), path, data)
}

func (c *Logical) WriteBytesWithContext(ctx context.Context, path string, data []byte) (*Secret, error) {
	r := c.c.NewRequest(http.MethodPut, "/v1/"+path)
	r.BodyBytes = data

	return c.write(ctx, path, r)
}

func (c *Logical) write(ctx context.Context, path string, request *Request) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	resp, err := c.c.rawRequestWithContext(ctx, request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, parseErr
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) writeRaw(ctx context.Context, request *Request) (*Response, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	resp, err := c.c.rawRequestWithContext(ctx, request)
	return resp, err
}

func (c *Logical) Delete(path string) (*Secret, error) {
	return c.DeleteWithContext(context.Background(), path)
}

func (c *Logical) DeleteWithContext(ctx context.Context, path string) (*Secret, error) {
	return c.DeleteWithDataWithContext(ctx, path, nil)
}

func (c *Logical) DeleteWithData(path string, data map[string][]string) (*Secret, error) {
	return c.DeleteWithDataWithContext(context.Background(), path, data)
}

// DeleteWithRequest returns a Secret for the given LogicalDeleteRequest. This is a
// more flexible version of DeleteWithContext, which allows for passing extra
// headers to the Vault server.
func (c *Logical) DeleteWithRequest(ctx context.Context, req LogicalDeleteRequest) (*Secret, error) {
	return c.DeleteWithDataWithContext(ctx, req.Path(), req.Values())
}

func (c *Logical) DeleteWithDataWithContext(ctx context.Context, path string, data map[string][]string) (*Secret, error) {
	return c.deleteWithDataWithContext(ctx, path, data)
}

func (c *Logical) deleteWithDataWithContext(ctx context.Context, path string, data map[string][]string) (*Secret, error) {
	return c.delete(ctx, NewDeleteRequest(path, data, nil))
}

func (c *Logical) delete(ctx context.Context, req LogicalDeleteRequest) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodDelete, "/v1/"+req.Path())

	if values := req.Values(); values != nil {
		r.Params = values
	}

	if err := c.addExtraHeaders(r, req.Headers()); err != nil {
		return nil, err
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, parseErr
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) Unwrap(wrappingToken string) (*Secret, error) {
	return c.UnwrapWithContext(context.Background(), wrappingToken)
}

func (c *Logical) UnwrapWithContext(ctx context.Context, wrappingToken string) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	var data map[string]interface{}
	wt := strings.TrimSpace(wrappingToken)
	if wrappingToken != "" {
		if c.c.Token() == "" {
			c.c.SetToken(wt)
		} else if wrappingToken != c.c.Token() {
			data = map[string]interface{}{
				"token": wt,
			}
		}
	}

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/wrapping/unwrap")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp == nil || resp.StatusCode != 404 {
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, nil
		}
		return ParseSecret(resp.Body)
	}

	// In the 404 case this may actually be a wrapped 404 error
	secret, parseErr := ParseSecret(resp.Body)
	switch parseErr {
	case nil:
	case io.EOF:
		return nil, nil
	default:
		return nil, parseErr
	}
	if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
		return secret, nil
	}

	// Otherwise this might be an old-style wrapping token so attempt the old
	// method
	if wrappingToken != "" {
		origToken := c.c.Token()
		defer c.c.SetToken(origToken)
		c.c.SetToken(wrappingToken)
	}

	secret, err = c.ReadWithContext(ctx, wrappedResponseLocation)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error reading %q: {{err}}", wrappedResponseLocation), err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no value found at %q", wrappedResponseLocation)
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("\"data\" not found in wrapping response")
	}
	if _, ok := secret.Data["response"]; !ok {
		return nil, fmt.Errorf("\"response\" not found in wrapping response \"data\" map")
	}

	wrappedSecret := new(Secret)
	buf := bytes.NewBufferString(secret.Data["response"].(string))
	dec := json.NewDecoder(buf)
	dec.UseNumber()
	if err := dec.Decode(wrappedSecret); err != nil {
		return nil, errwrap.Wrapf("error unmarshalling wrapped secret: {{err}}", err)
	}

	return wrappedSecret, nil
}
