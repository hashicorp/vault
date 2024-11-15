package rabbithole

import (
	"net/http"
	"net/url"
)

// FederationDefinition represents settings
// that will be used by federation links.
type FederationDefinition struct {
	Uri            URISet `json:"uri"`
	Expires        int    `json:"expires,omitempty"`
	MessageTTL     int32  `json:"message-ttl,omitempty"`
	MaxHops        int    `json:"max-hops,omitempty"`
	PrefetchCount  int    `json:"prefetch-count,omitempty"`
	ReconnectDelay int    `json:"reconnect-delay"`
	AckMode        string `json:"ack-mode,omitempty"`
	TrustUserId    bool   `json:"trust-user-id"`
	Exchange       string `json:"exchange,omitempty"`
	Queue          string `json:"queue,omitempty"`
}

// FederationUpstream represents a configured federation upstream.
type FederationUpstream struct {
	Name       string               `json:"name"`
	Vhost      string               `json:"vhost"`
	Component  string               `json:"component"`
	Definition FederationDefinition `json:"value"`
}

// FederationUpstreamComponent is the name of the runtime parameter component
// used by federation upstreams.
const FederationUpstreamComponent string = "federation-upstream"

//
// GET /api/parameters/federation-upstream
//

// ListFederationUpstreams returns a list of all federation upstreams.
func (c *Client) ListFederationUpstreams() (ups []FederationUpstream, err error) {
	req, err := newGETRequest(c, "parameters/"+FederationUpstreamComponent)
	if err != nil {
		return []FederationUpstream{}, err
	}

	if err = executeAndParseRequest(c, req, &ups); err != nil {
		return []FederationUpstream{}, err
	}

	return ups, nil
}

//
// GET /api/parameters/federation-upstream/{vhost}
//

// ListFederationUpstreamsIn returns a list of all federation upstreams in a vhost.
func (c *Client) ListFederationUpstreamsIn(vhost string) (ups []FederationUpstream, err error) {
	req, err := newGETRequest(c, "parameters/"+FederationUpstreamComponent+"/"+url.PathEscape(vhost))
	if err != nil {
		return []FederationUpstream{}, err
	}

	if err = executeAndParseRequest(c, req, &ups); err != nil {
		return []FederationUpstream{}, err
	}

	return ups, nil
}

//
// GET /api/parameters/federation-upstream/{vhost}/{upstream}
//

// GetFederationUpstream returns information about a federation upstream.
func (c *Client) GetFederationUpstream(vhost, name string) (up *FederationUpstream, err error) {
	req, err := newGETRequest(c, "parameters/"+FederationUpstreamComponent+"/"+url.PathEscape(vhost)+"/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &up); err != nil {
		return nil, err
	}

	return up, nil
}

//
// PUT /api/parameters/federation-upstream/{vhost}/{upstream}
//

// PutFederationUpstream creates or updates a federation upstream configuration.
func (c *Client) PutFederationUpstream(vhost, name string, def FederationDefinition) (res *http.Response, err error) {
	return c.PutRuntimeParameter(FederationUpstreamComponent, vhost, name, def)
}

//
// DELETE /api/parameters/federation-upstream/{vhost}/{name}
//

// DeleteFederationUpstream removes a federation upstream.
func (c *Client) DeleteFederationUpstream(vhost, name string) (res *http.Response, err error) {
	return c.DeleteRuntimeParameter(FederationUpstreamComponent, vhost, name)
}
