package listenerutil

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
)

// DefaultStatus is used to set default headers early before having a status code,
// for example, for /ui headers
const DefaultStatus = 1

func SetCustomResponseHeaders(hm map[string]map[string]string, w http.ResponseWriter, status int) error {
	// Removing X-Vault-Listener-Add header from ResponseWriter
	// This should be safe as the call to this function is right
	// before w.WriteHeader for which the status code is finalized and known
	w.Header().Del("X-Vault-Listener-Add")

	if hm == nil {
		return nil
	}

	// setter function to set the headers
	setter := func(hv map[string]string) {
		for h, v := range hv {
			w.Header().Set(h, v)
		}
	}

	// Checking the validity of the status code
	if status >= 600 || (status < 100 && status != DefaultStatus) {
		return fmt.Errorf("invalid status code")
	}

	// Setting the default headers first
	setter(hm["default"])

	// for DefaultStatus, we only set the default headers
	if status == DefaultStatus {
		return nil
	}

	// setting the Xyy pattern first
	d := fmt.Sprintf("%vxx", status / 100)
	if val, ok := hm[d]; ok {
		setter(val)
	}
	// Setting the specific headers
	if val, ok := hm[strconv.Itoa(status)]; ok {
		setter(val)
	}

	return nil
}

func FetchCustomResponseHeaderValue(hm map[string]map[string]string, th string, sc int) (string, error) {
	if hm == nil {
		return "", nil
	}
	if th == "" {
		return "", fmt.Errorf("invalid target header")
	}

	var h map[string]string
	if sc == DefaultStatus {
		h = hm["default"]
	}else {
		h = hm[strconv.Itoa(sc)]
	}

	hn := textproto.CanonicalMIMEHeaderKey(th)
	if v, ok := h[hn]; ok {
		return v, nil
	}
	return "", nil
}

func ExistHeader(hm map[string]map[string]string, th string, sl []int) bool {
	if len(sl) == 0 {
		return false
	}

	for _, s := range sl {
		chv, _ := FetchCustomResponseHeaderValue(hm, th, s)
		if chv != "" {
			return true
		}
	}

	return false
}