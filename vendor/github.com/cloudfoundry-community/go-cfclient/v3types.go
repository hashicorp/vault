package cfclient

// Pagination is used by the V3 apis
type Pagination struct {
	TotalResults int  `json:"total_results"`
	TotalPages   int  `json:"total_pages"`
	First        Link `json:"first"`
	Last         Link `json:"last"`
	Next         Link `json:"next"`
	Previous     Link `json:"previous"`
}

// Link is a HATEOAS-style link for v3 apis
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

// V3ToOneRelationship is a relationship to a single object
type V3ToOneRelationship struct {
	Data V3Relationship `json:"data,omitempty"`
}

// V3ToManyRelationships is a relationship to multiple objects
type V3ToManyRelationships struct {
	Data []V3Relationship `json:"data,omitempty"`
}

type V3Relationship struct {
	GUID string `json:"guid,omitempty"`
}

type V3Metadata struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
