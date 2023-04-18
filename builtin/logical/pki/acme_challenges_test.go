package pki

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type keyAuthorizationTestCase struct {
	keyAuthz   string
	token      string
	thumbprint string
	shouldFail bool
}

var keyAuthorizationTestCases = []keyAuthorizationTestCase{
	{
		// Entirely empty
		"",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Both empty
		".",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Not equal
		"non-.non-",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Empty thumbprint
		"non-.",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Empty token
		".non-",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Wrong order
		"non-empty-thumbprint.non-empty-token",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Too many pieces
		"one.two.three",
		"non-empty-token",
		"non-empty-thumbprint",
		true,
	},
	{
		// Valid
		"non-empty-token.non-empty-thumbprint",
		"non-empty-token",
		"non-empty-thumbprint",
		false,
	},
}

func TestAcmeValidateKeyAuthorization(t *testing.T) {
	t.Parallel()

	for index, tc := range keyAuthorizationTestCases {
		isValid, err := ValidateKeyAuthorization(tc.keyAuthz, tc.token, tc.thumbprint)
		if !isValid && err == nil {
			t.Fatalf("[%d] expected failure to give reason via err (%v / %v)", index, isValid, err)
		}

		expectedValid := !tc.shouldFail
		if expectedValid != isValid {
			t.Fatalf("[%d] got ret=%v, expected ret=%v (shouldFail=%v)", index, isValid, expectedValid, tc.shouldFail)
		}
	}
}

func TestAcmeValidateHTTP01Challenge(t *testing.T) {
	t.Parallel()

	for index, tc := range keyAuthorizationTestCases {
		validFunc := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tc.keyAuthz))
		}
		withPadding := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("  " + tc.keyAuthz + "     "))
		}
		withRedirect := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/.well-known/") {
				http.Redirect(w, r, "/my-http-01-challenge-response", 301)
				return
			}

			w.Write([]byte(tc.keyAuthz))
		}
		withSleep := func(w http.ResponseWriter, r *http.Request) {
			// Long enough to ensure any excessively short timeouts are hit,
			// not long enough to trigger a failure (hopefully).
			time.Sleep(5 * time.Second)
			w.Write([]byte(tc.keyAuthz))
		}

		validHandlers := []http.HandlerFunc{
			http.HandlerFunc(validFunc), http.HandlerFunc(withPadding),
			http.HandlerFunc(withRedirect), http.HandlerFunc(withSleep),
		}

		for handlerIndex, handler := range validHandlers {
			func() {
				ts := httptest.NewServer(handler)
				defer ts.Close()

				host := ts.URL[7:]
				isValid, err := ValidateHTTP01Challenge(host, tc.token, tc.thumbprint)
				if !isValid && err == nil {
					t.Fatalf("[tc=%d/handler=%d] expected failure to give reason via err (%v / %v)", handlerIndex, index, isValid, err)
				}

				expectedValid := !tc.shouldFail
				if expectedValid != isValid {
					t.Fatalf("[tc=%d/handler=%d] got ret=%v (err=%v), expected ret=%v (shouldFail=%v)", index, handlerIndex, isValid, err, expectedValid, tc.shouldFail)
				}
			}()
		}
	}

	// Negative test cases for various HTTP-specific scenarios.
	redirectLoop := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/my-http-01-challenge-response", 301)
	}
	publicRedirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://hashicorp.com/", 301)
	}
	noData := func(w http.ResponseWriter, r *http.Request) {}
	noContent := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	notFound := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	simulateHang := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Second)
		w.Write([]byte("my-token.my-thumbprint"))
	}
	tooLarge := func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 512; i++ {
			w.Write([]byte("my-token.my-thumbprint"))
		}
	}

	validHandlers := []http.HandlerFunc{
		http.HandlerFunc(redirectLoop), http.HandlerFunc(publicRedirect),
		http.HandlerFunc(noData), http.HandlerFunc(noContent),
		http.HandlerFunc(notFound), http.HandlerFunc(simulateHang),
		http.HandlerFunc(tooLarge),
	}
	for handlerIndex, handler := range validHandlers {
		func() {
			ts := httptest.NewServer(handler)
			defer ts.Close()

			host := ts.URL[7:]
			isValid, err := ValidateHTTP01Challenge(host, "my-token", "my-thumbprint")
			if isValid || err == nil {
				t.Fatalf("[handler=%d] expected failure validating challenge (%v / %v)", handlerIndex, isValid, err)
			}
		}()
	}
}
