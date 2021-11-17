package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Bootscript represents a  Bootscript
type Bootscript struct {
	Bootcmdargs string `json:"bootcmdargs,omitempty"`
	Dtb         string `json:"dtb,omitempty"`
	Initrd      string `json:"initrd,omitempty"`
	Kernel      string `json:"kernel,omitempty"`

	// Arch is the architecture target of the bootscript
	Arch string `json:"architecture,omitempty"`

	// Identifier is a unique identifier for the bootscript
	Identifier string `json:"id,omitempty"`

	// Organization is the owner of the bootscript
	Organization string `json:"organization,omitempty"`

	// Name is a user-defined name for the bootscript
	Title string `json:"title,omitempty"`

	// Public is true for public bootscripts and false for user bootscripts
	Public bool `json:"public,omitempty"`

	Default bool `json:"default,omitempty"`
}

type getBootscriptResponse struct {
	Bootscript Bootscript `json:"bootscript,omitempty"`
}

type getBootscriptsResponse struct {
	Bootscripts []Bootscript `json:"bootscripts,omitempty"`
}

// GetBootscripts gets the list of bootscripts from the API
func (s *API) GetBootscripts() ([]Bootscript, error) {
	query := url.Values{}

	resp, err := s.GetResponsePaginate(s.computeAPI, "bootscripts", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var bootscripts getBootscriptsResponse

	if err = json.Unmarshal(body, &bootscripts); err != nil {
		return nil, err
	}
	return bootscripts.Bootscripts, nil
}

// GetBootscript gets a bootscript from the API
func (s *API) GetBootscript(bootscriptID string) (*Bootscript, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "bootscripts/"+bootscriptID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneBootscript getBootscriptResponse

	if err = json.Unmarshal(body, &oneBootscript); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	return &oneBootscript.Bootscript, nil
}
