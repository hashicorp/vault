// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type ListenerCustomHeaders struct {
	Address             string
	StatusCodeHeaderMap map[string][]*logical.CustomHeader
	// ConfiguredHeadersStatusCodeMap field is introduced so that we would not need to loop through
	// StatusCodeHeaderMap to see if a header exists, the key for this map is the headers names
	configuredHeadersStatusCodeMap map[string][]string
}

func NewListenerCustomHeader(ln []*configutil.Listener, logger log.Logger, uiHeaders http.Header) []*ListenerCustomHeaders {
	var listenerCustomHeadersList []*ListenerCustomHeaders

	for _, l := range ln {
		listenerCustomHeaderStruct := &ListenerCustomHeaders{
			Address: l.Address,
		}
		listenerCustomHeaderStruct.StatusCodeHeaderMap = make(map[string][]*logical.CustomHeader)
		listenerCustomHeaderStruct.configuredHeadersStatusCodeMap = make(map[string][]string)
		for statusCode, headerValMap := range l.CustomResponseHeaders {
			var customHeaderList []*logical.CustomHeader
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
						logger.Warn(fmt.Sprintf("found a duplicate header in UI: header=%s. Headers defined in the server configuration take precedence.", headerName))
					}
				}

				// Checking if the header value is not an empty string
				if headerVal == "" {
					logger.Warn("header value is an empty string", "header", headerName, "value", headerVal)
					continue
				}

				ch := &logical.CustomHeader{
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
