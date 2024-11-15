package rabbithole

import (
	"net/url"
)

// ListVhostConnections returns the current connections to a specified vhost
func (c *Client) ListVhostConnections(vhostname string) (rec []ConnectionInfo, err error) {
	req, err := newGETRequest(c, "vhosts/"+url.PathEscape(vhostname)+"/connections")
	if err != nil {
		return []ConnectionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []ConnectionInfo{}, err
	}

	return rec, nil
}
