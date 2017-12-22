package okta

import (
	"encoding/json"
	"fmt"
	"github.com/chrismalek/oktasdk-go/okta"
	"github.com/hashicorp/vault/logical"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getMockClient(ts *httptest.Server) *okta.Client {
	var urlPtr url.URL
	tsURL, _ := urlPtr.Parse(ts.URL)
	return okta.NewClientWithBaseURL(nil, tsURL, "")
}

func successHandler(w http.ResponseWriter, req *http.Request) {
	fr := mfaResponse{
		FactorResult: "SUCCESS",
	}
	json.NewEncoder(w).Encode(fr)
}

func pollSuccessHandler(w http.ResponseWriter, req *http.Request) {
	mfaResponse := mfaResponse{
		FactorResult: "SUCCESS",
	}
	json.NewEncoder(w).Encode(mfaResponse)
}

func pollRejectHandler(w http.ResponseWriter, req *http.Request) {
	mfaResponse := mfaResponse{
		FactorResult: "REJECTED",
	}
	json.NewEncoder(w).Encode(mfaResponse)
}

func unauthorizedHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 Forbidden"))
}

func TestOktaHandlerOTPSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()

	successResp := &logical.Response{
		Auth: &logical.Auth{},
	}

	f := FactorLink{
		URL: ts.URL,
	}

	otpf := OTPFactor{
		Mfa: userMFAFactor{
			Links: struct {
				Self   FactorLink `json:"self"`
				Verify FactorLink `json:"verify"`
			}{Self: f, Verify: f},
		},
	}

	resp := otpf.DoAuth(getMockClient(ts), &oktaAuthRequest{successResp: successResp})
	fmt.Println(resp)
	if resp != successResp {
		t.Fatalf("Testing Okta authentication gave incorrect response (expected success, got: %v)", resp)
	}
}

func TestOktaHandlerOTPReject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(unauthorizedHandler))
	defer ts.Close()

	f := FactorLink{
		URL: ts.URL,
	}

	otpf := OTPFactor{
		Mfa: userMFAFactor{
			Links: struct {
				Self   FactorLink `json:"self"`
				Verify FactorLink `json:"verify"`
			}{Self: f, Verify: f},
		},
	}

	resp := otpf.DoAuth(getMockClient(ts), &oktaAuthRequest{})
	if resp.Data["error"] != "Invalid MFA passcode." {
		t.Fatalf("Testing Okta authentication gave incorrect response (expected deny, got: %v)", resp)
	}
}

func TestOktaPushVerified(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(pollSuccessHandler))
	defer ts.Close()

	flink := FactorLink{
		URL: ts.URL,
	}

	mfaLinks := struct {
		Self   FactorLink `json:"self"`
		Verify FactorLink `json:"verify"`
	}{Self: flink, Verify: flink}

	umfa := userMFAFactor{
		Links: mfaLinks,
	}

	pf := PushFactor{
		Mfa: umfa,
	}

	response := &oktaAuthRequest{}
	response.successResp = logical.ErrorResponse("Success")

	res := pf.DoAuth(getMockClient(ts), response)
	if res.Data["error"] != "Success" {
		t.Fatalf("Testing Okta Push authentication gave incorrect response (expected nil error, got: %v)", res)
	}
}

func TestOktaPushRejected(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(pollRejectHandler))
	defer ts.Close()

	flink := FactorLink{
		URL: ts.URL,
	}

	mfaLinks := struct {
		Self   FactorLink `json:"self"`
		Verify FactorLink `json:"verify"`
	}{Self: flink, Verify: flink}

	umfa := userMFAFactor{
		Links: mfaLinks,
	}

	pf := PushFactor{
		Mfa: umfa,
	}

	response := &oktaAuthRequest{}

	res := pf.DoAuth(getMockClient(ts), response)
	if res.Data["error"] != "Could not authenticate Okta user." {
		t.Fatalf("Testing Okta Push authentication gave incorrect response (expected error, got: %v)", res)
	}
}
