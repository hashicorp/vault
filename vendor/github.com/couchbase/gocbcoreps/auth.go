package gocbcoreps

import (
	"context"
	"encoding/base64"

	"google.golang.org/grpc/credentials"
)

type GrpcBasicAuth struct {
	EncodedData string
}

// NewJWTAccessFromKey creates PerRPCCredentials from the given jsonKey.
func NewGrpcBasicAuth(username, password string) (credentials.PerRPCCredentials, error) {
	basicAuth := username + ":" + password
	authValue := base64.StdEncoding.EncodeToString([]byte(basicAuth))
	return GrpcBasicAuth{authValue}, nil
}

func (j GrpcBasicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + j.EncodedData,
	}, nil
}

func (j GrpcBasicAuth) RequireTransportSecurity() bool {
	return false
}
