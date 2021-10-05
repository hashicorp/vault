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
	var listenerCustomHeadersList []*ListenerCustomHeaders

	for _, l := range ln {
		listenerCustomHeaderStruct := &ListenerCustomHeaders{
			Address: l.Address,
		}
		listenerCustomHeaderStruct.StatusCodeHeaderMap = make(map[string][]*CustomHeader)
		listenerCustomHeaderStruct.configuredHeadersStatusCodeMap = make(map[string][]string)
		for statusCode, headerValMap := range l.CustomResponseHeaders {
			var customHeaderList []*CustomHeader
			for headerName, headerVal := range headerValMap {
				// Sanitizing custom headers
				// X-Vault- prefix is reserved for Vault internal processes
				if strings.HasPrefix(headerName, "X-Vault-") {
					logger.Warn("custom headers starting with X-Vault are not valid", "header", headerName)
					continue
				}

				// Checking for UI headers, if any common header exists, we just log an error
				if uiHeaders != nil {
					exist := uiHeaders.Get(headerName)
					if exist != "" {
						logger.Warn("found a duplicate header in UI", "header:", headerName, "Headers defined in the server configuration take precedence.")
					}
				}

				// Checking if the header value is not an empty string
				if headerVal == "" {
					logger.Warn("header value is an empty string", "header", headerName, "value", headerVal)
					continue
				}

				ch := &CustomHeader{
					Name:  headerName,
					Value: headerVal,
				}

				customHeaderList = append(customHeaderList, ch)

				// setting up the reverse map of header to status code for easy lookups
				listenerCustomHeaderStruct.configuredHeadersStatusCodeMap[headerName] = append(listenerCustomHeaderStruct.configuredHeadersStatusCodeMap[headerName], statusCode)
			}
			listenerCustomHeaderStruct.StatusCodeHeaderMap[statusCode] = customHeaderList
		}
		listenerCustomHeadersList = append(listenerCustomHeadersList, listenerCustomHeaderStruct)
	}

	return listenerCustomHeadersList
}

func (l *ListenerCustomHeaders) findCustomHeaderMatchStatusCode(statusCode string, headerName string) string {
	getHeader := func(ch []*CustomHeader) string {
		for _, h := range ch {
			if h.Name == headerName {
				return h.Value
			}
		}
		return ""
	}

	headerMap := l.StatusCodeHeaderMap

	// starting with the most specific status code
	if customHeaderList, ok := headerMap[statusCode]; ok {
		h := getHeader(customHeaderList)
		if h != "" {
			return h
		}
	}

	// Checking for the Yxx pattern
	var firstDig string
	if len(statusCode) == 3 {
		firstDig = string(statusCode[0])
	}
	if firstDig != "" {
		s := fmt.Sprintf("%vxx", firstDig)
		if configutil.IsValidStatusCodeCollection(s) {
			if customHeaderList, ok := headerMap[s]; ok {
				h := getHeader(customHeaderList)
				if h != "" {
					return h
				}
			}
		}
	}

	// At this point, we could not find a match for the given status code in the config file
	// so, we just return the "default" ones
	h := getHeader(headerMap["default"])
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

	headerName := textproto.CanonicalMIMEHeaderKey(header)

	h := l.findCustomHeaderMatchStatusCode(sc, headerName)

	return h, nil
}

func (l *ListenerCustomHeaders) ExistCustomResponseHeader(header string) bool {
	if header == "" {
		return false
	}

	if l.StatusCodeHeaderMap == nil {
		return false
	}

	headerName := textproto.CanonicalMIMEHeaderKey(header)

	headerMap := l.configuredHeadersStatusCodeMap
	_, ok := headerMap[headerName]
	return ok
}
