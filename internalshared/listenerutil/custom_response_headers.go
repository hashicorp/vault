package listenerutil

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/internalshared/configutil"
)

type ListenerCustomHeaders struct {
	Address string
	StatusCodeHeaderMap map[string][]*CustomHeader
	// ConfiguredHeadersStatusCodeMap field is introduced so that we would not need to loop through
	// StatusCodeHeaderMap to see if a header exists, the key for this map is the headers names
	ConfiguredHeadersStatusCodeMap map[string][]string
}

type CustomHeader struct {
	Name string
	Value string
}

// ChangeListenerAddress is used for tests where the listener address (at least the port)
// is chosen at random
func (l *ListenerCustomHeaders) ChangeListenerAddress(la string) {
	l.Address = la
	return
}

func (l *ListenerCustomHeaders) SetCustomResponseHeaders(w http.ResponseWriter, status int) {
	if w == nil {
		fmt.Println("No ResponseWriter provided")
		return
	}

	sch := l.StatusCodeHeaderMap
	if sch == nil {
		fmt.Println("status code header map not configured")
		return
	}

	// setter function to set the headers
	setter := func(hvl []*CustomHeader) {
		for _, hv := range hvl {
			w.Header().Set(hv.Name, hv.Value)
		}
	}

	// Checking the validity of the status code
	if status >= 600 || status < 100 {
		fmt.Println("invalid status code")
		return
	}

	// Setting the default headers first
	setter(sch["default"])

	// setting the Xyy pattern first
	d := fmt.Sprintf("%vxx", status / 100)
	if val, ok := sch[d]; ok {
		setter(val)
	}

	// Setting the specific headers
	if val, ok := sch[strconv.Itoa(status)]; ok {
		setter(val)
	}

	return
}


func (l *ListenerCustomHeaders) findCustomHeaderMatchStatusCode(sc string, hn string) string {

	getHeader := func(ch []*CustomHeader) string {
		for _, h := range ch {
			if h.Name == hn {
				return h.Value
			}
		}
		return ""
	}

	hm := l.StatusCodeHeaderMap

	// starting with the most specific status code
	if ch, ok := hm[sc]; ok {
		h := getHeader(ch)
		if h != "" {
			return h
		}
	}

	// Checking for the Yxx pattern
	var firstDig string
	if len(sc) == 3 {
		firstDig = strings.Split(sc, "")[0]
	}
	if firstDig != "" {
		s := fmt.Sprintf("%vxx", firstDig)
		if configutil.IsValidStatusCodeCollection(s) {
			if ch, ok := hm[s]; ok {
				h := getHeader(ch)
				if h != "" {
					return h
				}
			}
		}
	}

	// At this point, we could not find a match for the given status code in the config file
	// so, we just return the "default" ones
	h := getHeader(hm["default"])
	if h != ""{
		return h
	}

	return ""
}

func(l *ListenerCustomHeaders) FetchHeaderForStatusCode(header, sc string) (string, error) {

	if header == "" {
		return "", fmt.Errorf("invalid target header")
	}

	if l.StatusCodeHeaderMap == nil {
		return "", fmt.Errorf("custom headers not configured")
	}

	if !configutil.IsValidStatusCode(sc) {
		return "", fmt.Errorf("failed to check if a header exist in config file due to invalid status code")
	}

	hn := textproto.CanonicalMIMEHeaderKey(header)

	h := l.findCustomHeaderMatchStatusCode(sc, hn)

	return h, nil
}

func (l *ListenerCustomHeaders) ExistCustomResponseHeader(header string) bool {

	if header == "" {
		return false
	}

	if l.StatusCodeHeaderMap == nil {
		return false
	}

	hn := textproto.CanonicalMIMEHeaderKey(header)

	hs := l.ConfiguredHeadersStatusCodeMap
	_, ok := hs[hn]
	return ok
}