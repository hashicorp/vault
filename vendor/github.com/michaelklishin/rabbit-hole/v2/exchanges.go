package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

//
// GET /api/exchanges
//

// IngressEgressStats represents common message flow metrics.
type IngressEgressStats struct {
	PublishIn        int          `json:"publish_in,omitempty"`
	PublishInDetails *RateDetails `json:"publish_in_details,omitempty"`

	PublishOut        int          `json:"publish_out,omitempty"`
	PublishOutDetails *RateDetails `json:"publish_out_details,omitempty"`
}

// ExchangeInfo represents and exchange and its properties.
type ExchangeInfo struct {
	Name       string                 `json:"name"`
	Vhost      string                 `json:"vhost,omitempty"`
	Type       string                 `json:"type"`
	Durable    bool                   `json:"durable"`
	AutoDelete AutoDelete             `json:"auto_delete"`
	Internal   bool                   `json:"internal"`
	Arguments  map[string]interface{} `json:"arguments"`

	MessageStats *IngressEgressStats `json:"message_stats,omitempty"`
}

// ExchangeSettings is a set of exchange properties. Use this type when declaring
// an exchange.
type ExchangeSettings struct {
	Type       string                 `json:"type"`
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete,omitempty"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

// ListExchanges lists all exchanges in a cluster. This only includes exchanges in the
// virtual hosts accessible to the user.
func (c *Client) ListExchanges() (rec []ExchangeInfo, err error) {
	req, err := newGETRequest(c, "exchanges")
	if err != nil {
		return []ExchangeInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ExchangeInfo{}, err
	}

	return rec, nil
}

//
// GET /api/exchanges/{vhost}
//

// ListExchangesIn lists all exchanges in a virtual host.
func (c *Client) ListExchangesIn(vhost string) (rec []ExchangeInfo, err error) {
	req, err := newGETRequest(c, "exchanges/"+url.PathEscape(vhost))
	if err != nil {
		return []ExchangeInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ExchangeInfo{}, err
	}

	return rec, nil
}

//
// GET /api/exchanges/{vhost}/{name}
//

// Example response:
//
// {
//   "incoming": [
//     {
//       "stats": {
//         "publish": 2760,
//         "publish_details": {
//           "rate": 20
//         }
//       },
//       "channel_details": {
//         "name": "127.0.0.1:46928 -> 127.0.0.1:5672 (2)",
//         "number": 2,
//         "connection_name": "127.0.0.1:46928 -> 127.0.0.1:5672",
//         "peer_port": 46928,
//         "peer_host": "127.0.0.1"
//       }
//     }
//   ],
//   "outgoing": [
//     {
//       "stats": {
//         "publish": 1280,
//         "publish_details": {
//           "rate": 20
//         }
//       },
//       "queue": {
//         "name": "amq.gen-7NhO_yRr4lDdp-8hdnvfuw",
//         "vhost": "rabbit\/hole"
//       }
//     }
//   ],
//   "message_stats": {
//     "publish_in": 2760,
//     "publish_in_details": {
//       "rate": 20
//     },
//     "publish_out": 1280,
//     "publish_out_details": {
//       "rate": 20
//     }
//   },
//   "name": "amq.fanout",
//   "vhost": "rabbit\/hole",
//   "type": "fanout",
//   "durable": true,
//   "auto_delete": false,
//   "internal": false,
//   "arguments": {
//   }
// }

// ExchangeIngressDetails represents ingress (inbound) message flow metrics of an exchange.
type ExchangeIngressDetails struct {
	Stats          MessageStats      `json:"stats"`
	ChannelDetails PublishingChannel `json:"channel_details"`
}

// PublishingChannel represents a channel and its basic properties.
type PublishingChannel struct {
	Number         int    `json:"number"`
	Name           string `json:"name"`
	ConnectionName string `json:"connection_name"`
	PeerPort       Port   `json:"peer_port"`
	PeerHost       string `json:"peer_host"`
}

// NameAndVhost repesents a named entity in a virtual host.
type NameAndVhost struct {
	Name  string `json:"name"`
	Vhost string `json:"vhost"`
}

// ExchangeEgressDetails represents egress (outbound) message flow metrics of an exchange.
type ExchangeEgressDetails struct {
	Stats MessageStats `json:"stats"`
	Queue NameAndVhost `json:"queue"`
}

// DetailedExchangeInfo represents an exchange with all of its properties and metrics.
type DetailedExchangeInfo struct {
	Name       string                 `json:"name"`
	Vhost      string                 `json:"vhost"`
	Type       string                 `json:"type"`
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete"`
	Internal   bool                   `json:"internal"`
	Arguments  map[string]interface{} `json:"arguments"`

	Incoming     []ExchangeIngressDetails `json:"incoming"`
	Outgoing     []ExchangeEgressDetails  `json:"outgoing"`
	PublishStats IngressEgressStats       `json:"message_stats"`
}

// GetExchange returns information about an exchange.
func (c *Client) GetExchange(vhost, exchange string) (rec *DetailedExchangeInfo, err error) {
	req, err := newGETRequest(c, "exchanges/"+url.PathEscape(vhost)+"/"+url.PathEscape(exchange))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/exchanges/{vhost}/{exchange}
//

// DeclareExchange declares an exchange.
func (c *Client) DeclareExchange(vhost, exchange string, info ExchangeSettings) (res *http.Response, err error) {
	if info.Arguments == nil {
		info.Arguments = make(map[string]interface{})
	}
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "exchanges/"+url.PathEscape(vhost)+"/"+url.PathEscape(exchange), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/exchanges/{vhost}/{name}
//

// DeleteExchange deletes an exchange.
func (c *Client) DeleteExchange(vhost, exchange string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "exchanges/"+url.PathEscape(vhost)+"/"+url.PathEscape(exchange), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
