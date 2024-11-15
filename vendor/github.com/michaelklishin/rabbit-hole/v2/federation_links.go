package rabbithole

import "net/url"

type FederationLinkMap = []map[string]interface{}

//
// GET /api/federation-links
//

// ListFederationLinks returns a list of all federation links.
func (c *Client) ListFederationLinks() (links FederationLinkMap, err error) {
	req, err := newGETRequest(c, "federation-links")
	if err != nil {
		return links, err
	}

	if err = executeAndParseRequest(c, req, &links); err != nil {
		return links, err
	}

	return links, nil
}

//
// GET /api/federation-links/{vhost}
//

// ListFederationLinksIn returns a list of federation links in a vhost.
func (c *Client) ListFederationLinksIn(vhost string) (links FederationLinkMap, err error) {
	req, err := newGETRequest(c, "federation-links/"+url.PathEscape(vhost))
	if err != nil {
		return links, err
	}

	if err = executeAndParseRequest(c, req, &links); err != nil {
		return links, err
	}

	return links, nil
}
