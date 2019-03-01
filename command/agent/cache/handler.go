package cache

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/helper/consts"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
)

func Handler(ctx context.Context, logger hclog.Logger, proxier Proxier, inmemSink sink.Sink) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request", "path", r.URL.Path, "method", r.Method)

		token := r.Header.Get(consts.AuthHeaderName)
		if token == "" && inmemSink != nil {
			logger.Debug("using auto auth token", "path", r.URL.Path, "method", r.Method)
			token = inmemSink.(sink.SinkReader).Token()
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
		r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		req := &SendRequest{
			Token:       token,
			Request:     r,
			RequestBody: reqBody,
		}

		resp, err := proxier.Send(ctx, req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errwrap.Wrapf("failed to get the response: {{err}}", err))
			return
		}

		err = processTokenLookupResponse(ctx, logger, inmemSink, req, resp)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errwrap.Wrapf("failed to process token lookup response: {{err}}", err))
			return
		}

		defer resp.Response.Body.Close()

		copyHeader(w.Header(), resp.Response.Header)
		w.WriteHeader(resp.Response.StatusCode)
		io.Copy(w, resp.Response.Body)
		return
	})
}

// processTokenLookupResponse checks if the request was one of token
// lookup-self. If the auto-auth token was used to perform lookup-self, the
// identifier of the token and its accessor same will be stripped off of the
// response.
func processTokenLookupResponse(ctx context.Context, logger hclog.Logger, inmemSink sink.Sink, req *SendRequest, resp *SendResponse) error {
	// If auto-auth token is not being used, there is nothing to do.
	if inmemSink == nil {
		return nil
	}
	autoAuthToken := inmemSink.(sink.SinkReader).Token()

	// If lookup responded with non 200 status, there is nothing to do.
	if resp.Response.StatusCode != http.StatusOK {
		return nil
	}

	_, path := deriveNamespaceAndRevocationPath(req)
	switch path {
	case vaultPathTokenLookupSelf:
		if req.Token != autoAuthToken {
			return nil
		}
	case vaultPathTokenLookup:
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return err
		}
		tokenRaw, ok := jsonBody["token"]
		if !ok {
			// Input error will be caught by the API
			return nil
		}
		token, ok := tokenRaw.(string)
		if !ok {
			// Input error will be caught by the API
			return nil
		}
		if token != "" && token != autoAuthToken {
			// Lookup is performed on the non-auto-auth token
			return nil
		}
	default:
		return nil
	}

	logger.Info("stripping auto-auth token from the response", "path", req.Request.URL.Path, "method", req.Request.Method)
	secret, err := api.ParseSecret(bytes.NewReader(resp.ResponseBody))
	if err != nil {
		return fmt.Errorf("failed to parse token lookup response: %v", err)
	}
	if secret == nil || secret.Data == nil {
		return nil
	}
	if secret.Data["id"] == nil && secret.Data["accessor"] == nil {
		return nil
	}

	delete(secret.Data, "id")
	delete(secret.Data, "accessor")

	bodyBytes, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	if resp.Response.Body != nil {
		resp.Response.Body.Close()
	}
	resp.Response.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	resp.Response.ContentLength = int64(len(bodyBytes))

	// Serialize and re-read the reponse
	var respBytes bytes.Buffer
	err = resp.Response.Write(&respBytes)
	if err != nil {
		return fmt.Errorf("failed to serialize the updated response: %v", err)
	}

	updatedResponse, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes.Bytes())), nil)
	if err != nil {
		return fmt.Errorf("failed to deserialize the updated response: %v", err)
	}

	resp.Response = &api.Response{
		Response: updatedResponse,
	}
	resp.ResponseBody = bodyBytes

	return nil
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
