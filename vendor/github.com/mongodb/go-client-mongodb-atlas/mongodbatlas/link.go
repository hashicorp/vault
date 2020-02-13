package mongodbatlas

import "net/url"

//Link is the link to sub-resources and/or related resources.
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

func (l *Link) getHrefURL() (*url.URL, error) {
	return url.Parse(l.Href)
}

func (l *Link) getHrefQueryParam(param string) (string, error) {
	hrefURL, err := l.getHrefURL()
	if err != nil {
		return "", err
	}
	return hrefURL.Query().Get(param), nil
}
