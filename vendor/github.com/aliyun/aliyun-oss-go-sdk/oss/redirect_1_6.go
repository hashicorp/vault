//go:build !go1.7
// +build !go1.7

package oss

import "net/http"

// http.ErrUseLastResponse only is defined go1.7 onward

func disableHTTPRedirect(client *http.Client) {

}
