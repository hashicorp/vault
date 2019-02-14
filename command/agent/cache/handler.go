package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
)

func Handler(ctx context.Context, logger hclog.Logger, proxier Proxier, useAutoAuthToken bool, client *api.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request", "path", r.URL.Path, "method", r.Method)

		token := r.Header.Get(consts.AuthHeaderName)
		if token == "" && useAutoAuthToken {
			logger.Debug("using auto auth token")
			token = client.Token()
		}

		// Parse and reset body.
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body")
			respondError(w, http.StatusInternalServerError, errors.New("failed to read request body"))
		}
		if r.Body != nil {
			r.Body.Close()
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		resp, err := proxier.Send(ctx, &SendRequest{
			Token:       token,
			Request:     r,
			RequestBody: reqBody,
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, errwrap.Wrapf("failed to get the response: {{err}}", err))
			return
		}

		defer resp.Response.Body.Close()

		copyHeader(w.Header(), resp.Response.Header)
		w.WriteHeader(resp.Response.StatusCode)
		io.Copy(w, resp.Response.Body)
		return
	})
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func respondError(w http.ResponseWriter, status int, err error) {
	logical.AdjustErrorStatusCode(&status, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &vaulthttp.ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
