// +build !windows

package ldif

import "net/url"

func toPath(u *url.URL) string {
	return u.Path
}
