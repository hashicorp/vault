package awsutil

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-cleanhttp"
)

var InstanceMetadataService = &instanceMetadataService{
	BaseURL:                "http://169.254.169.254",
	LatestEndpoint:         "/latest",
	IdentityEndpoint:       "/latest/dynamic/instance-identity",
	LatestMetadataEndpoint: "/latest/meta-data/",
	LatestTokenEndpoint:    "/latest/api/token",
	TokenTTLHeader:         "X-aws-ec2-metadata-token-ttl-seconds",
	DefaultTokenTTLSec:     21600, // 6 hours, default in AWS documentation.
	TokenHeader:            "X-aws-ec2-metadata-token",
}

type instanceMetadataService struct {
	BaseURL string

	// Endpoints.
	LatestEndpoint         string
	IdentityEndpoint       string
	LatestMetadataEndpoint string
	LatestTokenEndpoint    string

	// Headers and related values.
	TokenTTLHeader     string
	TokenHeader        string
	DefaultTokenTTLSec int
}

/*
PrepareRequest retrieves and attaches a token for
v2 of the instance metadata service if it's available, or leaves
the request unaltered if it's not available.

References:
- https://aws.amazon.com/blogs/security/defense-in-depth-open-firewalls-reverse-proxies-ssrf-vulnerabilities-ec2-instance-metadata-service/
- https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html
*/
func (s *instanceMetadataService) PrepareRequest(req *http.Request) error {
	token, err := s.attemptV2Token()
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set(s.TokenHeader, token)
	}
	return nil
}

func (s instanceMetadataService) attemptV2Token() (string, error) {
	tokenReq, err := http.NewRequest(http.MethodPut, s.BaseURL+s.LatestTokenEndpoint, nil)
	if err != nil {
		return "", err
	}
	tokenReq.Header.Set(s.TokenTTLHeader, strconv.Itoa(s.DefaultTokenTTLSec))
	resp, err := cleanhttp.DefaultClient().Do(tokenReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 400:
		return "", errors.New("the put request is not valid")
	case 403:
		// The instance metadata service is turned off.
		return "", nil
	}
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
