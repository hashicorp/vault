package gosnowflake

import (
	"net/url"
	"regexp"
)

var (
	matcher, _ = regexp.Compile(`^http(s?)\:\/\/[0-9a-zA-Z]([-.\w]*[0-9a-zA-Z@:])*(:(0-9)*)*(\/?)([a-zA-Z0-9\-\.\?\,\&\(\)\/\\\+&%\$#_=@]*)?$`)
)

func isValidURL(targetURL string) bool {
	if !matcher.MatchString(targetURL) {
		logger.Infof(" The provided URL is not a valid URL - " + targetURL)
		return false
	}
	return true
}

func urlEncode(targetString string) string {
	// We use QueryEscape instead of PathEscape here
	// for consistency across Drivers. For example:
	// QueryEscape escapes space as "+" whereas PE
	// it as %20F. PE also does not escape @ or &
	// either but QE does.
	// The behavior of QE in Golang is more in sync
	// with URL encoders in Python and Java hence the choice
	return url.QueryEscape(targetString)
}
