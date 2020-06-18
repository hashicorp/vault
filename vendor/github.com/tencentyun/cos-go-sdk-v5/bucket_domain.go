package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type BucketPutDomainOptions struct {
	XMLName           xml.Name `xml:"DomainConfiguration"`
	Status            string   `xml:"DomainRule>Status"`
	Name              string   `xml:"DomainRule>Name"`
	Type              string   `xml:"DomainRule>Type"`
	ForcedReplacement string   `xml:"DomainRule>ForcedReplacement,omitempty"`
}
type BucketGetDomainResult BucketPutDomainOptions

func (s *BucketService) PutDomain(ctx context.Context, opt *BucketPutDomainOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?domain",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetDomain(ctx context.Context) (*BucketGetDomainResult, *Response, error) {
	var res BucketGetDomainResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?domain",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &res, resp, err
}
