//
// Copyright 2020 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// This file is used for tracing HTTP requests (trace).
//
// All HTTP requests and responses through this transport will be printed to
// stderr.

// TraceRoundTripper to wrap a HTTP Transport.
func TraceRoundTripper(in http.RoundTripper) http.RoundTripper {
	return &traceRoundTripper{inner: in, logger: os.Stderr}
}

type traceRoundTripper struct {
	inner  http.RoundTripper
	logger io.Writer
}

func (d *traceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	d.dumpRequest(req)
	res, err := d.inner.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(d.logger, "\n\tERROR for request: %v\n", err)
	}
	if res != nil {
		d.dumpResponse(res)
	}
	return res, err
}

func (d *traceRoundTripper) dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		fmt.Fprintf(d.logger, "\n\tERROR dumping: %v\n", err)
		return
	}
	d.dump("REQUEST", dump)
}

func (d *traceRoundTripper) dumpResponse(r *http.Response) {
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		fmt.Fprintf(d.logger, "\n\tERROR dumping: %v\n", err)
		return
	}
	d.dump("RESPONSE", dump)
}

func (d *traceRoundTripper) dump(label string, dump []byte) {
	fmt.Fprintf(d.logger, "\n%s:\n--\n%s\n", label, string(dump))
	if label == "RESPONSE" {
		fmt.Fprintf(d.logger, "\n")
	}
}
