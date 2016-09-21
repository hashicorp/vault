package audit

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jefferai/jsonx"
)

// JSONxFormatWriter is an AuditFormatWriter implementation that structures data into
// a XML format.
type JSONxFormatWriter struct{}

func (f *JSONxFormatWriter) WriteRequest(w io.Writer, req *AuditRequestEntry) error {
	if req == nil {
		return fmt.Errorf("request entry was nil, cannot encode")
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	xmlBytes, err := jsonx.EncodeJSONBytes(jsonBytes)
	if err != nil {
		return err
	}

	_, err = w.Write(xmlBytes)
	return err
}

func (f *JSONxFormatWriter) WriteResponse(w io.Writer, resp *AuditResponseEntry) error {
	if resp == nil {
		return fmt.Errorf("response entry was nil, cannot encode")
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	xmlBytes, err := jsonx.EncodeJSONBytes(jsonBytes)
	if err != nil {
		return err
	}

	_, err = w.Write(xmlBytes)
	return err
}
