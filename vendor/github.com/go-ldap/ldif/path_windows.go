package ldif

import "net/url"
import "strings"

// toPath get the file path
// We use ioutil.ReadFile to read the content file.
// On windows,
// https://github.com/golang/go/blob/95a11c7381e01fdaaf34e25b82db0632081ab74e/src/net/url/url_test.go#L283-L292
func toPath(u *url.URL) string {
	return strings.TrimPrefix(u.Path, "/")
}
