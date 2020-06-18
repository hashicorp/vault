package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type WebsiteRoutingRule struct {
	ConditionErrorCode string `xml:"Condition>HttpErrorCodeReturnedEquals,omitempty"`
	ConditionPrefix    string `xml:"Condition>KeyPrefixEquals,omitempty"`

	RedirectProtocol         string `xml:"Redirect>Protocol,omitempty"`
	RedirectReplaceKey       string `xml:"Redirect>ReplaceKeyWith,omitempty"`
	RedirectReplaceKeyPrefix string `xml:"Redirect>ReplaceKeyPrefixWith,omitempty"`
}

type WebsiteRoutingRules struct {
	Rules []WebsiteRoutingRule `xml:"RoutingRule,omitempty"`
}

type ErrorDocument struct {
	Key string `xml:"Key,omitempty"`
}

type RedirectRequestsProtocol struct {
	Protocol string `xml:"Protocol,omitempty"`
}

type BucketPutWebsiteOptions struct {
	XMLName          xml.Name                  `xml:"WebsiteConfiguration"`
	Index            string                    `xml:"IndexDocument>Suffix"`
	RedirectProtocol *RedirectRequestsProtocol `xml:"RedirectAllRequestsTo,omitempty"`
	Error            *ErrorDocument            `xml:"ErrorDocument,omitempty"`
	RoutingRules     *WebsiteRoutingRules      `xml:"RoutingRules,omitempty"`
}

type BucketGetWebsiteResult BucketPutWebsiteOptions

func (s *BucketService) PutWebsite(ctx context.Context, opt *BucketPutWebsiteOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetWebsite(ctx context.Context) (*BucketGetWebsiteResult, *Response, error) {
	var res BucketGetWebsiteResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &res, resp, err
}

func (s *BucketService) DeleteWebsite(ctx context.Context) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?website",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}
