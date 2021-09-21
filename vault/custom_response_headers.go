package vault

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
)

type ListenerCustomHeaders struct {
	Address             string
	StatusCodeHeaderMap map[string][]*CustomHeader
	// ConfiguredHeadersStatusCodeMap field is introduced so that we would not need to loop through
	// StatusCodeHeaderMap to see if a header exists, the key for this map is the headers names
	configuredHeadersStatusCodeMap map[string][]string
}

type CustomHeader struct {
	Name  string
	Value string
}

func NewListenerCustomHeader(ln []*configutil.Listener, logger log.Logger, uiHeaders http.Header) []*ListenerCustomHeaders {
	var ll []*ListenerCustomHeaders

	for _, l := range ln {
		lc := &ListenerCustomHeaders{
			Address: l.Address,
		}
		lc.StatusCodeHeaderMap = make(map[string][]*CustomHeader)
		lc.configuredHeadersStatusCodeMap = make(map[string][]string)
		for sc, hv := range l.CustomResponseHeaders {
			var chl []*CustomHeader
			for h, v := range hv {
				// Sanitizing custom headers
				// X-Vault- prefix is reserved for Vault internal processes
				if strings.HasPrefix(h, "X-Vault-") {
					logger.Warn("custom headers starting with X-Vault are not valid", "header", h)
					continue
				}

				// Checking for UI headers, if any common header exists, we just log an error
				if uiHeaders != nil {
					exist := uiHeaders.Get(h)
					if exist != "" {
						logger.Warn("found a duplicate header in UI", "header:", h, "Headers defined in the server configuration take precedence.")
					}
				}

				// Checking if the header value is not an empty string
				if v == "" {
					logger.Warn("header value is an empty string", "header", h, "value", v)
					continue
				}

				ch := &CustomHeader{
					Name:  h,
					Value: v,
				}

				chl = append(chl, ch)

				// setting up the reverse map of header to status code for easy lookups
				lc.configuredHeadersStatusCodeMap[h] = append(lc.configuredHeadersStatusCodeMap[h], sc)
			}
			lc.StatusCodeHeaderMap[sc] = chl
		}
		ll = append(ll, lc)
	}

	return ll
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
		firstDig = string(sc[0])
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
	if h != "" {
		return h
	}

	return ""
}

func (l *ListenerCustomHeaders) FetchHeaderForStatusCode(header, sc string) (string, error) {
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

	hs := l.configuredHeadersStatusCodeMap
	_, ok := hs[hn]
	return ok
}
