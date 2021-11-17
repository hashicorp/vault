package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//
// GET /api/vhosts
//

// Example response:

// [
//   {
//     "message_stats": {
//       "publish": 78,
//       "publish_details": {
//         "rate": 0
//       }
//     },
//     "messages": 0,
//     "messages_details": {
//       "rate": 0
//     },
//     "messages_ready": 0,
//     "messages_ready_details": {
//       "rate": 0
//     },
//     "messages_unacknowledged": 0,
//     "messages_unacknowledged_details": {
//       "rate": 0
//     },
//     "recv_oct": 16653,
//     "recv_oct_details": {
//       "rate": 0
//     },
//     "send_oct": 40495,
//     "send_oct_details": {
//       "rate": 0
//     },
//     "name": "\/",
//	   "description": "myvhost",
//     "tags": "production,eu-west-1",
//     "tracing": false
//   },
//   {
//     "name": "29dd51888b834698a8b5bc3e7f8623aa1c9671f5",
//     "tracing": false
//   }
// ]

// VhostInfo represents a virtual host, its properties and key metrics.
type VhostInfo struct {
	// Virtual host name
	Name string `json:"name"`
	// Virtual host description
	Description string `json:"description"`
	// Virtual host tags
	Tags VhostTags `json:"tags"`
	// True if tracing is enabled for this virtual host
	Tracing bool `json:"tracing"`

	// Total number of messages in queues of this virtual host
	Messages        int         `json:"messages"`
	MessagesDetails RateDetails `json:"messages_details"`

	// Total number of messages ready to be delivered in queues of this virtual host
	MessagesReady        int         `json:"messages_ready"`
	MessagesReadyDetails RateDetails `json:"messages_ready_details"`

	// Total number of messages pending acknowledgement from consumers in this virtual host
	MessagesUnacknowledged        int         `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails RateDetails `json:"messages_unacknowledged_details"`

	// Octets received
	RecvOct uint64 `json:"recv_oct"`
	// Octets sent
	SendOct        uint64      `json:"send_oct"`
	RecvCount      uint64      `json:"recv_cnt"`
	SendCount      uint64      `json:"send_cnt"`
	SendPending    uint64      `json:"send_pend"`
	RecvOctDetails RateDetails `json:"recv_oct_details"`
	SendOctDetails RateDetails `json:"send_oct_details"`

	// Cluster State
	ClusterState map[string]string `json:"cluster_state"`
}

type VhostTags []string

// MarshalJSON can marshal an array of strings or a comma-separated list in a string
func (d VhostTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(d, ","))
}

// UnmarshalJSON can unmarshal an array of strings or a comma-separated list in a string
func (d *VhostTags) UnmarshalJSON(b []byte) error {
	// the value is a comma-separated string
	t, _ := strconv.Unquote(string(b))
	if b[0] == '"' {
		quotedTags := strings.Split(t, ",")
		var tags []string
		for _, qt := range quotedTags {
			tags = append(tags, qt)
		}
		*d = tags
		return nil
	}
	// the value is an array
	var ary []string
	if err := json.Unmarshal(b, &ary); err != nil {
		return err
	}
	*d = ary
	return nil
}

// ListVhosts returns a list of virtual hosts.
func (c *Client) ListVhosts() (rec []VhostInfo, err error) {
	req, err := newGETRequest(c, "vhosts")
	if err != nil {
		return []VhostInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []VhostInfo{}, err
	}

	return rec, nil
}

//
// GET /api/vhosts/{name}
//

// GetVhost returns information about a specific virtual host.
func (c *Client) GetVhost(vhostname string) (rec *VhostInfo, err error) {
	req, err := newGETRequest(c, "vhosts/"+url.PathEscape(vhostname))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/vhosts/{name}
//

// VhostSettings are properties used to create or modify virtual hosts.
type VhostSettings struct {
	// Virtual host description
	Description string `json:"description"`
	// Virtual host tags
	Tags VhostTags `json:"tags"`
	// True if tracing should be enabled.
	Tracing bool `json:"tracing"`
}

// PutVhost creates or updates a virtual host.
func (c *Client) PutVhost(vhostname string, settings VhostSettings) (res *http.Response, err error) {
	body, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "vhosts/"+url.PathEscape(vhostname), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/vhosts/{name}
//

// DeleteVhost deletes a virtual host.
func (c *Client) DeleteVhost(vhostname string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "vhosts/"+url.PathEscape(vhostname), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
