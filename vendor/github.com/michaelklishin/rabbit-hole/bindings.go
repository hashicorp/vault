package rabbithole

import (
	"encoding/json"
	"net/http"
)

type BindingVertex string

const (
	BindingSource      BindingVertex = "source"
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

type BindingInfo struct {
	// Binding source (exchange name)
	Source string `json:"source"`
	Vhost  string `json:"vhost"`
	// Binding destination (queue or exchange name)
	Destination string `json:"destination"`
	// Destination type, either "queue" or "exchange"
	DestinationType string                 `json:"destination_type"`
	RoutingKey      string                 `json:"routing_key"`
	Arguments       map[string]interface{} `json:"arguments"`
	PropertiesKey   string                 `json:"properties_key"`
}

// Returns all bindings
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

// Returns all bindings in a virtual host.
func (c *Client) ListBindingsIn(vhost string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + PathEscape(vhost))
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

// Returns all bindings of individual queue.
func (c *Client) ListQueueBindings(vhost, queue string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("queues/" + PathEscape(vhost) + "/" + PathEscape(queue) + "/bindings")
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/source
//

func (c *Client) ListExchangeBindingsWithSource(vhost, exchange string) (rec []BindingInfo, err error) {
	return c.ListExchangeBindings(vhost, exchange, BindingSource)
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/destination
//

func (c *Client) ListExchangeBindingsWithDestination(vhost, exchange string) (rec []BindingInfo, err error) {
	return c.ListExchangeBindings(vhost, exchange, BindingDestination)
}

//
// GET /api/exchanges/{vhost}/{exchange}/bindings/{source-or-destination}
//

// Returns all bindings having the exchange as source or destination as defined by the Target
func (c *Client) ListExchangeBindings(vhost, exchange string, sourceOrDestination BindingVertex) (rec []BindingInfo, err error) {
	return c.listBindingsVia("exchanges/" + PathEscape(vhost) + "/" + PathEscape(exchange) + "/bindings/" + sourceOrDestination.String())
}

//
// GET /api/bindings/{vhost}/e/{source}/e/{destination}
//

func (c *Client) ListExchangeBindingsBetween(vhost, source string, destination string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + PathEscape(vhost) + "/e/" + PathEscape(source) + "/e/" + destination)
}

//
// GET /api/bindings/{vhost}/e/{exchange}/q/{queue}
//

func (c *Client) ListQueueBindingsBetween(vhost, exchange string, queue string) (rec []BindingInfo, err error) {
	return c.listBindingsVia("bindings/" + PathEscape(vhost) + "/e/" + PathEscape(exchange) + "/q/" + queue)
}

//
// POST /api/bindings/{vhost}/e/{source}/{destination_type}/{destination}
//

// DeclareBinding updates information about a binding between a source and a target
func (c *Client) DeclareBinding(vhost string, info BindingInfo) (res *http.Response, err error) {
	info.Vhost = vhost

	if info.Arguments == nil {
		info.Arguments = make(map[string]interface{})
	}
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "POST", "bindings/"+PathEscape(vhost)+
		"/e/"+PathEscape(info.Source)+"/"+PathEscape(string(info.DestinationType[0]))+
		"/"+PathEscape(info.Destination), body)

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
// DELETE /api/bindings/{vhost}/e/{source}/{destination_type}/{destination}/{props}
//

// DeleteBinding delets an individual binding
func (c *Client) DeleteBinding(vhost string, info BindingInfo) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "bindings/"+PathEscape(vhost)+
		"/e/"+PathEscape(info.Source)+"/"+PathEscape(string(info.DestinationType[0]))+
		"/"+PathEscape(info.Destination)+"/"+PathEscape(info.PropertiesKey), nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
