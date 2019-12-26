package cfclient

// Pagination is used by the V3 apis
type Pagination struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	First        Link        `json:"first"`
	Last         Link        `json:"last"`
	Next         interface{} `json:"next"`
	Previous     interface{} `json:"previous"`
}

// Link is a HATEOAS-style link for v3 apis
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}
