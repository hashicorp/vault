package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// BindingVertex represents one end (vertex) of a binding,
// a source or destination. This is primarily relevant for
// exchange-to-exchange bindings (E2Es).
type BindingVertex string

const (
	// BindingSource indicates the source vertex of a binding
	BindingSource BindingVertex = "source"
	// BindingDestination indicates the source vertex of a binding
	BindingDestination BindingVertex = "destination"
)

func (v BindingVertex) String() string {
	return string(v)
}

//
// GET /api/bindings
//

// Example response:
//
// [
//   {
//     "source": "",
//     "vhost": "\/",
//     "destination": "amq.gen-Dzw36tPTm_VsmILY9oTG9w",
//     "destination_type": "queue",
//     "routing_key": "amq.gen-Dzw36tPTm_VsmILY9oTG9w",
//     "arguments": {
//
//     },
//     "properties_key": "amq.gen-Dzw36tPTm_VsmILY9oTG9w"
//   }
// ]

// BindingInfo represents details of a binding.
type BindingInfo struct {
	// Binding source (exchange name)
	Source string `json:"source"`
	Vhost  string `json:"vhost,omitempty"`
	// Binding destination (queue or exchange name)
	Destination string `json:"destination"`
	// Destination type, either "queue" or "exchange"
	DestinationType string                 `json:"destination_type"`
	RoutingKey      string                 `json:"routing_key"`
	Arguments       map[string]interface{} `json:"arguments"`
	PropertiesKey   string                 `json:"properties_key,omitempty"`
}

// ListBindings returns all bindings
func (c *Client) ListBindings() (rec []BindingInfo, err error) {
	req, err := newGETRequest(c, "bindings/")
	if err != nil {
		return []BindingInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []BindingInfo{}, err
	}

	return rec, nil
}

func (c *Client) listBindingsVia(path string) (rec []BindingInfo, err error) {
	req, err := newGETRequest(c, path)
	if err != nil {
		return []BindingInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []BindingInfo{}, err
	}

	return rec, nil
}

//
// GET /api/bindings/{vhost}
//

// ListBindingsIn returns all bindings in a virtual host.
func (c *Client) ListBindingsIn(vhost string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + url.PathEscape(vhost))
}

//
// GET /api/queues/{vhost}/{queue}/bindings
//

// Example response:
// [
//   {"source":"",
//    "vhost":"/",
//    "destination":"amq.gen-H0tnavWatL7g7uU2q5cAPA",
//    "destination_type":"queue",
//    "routing_key":"amq.gen-H0tnavWatL7g7uU2q5cAPA",
//    "arguments":{},
//    "properties_key":"amq.gen-H0tnavWatL7g7uU2q5cAPA"},
//   {"source":"temp",
//    "vhost":"/",
//    "destination":"amq.gen-H0tnavWatL7g7uU2q5cAPA",
//    "destination_type":"queue",
//    "routing_key":"",
//    "arguments":{},
//    "properties_key":"~"}
// ]

// ListQueueBindings returns all bindings of individual queue.
func (c *Client) ListQueueBindings(vhost, queue string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("queues/" + url.PathEscape(vhost) + "/" + url.PathEscape(queue) + "/bindings")
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/source
//

// ListExchangeBindingsWithSource returns exchange-to-exchange (E2E) bindings where
// the given exchange is the source.
func (c *Client) ListExchangeBindingsWithSource(vhost, exchange string) (rec []BindingInfo, err error) {
	return c.ListExchangeBindings(vhost, exchange, BindingSource)
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/destination
//

// ListExchangeBindingsWithDestination returns exchange-to-exchange (E2E) bindings where
// the given exchange is the destination.
func (c *Client) ListExchangeBindingsWithDestination(vhost, exchange string) (rec []BindingInfo, err error) {
	return c.ListExchangeBindings(vhost, exchange, BindingDestination)
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/{source-or-destination}
//

// ListExchangeBindings returns all bindings having the exchange as source or destination as defined by the Target
func (c *Client) ListExchangeBindings(vhost, exchange string, sourceOrDestination BindingVertex) (rec []BindingInfo, err error) {
	return c.listBindingsVia("exchanges/" + url.PathEscape(vhost) + "/" + url.PathEscape(exchange) + "/bindings/" + sourceOrDestination.String())
}

//
// GET /api/bindings/{vhost}/e/{source}/e/{destination}
//

// ListExchangeBindingsBetween returns a set of bindings between two exchanges.
func (c *Client) ListExchangeBindingsBetween(vhost, source string, destination string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + url.PathEscape(vhost) + "/e/" + url.PathEscape(source) + "/e/" + url.PathEscape(destination))
}

//
// GET /api/bindings/{vhost}/e/{exchange}/q/{queue}
//

// ListQueueBindingsBetween returns a set of bindings between an exchange and a queue.
func (c *Client) ListQueueBindingsBetween(vhost, exchange string, queue string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + url.PathEscape(vhost) + "/e/" + url.PathEscape(exchange) + "/q/" + url.PathEscape(queue))
}

//
// POST /api/bindings/{vhost}/e/{source}/{destination_type}/{destination}
//

// DeclareBinding adds a new binding
func (c *Client) DeclareBinding(vhost string, info BindingInfo) (res *http.Response, err error) {
	info.Vhost = vhost

	if info.Arguments == nil {
		info.Arguments = make(map[string]interface{})
	}
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "POST", c.newBindingPath(vhost, info), body)

	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) newBindingPath(vhost string, info BindingInfo) string {
	if info.DestinationType == "queue" {
		// /api/bindings/{vhost}/e/{exchange}/q/{queue}
		return "bindings/" + url.PathEscape(vhost) +
			"/e/" + url.PathEscape(info.Source) +
			"/q/" + url.PathEscape(info.Destination)
	}
	// /api/bindings/{vhost}/e/{source}/e/{destination}
	return "bindings/" + url.PathEscape(vhost) +
		"/e/" + url.PathEscape(info.Source) +
		"/e/" + url.PathEscape(info.Destination)
}

//
// DELETE /api/bindings/{vhost}/e/{source}/{destination_type}/{destination}/{props}
//

// DeleteBinding deletes an individual binding
func (c *Client) DeleteBinding(vhost string, info BindingInfo) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "bindings/"+url.PathEscape(vhost)+
		"/e/"+url.PathEscape(info.Source)+"/"+url.PathEscape(string(info.DestinationType[0]))+
		"/"+url.PathEscape(info.Destination)+"/"+url.PathEscape(info.PropertiesKey), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
