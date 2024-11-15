package godo

import (
	"context"
	"net/url"
	"strconv"
)

// Links manages links that are returned along with a List
type Links struct {
	Pages   *Pages       `json:"pages,omitempty"`
	Actions []LinkAction `json:"actions,omitempty"`
}

// Pages are pages specified in Links
type Pages struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
}

// LinkAction is a pointer to an action
type LinkAction struct {
	ID   int    `json:"id,omitempty"`
	Rel  string `json:"rel,omitempty"`
	HREF string `json:"href,omitempty"`
}

// CurrentPage is current page of the list
func (l *Links) CurrentPage() (int, error) {
	return l.Pages.current()
}

// NextPageToken is the page token to request the next page of the list
func (l *Links) NextPageToken() (string, error) {
	return l.Pages.nextPageToken()
}

// PrevPageToken is the page token to request the previous page of the list
func (l *Links) PrevPageToken() (string, error) {
	return l.Pages.prevPageToken()
}

func (p *Pages) current() (int, error) {
	switch {
	case p == nil:
		return 1, nil
	case p.Prev == "" && p.Next != "":
		return 1, nil
	case p.Prev != "":
		prevPage, err := pageForURL(p.Prev)
		if err != nil {
			return 0, err
		}

		return prevPage + 1, nil
	}

	return 0, nil
}

func (p *Pages) nextPageToken() (string, error) {
	if p == nil || p.Next == "" {
		return "", nil
	}
	token, err := pageTokenFromURL(p.Next)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *Pages) prevPageToken() (string, error) {
	if p == nil || p.Prev == "" {
		return "", nil
	}
	token, err := pageTokenFromURL(p.Prev)
	if err != nil {
		return "", err
	}
	return token, nil
}

// IsLastPage returns true if the current page is the last
func (l *Links) IsLastPage() bool {
	if l.Pages == nil {
		return true
	}
	return l.Pages.isLast()
}

func (p *Pages) isLast() bool {
	return p.Next == ""
}

func pageForURL(urlText string) (int, error) {
	u, err := url.ParseRequestURI(urlText)
	if err != nil {
		return 0, err
	}

	pageStr := u.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}

	return page, nil
}

func pageTokenFromURL(urlText string) (string, error) {
	u, err := url.ParseRequestURI(urlText)
	if err != nil {
		return "", err
	}
	return u.Query().Get("page_token"), nil
}

// Get a link action by id.
func (la *LinkAction) Get(ctx context.Context, client *Client) (*Action, *Response, error) {
	return client.Actions.Get(ctx, la.ID)
}
