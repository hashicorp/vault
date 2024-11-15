package rabbithole

import (
	"net/http"
	"net/url"
)

// ConnectionInfo provides information about connection to a RabbitMQ node.
type ConnectionInfo struct {
	// Connection name
	Name string `json:"name"`
	// Node the client is connected to
	Node string `json:"node"`
	// Number of open channels
	Channels int `json:"channels"`
	// Connection state
	State string `json:"state"`
	// Connection type, network (via AMQP client) or direct (via direct Erlang client)
	Type string `json:"type"`

	// Server port
	Port Port `json:"port"`
	// Client port
	PeerPort Port `json:"peer_port"`

	// Server host
	Host string `json:"host"`
	// Client host
	PeerHost string `json:"peer_host"`

	// Last connection blocking reason, if any
	LastBlockedBy string `json:"last_blocked_by"`
	// When connection was last blocked
	LastBlockedAge string `json:"last_blocked_age"`

	// True if connection uses TLS/SSL
	UsesTLS bool `json:"ssl"`
	// Client certificate subject
	PeerCertSubject string `json:"peer_cert_subject"`
	// Client certificate validity
	PeerCertValidity string `json:"peer_cert_validity"`
	// Client certificate issuer
	PeerCertIssuer string `json:"peer_cert_issuer"`

	// TLS/SSL protocol and version
	SSLProtocol string `json:"ssl_protocol"`
	// Key exchange mechanism
	SSLKeyExchange string `json:"ssl_key_exchange"`
	// SSL cipher suite used
	SSLCipher string `json:"ssl_cipher"`
	// SSL hash
	SSLHash string `json:"ssl_hash"`

	// Protocol, e.g. AMQP 0-9-1 or MQTT 3-1
	Protocol string `json:"protocol"`
	User     string `json:"user"`
	// Virtual host
	Vhost string `json:"vhost"`

	// Heartbeat timeout
	Timeout int `json:"timeout"`
	// Maximum frame size (AMQP 0-9-1)
	FrameMax int `json:"frame_max"`

	// A map of client properties (name, version, capabilities, etc)
	ClientProperties Properties `json:"client_properties"`

	// Octets received
	RecvOct uint64 `json:"recv_oct"`
	// Octets sent
	SendOct     uint64 `json:"send_oct"`
	RecvCount   uint64 `json:"recv_cnt"`
	SendCount   uint64 `json:"send_cnt"`
	SendPending uint64 `json:"send_pend"`
	// Ingress data rate
	RecvOctDetails RateDetails `json:"recv_oct_details"`
	// Egress data rate
	SendOctDetails RateDetails `json:"send_oct_details"`

	// Connection timestamp
	ConnectedAt uint64 `json:"connected_at,omitempty"`
}

// Connection of a specific user. This provides just enough information
// to the monitoring tools.
type UserConnectionInfo struct {
	// Connection name
	Name string `json:"name"`
	// Node the client is connected to
	Node string `json:"node"`
	// Username
	User string `json:"user"`
	// Virtual host
	Vhost string `json:"vhost"`
}

//
// GET /api/connections
//

// ListConnections returns a list of client connections to target node.
func (c *Client) ListConnections() (rec []ConnectionInfo, err error) {
	req, err := newGETRequest(c, "connections")
	if err != nil {
		return []ConnectionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ConnectionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/connections/username/{username}
//

// ListConnections returns a list of client connections to target node.
func (c *Client) ListConnectionsOfUser(username string) (rec []UserConnectionInfo, err error) {
	req, err := newGETRequest(c, "connections/username/"+url.PathEscape(username))
	if err != nil {
		return []UserConnectionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []UserConnectionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/connections/{name}
//

// GetConnection retrieves information about a connection.
func (c *Client) GetConnection(name string) (rec *ConnectionInfo, err error) {
	req, err := newGETRequest(c, "connections/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// DELETE /api/connections/{name}
//

// CloseConnection closes a connection.
func (c *Client) CloseConnection(name string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "connections/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/connections/username/{username}
//

// CloseConnection closes a connection.
func (c *Client) CloseAllConnectionsOfUser(username string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "connections/username/"+url.PathEscape(username), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
