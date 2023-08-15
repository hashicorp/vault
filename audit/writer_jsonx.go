// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jefferai/jsonx"
)

var _ Writer = (*JSONxWriter)(nil)

// JSONxWriter is a Writer implementation that structures data into an XML format.
type JSONxWriter struct {
	Prefix string
}

func (f *JSONxWriter) WriteRequest(w io.Writer, req *RequestEntry) error {
	if req == nil {
		return fmt.Errorf("request entry was nil, cannot encode")
	}

	if len(f.Prefix) > 0 {
		_, err := w.Write([]byte(f.Prefix))
		if err != nil {
			return err
		}
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

func (f *JSONxWriter) WriteResponse(w io.Writer, resp *ResponseEntry) error {
	if resp == nil {
		return fmt.Errorf("response entry was nil, cannot encode")
	}

	if len(f.Prefix) > 0 {
		_, err := w.Write([]byte(f.Prefix))
		if err != nil {
			return err
		}
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
