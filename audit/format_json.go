package audit

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONFormatWriter is an AuditFormatWriter implementation that structures data into
// a JSON format.
type JSONFormatWriter struct{}

func (f *JSONFormatWriter) WriteRequest(w io.Writer, req *AuditRequestEntry) error {
	if req == nil {
		return fmt.Errorf("request entry was nil, cannot encode")
	}

	enc := json.NewEncoder(w)
	return enc.Encode(req)
}

func (f *JSONFormatWriter) WriteResponse(w io.Writer, resp *AuditResponseEntry) error {
	if resp == nil {
		return fmt.Errorf("response entry was nil, cannot encode")
	}

	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}
