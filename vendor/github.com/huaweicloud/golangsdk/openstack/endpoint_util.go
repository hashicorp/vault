package openstack

import "regexp"

// A regular expression used to verify whether or not contains a project id in an endpoint url
var endpointProjectIdMatcher = regexp.MustCompile(`http[s]?://.+/(?:[V|v]\d+|[V|v]\d+\.\d+)/([a-z|A-Z|0-9]{32})(?:/|$)`)

// ContainsProjectId detects whether or not  the encpoint url contains a projectid
func ContainsProjectId(endpointUrl string) bool {
	return endpointProjectIdMatcher.Match([]byte(endpointUrl))
}
