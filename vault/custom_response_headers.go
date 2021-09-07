package vault

import (
	"fmt"
	log "github.com/hashicorp/go-hclog"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/internalshared/configutil"
)

type ListenersCustomResponseHeadersList struct {
	logger log.Logger
	CustomHeadersList []*ListenerCustomHeaders
}

type ListenerCustomHeaders struct {
	Address string
	StatusCodeHeaderMap map[string][]*CustomHeader
}

type CustomHeader struct {
	Name string
	Value string
}

func NewListenerCustomHeader(ln []*configutil.Listener, logger log.Logger, uiHeaders http.Header) *ListenersCustomResponseHeadersList {

	if ln == nil {
		return nil
	}

	ll := &ListenersCustomResponseHeadersList{
		logger: logger,
	}

	for _, l := range ln {
		lc := &ListenerCustomHeaders{
			Address: l.Address,
		}
		lc.StatusCodeHeaderMap = make(map[string][]*CustomHeader)
		for sc, hv := range l.CustomResponseHeaders {
			var chl []*CustomHeader
			for h, v := range hv {
				// Sanitizing custom headers
				// X-Vault- prefix is reserved for Vault internal processes
				if strings.HasPrefix(h, "X-Vault-") {
					logger.Error("Custom headers starting with X-Vault are not valid", "header", h)
					continue
				}

				// Checking for UI headers, if any common header exists, we just log an error
				if uiHeaders != nil {
					exist := uiHeaders.Get(h)
					if exist != "" {
						logger.Error("found a duplicate header in UI, note that config file headers take precedence.", "header:", h)
					}
				}

				ch := &CustomHeader{
					Name: h,
					Value: v,
				}

				chl = append(chl, ch)
			}
			lc.StatusCodeHeaderMap[sc] = chl
		}
		ll.CustomHeadersList = append(ll.CustomHeadersList, lc)
	}

	return ll
}

func (c *ListenersCustomResponseHeadersList) SetCustomResponseHeaders(w http.ResponseWriter, status int) {
	if w == nil {
		c.logger.Error("No ResponseWriter provided")
	}

	// Getting the listener address to set its corresponding custom headers
	la := w.Header().Get("X-Vault-Listener-Add")
	if la == "" {
		c.logger.Error("X-Vault-Listener-Add was not set in the ResponseWriter")
		return
	}

	// Removing X-Vault-Listener-Add header from ResponseWriter
	// This should be safe as the call to this function is right
	// before w.WriteHeader for which the status code is finalized and known
	w.Header().Del("X-Vault-Listener-Add")

	lch := c.getListenerMap(la)
    if lch == nil {
        c.logger.Warn("no listener config found", "address", la)
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
		c.logger.Error("invalid status code")
		return
	}

	// Setting the default headers first
	setter(lch["default"])

	// setting the Xyy pattern first
	d := fmt.Sprintf("%vxx", status / 100)
	if val, ok := lch[d]; ok {
		setter(val)
	}
	// Setting the specific headers
	if val, ok := lch[strconv.Itoa(status)]; ok {
		setter(val)
	}

	return
}

func (c *ListenersCustomResponseHeadersList) getListenerMap(address string) map[string][]*CustomHeader {
	if c.CustomHeadersList == nil {
		return nil
	}
	for _, l := range c.CustomHeadersList {
		if l.Address == address {
			return l.StatusCodeHeaderMap
		}
	}
	return nil
}

func (c *ListenersCustomResponseHeadersList) findCustomHeaderMatchStatusCode(hm map[string][]*CustomHeader, sc int, hn string) string {

	getHeader := func(ch []*CustomHeader) string {
		for _, h := range ch {
			if h.Name == hn {
				return h.Value
			}
		}
		return ""
	}

	// starting with the most specific status code
	if ch, ok := hm[strconv.Itoa(sc)]; ok {
		h := getHeader(ch)
		if h != "" {
			return h
		}
	}

	s := fmt.Sprintf("%vxx", sc/100)
	if configutil.IsValidStatusCodeCollection(s) {
		if ch, ok := hm[s]; ok {
			h := getHeader(ch)
			if h != "" {
				return h
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

func (c *ListenersCustomResponseHeadersList) FetchCustomResponseHeaderValue(header string, sc int, la string) ([]string, error) {

	if header == "" {
		return nil, fmt.Errorf("invalid target header")
	}

	// either looking for a specific listener, or if listener address isn't given,
	// checking for all available listeners
	var lch []*ListenerCustomHeaders
	if la == "" {
		lch = c.CustomHeadersList
	} else {
		for _, l := range c.CustomHeadersList {
			if l.Address == la {
				lch = append(lch, l)
			}
		}
		if len(lch) == 0 {
			return nil, fmt.Errorf("no listener found with address:%v", la)
		}
	}

	var headers []string
	var err error
	hn := textproto.CanonicalMIMEHeaderKey(header)
	for _, l := range lch {
		h := c.findCustomHeaderMatchStatusCode(l.StatusCodeHeaderMap, sc, hn)
		if h == "" {
			continue
		}
		headers = append(headers, h)
	}

	return headers, err
}

func(c *ListenersCustomResponseHeadersList) ExistHeader(th string, sc int, la string) bool {
	if !configutil.IsValidStatusCode(strconv.Itoa(sc)) {
		c.logger.Error("failed to check if a header exist in config file due to invalid status code")
		return false
	}

	chv, _ := c.FetchCustomResponseHeaderValue(th, sc, la)
	if chv != nil {
		return true
	}

	return false
}
