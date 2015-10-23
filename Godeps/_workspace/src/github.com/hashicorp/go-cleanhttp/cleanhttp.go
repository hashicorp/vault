package cleanhttp

import (
	"net"
	"net/http"
	"time"
)

// DefaultTransport returns a new http.Transport with the same default values
// as http.DefaultTransport
func DefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}

// DefaultClient returns a new http.Client with the same default values as
// http.Client, but with a non-shared Transport
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: DefaultTransport(),
	}
}
