// +build go1.7

package oss

import "net/http"

// http.ErrUseLastResponse only is defined go1.7 onward
func disableHTTPRedirect(client *http.Client) {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
