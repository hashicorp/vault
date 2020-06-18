package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type BucketPutRefererOptions struct {
	XMLName                 xml.Name `xml:"RefererConfiguration"`
	Status                  string   `xml:"Status"`
	RefererType             string   `xml:"RefererType"`
	DomainList              []string `xml:"DomainList>Domain"`
	EmptyReferConfiguration string   `xml:"EmptyReferConfiguration,omitempty"`
}

type BucketGetRefererResult BucketPutRefererOptions

func (s *BucketService) PutReferer(ctx context.Context, opt *BucketPutRefererOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?referer",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetReferer(ctx context.Context) (*BucketGetRefererResult, *Response, error) {
	var res BucketGetRefererResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?referer",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &res, resp, err
}
