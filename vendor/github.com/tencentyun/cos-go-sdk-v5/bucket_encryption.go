package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type BucketEncryptionConfiguration struct {
	SSEAlgorithm string `xml:"SSEAlgorithm"`
}

type BucketPutEncryptionOptions struct {
	XMLName xml.Name                       `xml:"ServerSideEncryptionConfiguration"`
	Rule    *BucketEncryptionConfiguration `xml:"Rule>ApplySideEncryptionConfiguration"`
}

type BucketGetEncryptionResult BucketPutEncryptionOptions

func (s *BucketService) PutEncryption(ctx context.Context, opt *BucketPutEncryptionOptions) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?encryption",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetEncryption(ctx context.Context) (*BucketGetEncryptionResult, *Response, error) {
	var res BucketGetEncryptionResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?encryption",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return &res, resp, err
}

func (s *BucketService) DeleteEncryption(ctx context.Context) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?encryption",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}
