package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// ShovelInfo contains the configuration of a dynamic Shovel
type ShovelInfo struct {
	// Shovel name
	Name string `json:"name"`
	// Virtual host this Shovel belongs to
	Vhost string `json:"vhost"`
	// Runtime component of the Shovel
	Component string `json:"component"`
	// Details the configuration values of the Shovel
	Definition ShovelDefinition `json:"value"`
}

// ShovelStatus represents the status of a Shovel
type ShovelStatus struct {
	// Shovel name
	Name string `json:"name"`
	// Virtual host this shovel belongs to
	Vhost string `json:"vhost"`
	// Type of this shovel
	Type string `json:"type"`
	// State of this shovel
	State string `json:"state"`
	// Timestamp is the moment when this Shovel last reported its state change (e.g. was started)
	Timestamp string `json:"timestamp"`
}

// DeleteAfter after can hold a delete-after value which may be a string (eg. "never") or an integer
type DeleteAfter string

// MarshalJSON can marshal a string or an integer
func (d DeleteAfter) MarshalJSON() ([]byte, error) {
	deleteAfterInt, err := strconv.Atoi(string(d))
	if err != nil {
		return json.Marshal(string(d))
	}
	return json.Marshal(deleteAfterInt)
}

// UnmarshalJSON can unmarshal a string or an integer
func (d *DeleteAfter) UnmarshalJSON(b []byte) error {
	// delete-after is a string, such as "never"
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*d = DeleteAfter(s)
		return nil
	}

	// delete-after is a number
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*d = DeleteAfter(strconv.Itoa(i))
	return nil
}

// ShovelDefinition contains the details of the shovel configuration
type ShovelDefinition struct {
	DestinationURI URISet `json:"dest-uri"`
	SourceURI      URISet `json:"src-uri"`

	AckMode                          string                 `json:"ack-mode,omitempty"`
	AddForwardHeaders                bool                   `json:"add-forward-headers,omitempty"`
	DeleteAfter                      DeleteAfter            `json:"delete-after,omitempty"`
	DestinationAddForwardHeaders     bool                   `json:"dest-add-forward-headers,omitempty"`
	DestinationAddTimestampHeader    bool                   `json:"dest-add-timestamp-header,omitempty"`
	DestinationAddress               string                 `json:"dest-address,omitempty"`
	DestinationApplicationProperties map[string]interface{} `json:"dest-application-properties,omitempty"`
	DestinationExchange              string                 `json:"dest-exchange,omitempty"`
	DestinationExchangeKey           string                 `json:"dest-exchange-key,omitempty"`
	DestinationProperties            map[string]interface{} `json:"dest-properties,omitempty"`
	DestinationProtocol              string                 `json:"dest-protocol,omitempty"`
	DestinationPublishProperties     map[string]interface{} `json:"dest-publish-properties,omitempty"`
	DestinationQueue                 string                 `json:"dest-queue,omitempty"`
	DestinationQueueArgs             map[string]interface{} `json:"dest-queue-args,omitempty"`
	DestinationMessageAnnotations    map[string]interface{} `json:"dest-message-annotations,omitempty"`
	PrefetchCount                    int                    `json:"prefetch-count,omitempty"`
	ReconnectDelay                   int                    `json:"reconnect-delay,omitempty"`
	SourceAddress                    string                 `json:"src-address,omitempty"`
	SourceDeleteAfter                DeleteAfter            `json:"src-delete-after,omitempty"`
	SourceExchange                   string                 `json:"src-exchange,omitempty"`
	SourceExchangeKey                string                 `json:"src-exchange-key,omitempty"`
	SourcePrefetchCount              int                    `json:"src-prefetch-count,omitempty"`
	SourceProtocol                   string                 `json:"src-protocol,omitempty"`
	SourceQueue                      string                 `json:"src-queue,omitempty"`
	SourceQueueArgs                  map[string]interface{} `json:"src-queue-args,omitempty"`
	SourceConsumerArgs               map[string]interface{} `json:"src-consumer-args,omitempty"`
}

// ShovelDefinitionDTO provides a data transfer object
type ShovelDefinitionDTO struct {
	Definition ShovelDefinition `json:"value"`
}

//
// GET /api/parameters/shovel
//

// ListShovels returns all shovels
func (c *Client) ListShovels() (rec []ShovelInfo, err error) {
	req, err := newGETRequest(c, "parameters/shovel")
	if err != nil {
		return []ShovelInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ShovelInfo{}, err
	}

	return rec, nil
}

//
// GET /api/parameters/shovel/{vhost}
//

// ListShovelsIn returns all shovels in a vhost
func (c *Client) ListShovelsIn(vhost string) (rec []ShovelInfo, err error) {
	req, err := newGETRequest(c, "parameters/shovel/"+url.PathEscape(vhost))
	if err != nil {
		return []ShovelInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ShovelInfo{}, err
	}

	return rec, nil
}

//
// GET /api/parameters/shovel/{vhost}/{name}
//

// GetShovel returns a shovel configuration
func (c *Client) GetShovel(vhost, shovel string) (rec *ShovelInfo, err error) {
	req, err := newGETRequest(c, "parameters/shovel/"+url.PathEscape(vhost)+"/"+url.PathEscape(shovel))

	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/parameters/shovel/{vhost}/{name}
//

// DeclareShovel creates a shovel
func (c *Client) DeclareShovel(vhost, shovel string, info ShovelDefinition) (res *http.Response, err error) {
	shovelDTO := ShovelDefinitionDTO{Definition: info}

	body, err := json.Marshal(shovelDTO)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "parameters/shovel/"+url.PathEscape(vhost)+"/"+url.PathEscape(shovel), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/parameters/shovel/{vhost}/{name}
//

// DeleteShovel a shovel
func (c *Client) DeleteShovel(vhost, shovel string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "parameters/shovel/"+url.PathEscape(vhost)+"/"+url.PathEscape(shovel), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// GET /api/shovels/{vhost}
//

// ListShovelStatus returns status of all shovels in a vhost
func (c *Client) ListShovelStatus(vhost string) (rec []ShovelStatus, err error) {
	req, err := newGETRequest(c, "shovels/"+url.PathEscape(vhost))
	if err != nil {
		return []ShovelStatus{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ShovelStatus{}, err
	}

	return rec, nil
}
