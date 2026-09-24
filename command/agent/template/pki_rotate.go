// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package template

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type pkiRotationRequest struct {
	ctx         context.Context
	destination string
	result      chan int
}

// PKIRotationHandler exposes an opt-in, asynchronous local administrative API.
// The listener must be restricted to trusted callers: rotation can issue new
// certificates and execute the commands already configured for the template.
func (ts *Server) PKIRotationHandler(enabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enabled {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Destination string `json:"destination"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil || body.Destination == "" {
			http.Error(w, "a destination is required", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			http.Error(w, "expected one JSON object", http.StatusBadRequest)
			return
		}
		if ts == nil || ts.config == nil || ts.config.AgentConfig == nil {
			http.Error(w, "template server unavailable", http.StatusServiceUnavailable)
			return
		}
		found := false
		for _, tmpl := range ts.config.AgentConfig.Templates {
			if tmpl.Destination != nil && *tmpl.Destination == body.Destination {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "unknown template destination", http.StatusNotFound)
			return
		}
		if ts.exitAfterAuth {
			http.Error(w, "rotation is unavailable with exit_after_auth", http.StatusConflict)
			return
		}
		if !ts.runnerStarted.Load() {
			http.Error(w, "template runner is not ready", http.StatusServiceUnavailable)
			return
		}
		// The server loop serializes this request with runner replacement on reauth.
		// This timeout bounds acceptance, not certificate issuance or command execution.
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		request := pkiRotationRequest{ctx: ctx, destination: body.Destination, result: make(chan int, 1)}
		select {
		case ts.rotateRequests <- request:
		case <-ts.DoneCh:
			http.Error(w, "template server stopped", http.StatusServiceUnavailable)
			return
		case <-ctx.Done():
			http.Error(w, "template server unavailable", http.StatusServiceUnavailable)
			return
		}
		select {
		case count := <-request.result:
			if count == 0 {
				http.Error(w, "no active pkiCert dependency for destination", http.StatusConflict)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]any{"destination": body.Destination, "status": "accepted", "dependencies": count})
		case <-ts.DoneCh:
			http.Error(w, "template server stopped", http.StatusServiceUnavailable)
		case <-ctx.Done():
			http.Error(w, "template server unavailable", http.StatusServiceUnavailable)
		}
	})
}
