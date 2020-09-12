package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type BucketStatement struct {
	Principal map[string][]string               `json:"principal,omitempty"`
	Action    []string                          `json:"action,omitempty"`
	Effect    string                            `json:"effect,omitempty"`
	Resource  []string                          `json:"resource,omitempty"`
	Condition map[string]map[string]interface{} `json:"condition,omitempty"`
}

type BucketPutPolicyOptions struct {
	Statement []BucketStatement   `json:"statement,omitempty"`
	Version   string              `json:"version,omitempty"`
	Principal map[string][]string `json:"principal,omitempty"`
}

type BucketGetPolicyResult BucketPutPolicyOptions

func (s *BucketService) PutPolicy(ctx context.Context, opt *BucketPutPolicyOptions) (*Response, error) {
	var f *strings.Reader
	if opt != nil {
		bs, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		body := string(bs)
		f = strings.NewReader(body)
	}
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?policy",
		method:  http.MethodPut,
		body:    f,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}

func (s *BucketService) GetPolicy(ctx context.Context) (*BucketGetPolicyResult, *Response, error) {
	var bs bytes.Buffer
	var res BucketGetPolicyResult
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?policy",
		method:  http.MethodGet,
		result:  &bs,
	}
	resp, err := s.client.send(ctx, sendOpt)
	if err == nil {
		err = json.Unmarshal(bs.Bytes(), &res)
	}
	return &res, resp, err
}

func (s *BucketService) DeletePolicy(ctx context.Context) (*Response, error) {
	sendOpt := &sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?policy",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, sendOpt)
	return resp, err
}
