package rabbithole

import (
	"encoding/json"
	"net/http"
)

// ShovelInfo contains the configuration of a shovel
type ShovelInfo struct {
	// Shovel name
	Name string `json:"name"`
	// Virtual host this shovel belongs to
	Vhost string `json:"vhost"`
	// Component shovels belong to
	Component string `json:"component"`
	// Details the configuration values of the shovel
	Definition ShovelDefinition `json:"value"`
}

// ShovelDefinition contains the details of the shovel configuration
type ShovelDefinition struct {
	SourceURI              string `json:"src-uri"`
	SourceExchange         string `json:"src-exchange,omitempty"`
	SourceExchangeKey      string `json:"src-exchange-key,omitempty"`
	SourceQueue            string `json:"src-queue,omitempty"`
	DestinationURI         string `json:"dest-uri"`
	DestinationExchange    string `json:"dest-exchange,omitempty"`
	DestinationExchangeKey string `json:"dest-exchange-key,omitempty"`
	DestinationQueue       string `json:"dest-queue,omitempty"`
	PrefetchCount          int    `json:"prefetch-count,omitempty"`
	ReconnectDelay         int    `json:"reconnect-delay,omitempty"`
	AddForwardHeaders      bool   `json:"add-forward-headers"`
	AckMode                string `json:"ack-mode"`
	DeleteAfter            string `json:"delete-after"`
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
	req, err := newGETRequest(c, "parameters/shovel/"+PathEscape(vhost))
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
	req, err := newGETRequest(c, "parameters/shovel/"+PathEscape(vhost)+"/"+PathEscape(shovel))

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

	req, err := newRequestWithBody(c, "PUT", "parameters/shovel/"+PathEscape(vhost)+"/"+PathEscape(shovel), body)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/parameters/shovel/{vhost}/{name}
//

// DeleteShovel a shovel
func (c *Client) DeleteShovel(vhost, shovel string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "parameters/shovel/"+PathEscape(vhost)+"/"+PathEscape(shovel), nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
