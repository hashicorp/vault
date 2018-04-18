package util

import "errors"

func ParseRoleName(prefix string, reqPath string) (string, error) {
	prefixLen := len(prefix)
	if len(reqPath) <= prefixLen {
		return "", errors.New("role name must be provided")
	}
	return reqPath[prefixLen:], nil
}
