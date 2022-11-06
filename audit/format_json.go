package audit

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/hashicorp/vault/sdk/helper/salt"
)

// JSONFormatWriter is an AuditFormatWriter implementation that structures data into
// a JSON format.
type JSONFormatWriter struct {
	Prefix   string
	SaltFunc func(context.Context) (*salt.Salt, error)
}

func (f *JSONFormatWriter) WriteRequest(w io.Writer, req *AuditRequestEntry) error {
	if req == nil {
		return errors.New("request entry was nil, cannot encode")
	}

	if len(f.Prefix) > 0 {
		if _, err := w.Write([]byte(f.Prefix)); err != nil {
			return err
		}
	}

	enc := json.NewEncoder(w)
	return enc.Encode(req)
}

func (f *JSONFormatWriter) WriteResponse(w io.Writer, resp *AuditResponseEntry) error {
	if resp == nil {
		return errors.New("response entry was nil, cannot encode")
	}

	if len(f.Prefix) > 0 {
		if _, err := w.Write([]byte(f.Prefix)); err != nil {
			return err
		}
	}

	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}

func (f *JSONFormatWriter) Salt(ctx context.Context) (*salt.Salt, error) {
	return f.SaltFunc(ctx)
}
