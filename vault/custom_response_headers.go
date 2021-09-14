package vault

import (
	log "github.com/hashicorp/go-hclog"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
)

type ListenersCustomResponseHeadersList struct {
	CustomHeadersList []*listenerutil.ListenerCustomHeaders
}

func NewListenerCustomHeader(ln []*configutil.Listener, logger log.Logger, uiHeaders http.Header) *ListenersCustomResponseHeadersList {

	if ln == nil {
		return nil
	}

	ll := &ListenersCustomResponseHeadersList{}

	for _, l := range ln {
		lc := &listenerutil.ListenerCustomHeaders{
			Address: l.Address,
		}
		lc.StatusCodeHeaderMap = make(map[string][]*listenerutil.CustomHeader)
		lc.ConfiguredHeadersStatusCodeMap = make(map[string][]string)
		for sc, hv := range l.CustomResponseHeaders {
			var chl []*listenerutil.CustomHeader
			for h, v := range hv {
				// Sanitizing custom headers
				// X-Vault- prefix is reserved for Vault internal processes
				if strings.HasPrefix(h, "X-Vault-") {
					logger.Warn("Custom headers starting with X-Vault are not valid", "header", h)
					continue
				}

				// Checking for UI headers, if any common header exists, we just log an error
				if uiHeaders != nil {
					exist := uiHeaders.Get(h)
					if exist != "" {
						logger.Warn("found a duplicate header in UI, note that config file headers take precedence.", "header:", h)
					}
				}

				// Checking if the header value is not an empty string
				if v == "" {
					logger.Warn("header value is an empty string", "header", h, "value", v)
					continue
				}

				ch := &listenerutil.CustomHeader{
					Name: h,
					Value: v,
				}

				chl = append(chl, ch)

				// setting up the reverse map of header to status code for easy lookups
				lc.ConfiguredHeadersStatusCodeMap[h] = append(lc.ConfiguredHeadersStatusCodeMap[h], sc)
			}
			lc.StatusCodeHeaderMap[sc] = chl
		}
		ll.CustomHeadersList = append(ll.CustomHeadersList, lc)
	}

	return ll
}

func (c *ListenersCustomResponseHeadersList) getListenerMap(address string) []*listenerutil.ListenerCustomHeaders {
	if c.CustomHeadersList == nil {
		return nil
	}

	// either looking for a specific listener, or if listener address isn't given,
	// checking for all available listeners
	var lch []*listenerutil.ListenerCustomHeaders
	if address == "" {
		lch = c.CustomHeadersList
	} else {
		for _, l := range c.CustomHeadersList {
			if l.Address == address {
				lch = append(lch, l)
			}
		}
		if len(lch) == 0 {
			return nil
		}
	}
	return lch
}
