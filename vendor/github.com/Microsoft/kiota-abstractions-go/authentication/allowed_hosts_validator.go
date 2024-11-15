package authentication

import (
	"errors"
	u "net/url"
	"strings"
)

// AllowedHostsValidator maintains a list of valid hosts and allows authentication providers to check whether a host is valid before authenticating a request
type AllowedHostsValidator struct {
	validHosts map[string]bool
}

// ErrInvalidHostPrefix indicates that a host should not contain the http or https prefix.
var ErrInvalidHostPrefix = errors.New("host should not contain http or https prefix")

// Deprecated: NewAllowedHostsValidator creates a new AllowedHostsValidator object with provided values.
func NewAllowedHostsValidator(validHosts []string) AllowedHostsValidator {
	result := AllowedHostsValidator{}
	result.SetAllowedHosts(validHosts)
	return result
}

// NewAllowedHostsValidatorErrorCheck creates a new AllowedHostsValidator object with provided values and performs error checking.
func NewAllowedHostsValidatorErrorCheck(validHosts []string) (*AllowedHostsValidator, error) {
	result := &AllowedHostsValidator{}
	if err := result.SetAllowedHostsErrorCheck(validHosts); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllowedHosts returns the list of valid hosts.
func (v *AllowedHostsValidator) GetAllowedHosts() map[string]bool {
	hosts := make(map[string]bool, len(v.validHosts))
	for host := range v.validHosts {
		hosts[host] = true
	}
	return hosts
}

// Deprecated: SetAllowedHosts sets the list of valid hosts.
func (v *AllowedHostsValidator) SetAllowedHosts(hosts []string) {
	v.validHosts = make(map[string]bool, len(hosts))
	if len(hosts) > 0 {
		for _, host := range hosts {
			v.validHosts[strings.ToLower(host)] = true
		}
	}
}

// SetAllowedHostsErrorCheck sets the list of valid hosts with error checking.
func (v *AllowedHostsValidator) SetAllowedHostsErrorCheck(hosts []string) error {
	v.validHosts = make(map[string]bool, len(hosts))
	if len(hosts) > 0 {
		for _, host := range hosts {
			lowerHost := strings.ToLower(host)
			if strings.HasPrefix(lowerHost, "http://") || strings.HasPrefix(lowerHost, "https://") {
				return ErrInvalidHostPrefix
			}
			v.validHosts[lowerHost] = true
		}
	}
	return nil
}

// IsValidHost returns true if the host is valid.
func (v *AllowedHostsValidator) IsUrlHostValid(uri *u.URL) bool {
	if uri == nil {
		return false
	}
	host := uri.Hostname()
	if host == "" {
		return false
	}
	return len(v.validHosts) == 0 || v.validHosts[strings.ToLower(host)]
}
