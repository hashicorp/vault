package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type BucketPutOriginOptions struct {
	XMLName xml.Name           `xml:"OriginConfiguration"`
	Rule    []BucketOriginRule `xml:"OriginRule"`
}

type BucketOriginRule struct {
	OriginType      string                 `xml:"OriginType"`
	OriginCondition *BucketOriginCondition `xml:"OriginCondition"`
	OriginParameter *BucketOriginParameter `xml:"OriginParameter"`
	OriginInfo      *BucketOriginInfo      `xml:"OriginInfo"`
}

type BucketOriginCondition struct {
	HTTPStatusCode string `xml:"HTTPStatusCode,omitempty"`
	Prefix         string `xml:"Prefix,omitempty"`
}

type BucketOriginParameter struct {
	Protocol          string                  `xml:"Protocol,omitempty"`
	FollowQueryString bool                    `xml:"FollowQueryString,omitempty"`
	HttpHeader        *BucketOriginHttpHeader `xml:"HttpHeader,omitempty"`
	FollowRedirection bool                    `xml:"FollowRedirection,omitempty"`
	HttpRedirectCode  string                  `xml:"HttpRedirectCode,omitempty"`
	CopyOriginData    bool                    `xml:"CopyOriginData,omitempty"`
}

type BucketOriginHttpHeader struct {
	// 目前还不支持 FollowAllHeaders
	// FollowAllHeaders  bool              `xml:"FollowAllHeaders,omitempty"`
	NewHttpHeaders    []OriginHttpHeader `xml:"NewHttpHeaders>Header,omitempty"`
	FollowHttpHeaders []OriginHttpHeader `xml:"FollowHttpHeaders>Header,omitempty"`
}

type OriginHttpHeader struct {
	Key   string `xml:"Key,omitempty"`
	Value string `xml:"Value,omitempty"`
}

type BucketOriginInfo struct {
	HostInfo string                `xml:"HostInfo>HostName,omitempty"`
	FileInfo *BucketOriginFileInfo `xml:"FileInfo,omitempty"`
}
type BucketOriginFileInfo struct {
	PrefixDirective bool   `xml:"PrefixDirective,omitempty"`
	Prefix          string `xml:"Prefix,omitempty"`
	Suffix          string `xml:"Suffix,omitempty"`
}

type BucketGetOriginResult BucketPutOriginOptions

func (s *BucketService) PutOrigin(ctx context.Context, opt *BucketPutOriginOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetOrigin(ctx context.Context) (*BucketGetOriginResult, *Response, error) {
	var res BucketGetOriginResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &res, resp, err
}

func (s *BucketService) DeleteOrigin(ctx context.Context) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?origin",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}
