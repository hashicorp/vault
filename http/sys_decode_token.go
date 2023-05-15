// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/hashicorp/vault/sdk/helper/xor"
	"github.com/hashicorp/vault/vault"
)

func handleSysDecodeToken(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT", "POST":
			handleSysDecodeTokenPut(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysDecodeTokenPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req DecodeTokenRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil && err != io.EOF {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.EncodedToken == "" {
		respondError(
			w, http.StatusBadRequest,
			errors.New("'encoded_token' must be specified in request body as JSON"))
		return
	}

	if req.OTP == "" {
		respondError(
			w, http.StatusBadRequest,
			errors.New("'otp' must be specified in request body as JSON"))
		return
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(req.EncodedToken)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(req.OTP))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	token := string(tokenBytes)

	status := &DecodeTokenResponse{
		Token: token,
	}

	respondOk(w, status)
}

type DecodeTokenRequest struct {
	EncodedToken string `json:"encoded_token"`
	OTP          string `json:"otp"`
}

type DecodeTokenResponse struct {
	Token string `json:"token"`
}
