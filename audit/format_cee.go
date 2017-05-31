package audit

import (
	"encoding/json"
	"fmt"
	"io"
)

// CEEFormatWriter is an AuditFormatWriter implementation that structures data into
// a JSON format.
type CEEFormatWriter struct{}

const ceePrefix = "@cee:"

func (f *CEEFormatWriter) WriteRequest(w io.Writer, req *AuditRequestEntry) error {
	if req == nil {
		return fmt.Errorf("request entry was nil, cannot encode")
	}
	_, err := w.Write([]byte(ceePrefix))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	return enc.Encode(req)
}

func (f *CEEFormatWriter) WriteResponse(w io.Writer, resp *AuditResponseEntry) error {
	if resp == nil {
		return fmt.Errorf("response entry was nil, cannot encode")
	}
	_, err := w.Write([]byte(ceePrefix))
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}
