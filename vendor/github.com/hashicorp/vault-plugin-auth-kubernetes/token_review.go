// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	authv1 "k8s.io/api/authentication/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// This is the result from the token review
type tokenReviewResult struct {
	Name      string
	Namespace string
	UID       string
}

// This exists so we can use a mock TokenReview when running tests
type tokenReviewer interface {
	Review(context.Context, *http.Client, string, []string) (*tokenReviewResult, error)
}

type tokenReviewFactory func(*kubeConfig) tokenReviewer

// This is the real implementation that calls the kubernetes API
type tokenReviewAPI struct {
	config *kubeConfig
}

func tokenReviewAPIFactory(config *kubeConfig) tokenReviewer {
	return &tokenReviewAPI{
		config: config,
	}
}

func (t *tokenReviewAPI) Review(ctx context.Context, client *http.Client, jwt string, aud []string) (*tokenReviewResult, error) {
	// Create the TokenReview Object and marshal it into json
	trReq := &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token:     jwt,
			Audiences: aud,
		},
	}
	trJSON, err := json.Marshal(trReq)
	if err != nil {
		return nil, err
	}

	// Build the request to the token review API
	url := fmt.Sprintf("%s/apis/authentication.k8s.io/v1/tokenreviews", strings.TrimSuffix(t.config.Host, "/"))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(trJSON))
	if err != nil {
		return nil, err
	}

	// If we have a configured TokenReviewer JWT use it as the bearer, otherwise
	// try to use the passed in JWT.
	bearer := fmt.Sprintf("Bearer %s", jwt)
	if len(t.config.TokenReviewerJWT) > 0 {
		bearer = fmt.Sprintf("Bearer %s", t.config.TokenReviewerJWT)
	}
	setRequestHeader(req, bearer)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse the resp into a tokenreview object or a kubernetes error type
	r, err := parseResponse(resp)
	switch {
	case kubeerrors.IsUnauthorized(err):
		// If the err is unauthorized that means the token has since been deleted;
		// this can happen if the service account is deleted, and even if it has
		// since been recreated the token will have changed, which means our
		// caller will need to be updated accordingly.
		return nil, errors.New("lookup failed: service account unauthorized; this could mean it has been deleted or recreated with a new token")
	case err != nil:
		return nil, err
	}

	if r.Status.Error != "" {
		return nil, fmt.Errorf("lookup failed: %s", r.Status.Error)
	}

	if !r.Status.Authenticated {
		return nil, errors.New("lookup failed: service account jwt not valid")
	}

	// Ensure the token review endpoint is audience-aware if we requested
	// audience validation.
	wantAud := trReq.Spec.Audiences
	if len(wantAud) != 0 {
		intersectionFound := false
		for _, aud := range trReq.Spec.Audiences {
			if strutil.StrListContains(r.Status.Audiences, aud) {
				intersectionFound = true
				break
			}
		}
		if !intersectionFound {
			return nil, fmt.Errorf("lookup failed: service account jwt valid for audience(s) %v, but wanted %v", r.Status.Audiences, wantAud)
		}
	}

	// The username is of format: system:serviceaccount:(NAMESPACE):(SERVICEACCOUNT)
	parts := strings.Split(r.Status.User.Username, ":")
	if len(parts) != 4 {
		return nil, errors.New("lookup failed: unexpected username format")
	}

	// Validate the user that comes back from token review is a service account
	if parts[0] != "system" || parts[1] != "serviceaccount" {
		return nil, errors.New("lookup failed: username returned is not a service account")
	}

	return &tokenReviewResult{
		Name:      parts[3],
		Namespace: parts[2],
		UID:       string(r.Status.User.UID),
	}, nil
}

// parseResponse takes the API response and either returns the appropriate error
// or the TokenReview Object.
func parseResponse(resp *http.Response) (*authv1.TokenReview, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If the request was not a success create a kuberenets error
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return nil, kubeerrors.NewGenericServerResponse(resp.StatusCode, "POST", schema.GroupResource{}, "", strings.TrimSpace(string(body)), 0, true)
	}

	// If we can successfully Unmarshal into a status object that means there is
	// an error to return
	errStatus := &metav1.Status{}
	err = json.Unmarshal(body, errStatus)
	if err == nil && errStatus.Status != metav1.StatusSuccess {
		return nil, kubeerrors.FromObject(runtime.Object(errStatus))
	}

	// Unmarshal the resp body into a TokenReview Object
	trResp := &authv1.TokenReview{}
	err = json.Unmarshal(body, trResp)
	if err != nil {
		return nil, err
	}

	return trResp, nil
}

// mock review is used while testing
type mockTokenReview struct {
	saName      string
	saNamespace string
	saUID       string
}

func mockTokenReviewFactory(name, namespace, UID string) tokenReviewFactory {
	return func(config *kubeConfig) tokenReviewer {
		return &mockTokenReview{
			saName:      name,
			saNamespace: namespace,
			saUID:       UID,
		}
	}
}

func (t *mockTokenReview) Review(ctx context.Context, client *http.Client, cjwt string, aud []string) (*tokenReviewResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	httpTransport, ok := client.Transport.(*http.Transport)
	if !ok {
		return nil, errors.New("failed to check whether DisableKeepAlives is false as Transport is not *http.Transport")
	}
	if httpTransport.DisableKeepAlives {
		return nil, errors.New("expected DisableKeepAlives to be false but was true")
	}

	return &tokenReviewResult{
		Name:      t.saName,
		Namespace: t.saNamespace,
		UID:       t.saUID,
	}, nil
}
