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

// DefaultCustomResponseStatus is used to set default headers early before having a status code,
// for example, for /ui headers
const DefaultCustomResponseStatus = 1

type ListenersCustomHeaderList struct {
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

func NewListenerCustomHeader(ln []*configutil.Listener, logger log.Logger, uiHeaders http.Header) *ListenersCustomHeaderList {

	if ln == nil {
		return nil
	}

	ll := &ListenersCustomHeaderList{
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

				// X-Vault- prefix is reserved for Vault internal processes
				if strings.HasPrefix(h, "X-Vault-") {
					logger.Error("Custom headers starting with X-Vault are not valid", "header", h)
					continue
				}

				// Checking for UI headers, if any common header exist, HCL headers take precedence
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

func (c *ListenersCustomHeaderList) SetCustomResponseHeaders(w http.ResponseWriter, status int) {
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
    	c.logger.Warn("no listener config found")
    	return
	}

	// setter function to set the headers
	setter := func(hvl []*CustomHeader) {
		for _, hv := range hvl {
			w.Header().Set(hv.Name, hv.Value)
		}
	}

	// Checking the validity of the status code
	if status >= 600 || (status < 100 && status != DefaultCustomResponseStatus) {
		c.logger.Error("invalid status code")
		return
	}

	// Setting the default headers first
	setter(lch["default"])

	// for DefaultCustomResponseStatus, we only set the default headers
	if status == DefaultCustomResponseStatus {
		return
	}

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

func (c *ListenersCustomHeaderList) getListenerMap(address string) map[string][]*CustomHeader {
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

func (c *ListenersCustomHeaderList) findCustomHeaderMatchStatusCode(hm map[string][]*CustomHeader, sc int) ([]*CustomHeader, error) {

	if sc == DefaultCustomResponseStatus {
		return hm["default"], nil
	}

	if h, ok := hm[strconv.Itoa(sc)]; ok {
		return h, nil
	}

	d := fmt.Sprintf("%vxx", sc / 100)
	for _, s := range configutil.ValidCustomStatusCodeCollection {
		if s == d {
			if h, ok := hm[s]; ok {
				return h, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to find a match for the given status code:%v", sc)
}

func (c *ListenersCustomHeaderList) FetchCustomResponseHeaderValue(header string, sc int, la string) ([]string, error) {

	if header == "" {
		return nil, fmt.Errorf("invalid target header")
	}

	getHeader := func(hm map[string][]*CustomHeader) (string, error){
		ch, err := c.findCustomHeaderMatchStatusCode(hm, sc)
		if err != nil {
			return "", err
		}

		if ch == nil {
			return "", nil
		}

		hn := textproto.CanonicalMIMEHeaderKey(header)
		for _, h := range ch {
			if h.Name == hn {
				return h.Value, nil
			}
		}

		return "", nil
	}

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
	for _, l := range lch {
		h, err := getHeader(l.StatusCodeHeaderMap)
		if err != nil || h == "" {
			continue
		}
		headers = append(headers, h)
	}

	return headers, err
}

func(c *ListenersCustomHeaderList) ExistHeader(th string, sl []int, la string) bool {
	if len(sl) == 0 {
		return false
	}

	for _, s := range sl {
		chv, _ := c.FetchCustomResponseHeaderValue(th, s, la)
		if chv != nil {
			return true
		}
	}

	return false
}
