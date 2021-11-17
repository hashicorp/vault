package rabbithole

import "net/url"

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
