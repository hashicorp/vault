package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ExportedDefinitions represents definitions exported from a RabbitMQ cluster
type ExportedDefinitions struct {
	RabbitVersion    string                    `json:"rabbit_version,omitempty"`
	RabbitMQVersion  string                    `json:"rabbitmq_version,omitempty"`
	ProductName      string                    `json:"product_name,omitempty"`
	ProductVersion   string                    `json:"product_version,omitempty"`
	Users            *[]UserInfo               `json:"users,omitempty"`
	Vhosts           *[]VhostInfo              `json:"vhosts,omitempty"`
	Permissions      *[]Permissions            `json:"permissions,omitempty"`
	TopicPermissions *[]TopicPermissionInfo    `json:"topic_permissions,omitempty"`
	Parameters       *[]RuntimeParameter       `json:"paramaters,omitempty"`
	GlobalParameters *[]GlobalRuntimeParameter `json:"global_parameters,omitempty"`
	Policies         *[]PolicyDefinition       `json:"policies"`
	Queues           *[]QueueInfo              `json:"queues"`
	Exchanges        *[]ExchangeInfo           `json:"exchanges"`
	Bindings         *[]BindingInfo            `json:"bindings"`
}

//
// GET /api/definitions
//

// ListDefinitions returns a set of definitions exported from a RabbitMQ cluster.
func (c *Client) ListDefinitions() (p *ExportedDefinitions, err error) {
	req, err := newGETRequest(c, "definitions")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &p); err != nil {
		return nil, err
	}

	return p, nil
}

//
// GET /api/definitions/vhost
//

// ListVhostDefinitions returns a set of definitions for a specific vhost.
func (c *Client) ListVhostDefinitions(vhost string) (p *ExportedDefinitions, err error) {
	req, err := newGETRequest(c, "definitions/"+url.QueryEscape(vhost))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &p); err != nil {
		return nil, err
	}

	return p, nil
}

//
// POST /api/definitions
//

// UploadDefinitions uploads a set of definitions and returns an error indicating if the operation was a failure
func (c *Client) UploadDefinitions(p *ExportedDefinitions) (res *http.Response, err error) {
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	req, err := newRequestWithBody(c, http.MethodPost, "definitions", body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}
	return res, nil
}

//
// POST /api/definitions/vhost
//

// UploadDefinitions uploads a set of definitions and returns an error indicating if the operation was a failure
func (c *Client) UploadVhostDefinitions(p *ExportedDefinitions, vhost string) (res *http.Response, err error) {
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	req, err := newRequestWithBody(c, http.MethodPost, "definitions/"+url.QueryEscape(vhost), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}
	return res, nil
}
