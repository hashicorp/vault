package vault

import (
	"net/url"
	"regexp"
	"strings"
)

func isValidRedirectURI(uri string, validUris []string) bool {
	requestedUri, err := url.Parse(uri)
	if err != nil {
		return false
	}

	for _, validUri := range validUris {
		if strings.ToLower(validUri) == strings.ToLower(uri) || isLoopbackURI(requestedUri, validUri) {
			return true
		}
	}

	return false
}

func isLoopbackURI(requestUri *url.URL, validUri string) bool {
	registeredUri, err := url.Parse(validUri)
	if err != nil {
		return false
	}

	// Verifies that the valid URL is HTTP and is the loopback address before
	// proceeding, otherwise return false
	if registeredUri.Scheme != "http" || !isLoopbackAddress(registeredUri.Host) {
		return false
	}

	// Returns true if the path after the IP/port is the same
	// Request URL and valid URL have already been validated as loopback
	if requestUri.Scheme == "http" && isLoopbackAddress(requestUri.Host) && registeredUri.Path == requestUri.Path {
		return true
	}

	return false
}

// Returns true if the address hostname is the IPv4 or IPv6 loopback address and ignores port
func isLoopbackAddress(address string) bool {
	match, _ := regexp.MatchString("^(127.0.0.1|\\[::1\\])(:?)(\\d*)$", address)
	return match
}
