// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

func NewSTSClient(sdkConfig *sdk.Config, key, secret string) (*STSClient, error) {
	creds, err := chainedCreds(key, secret)
	if err != nil {
		return nil, err
	}
	// We hard-code a region here because there's only one RAM endpoint regardless of the
	// region you're in.
	client, err := sts.NewClientWithOptions("us-east-1", sdkConfig, creds)
	if err != nil {
		return nil, err
	}
	return &STSClient{client: client}, nil
}

type STSClient struct {
	client *sts.Client
}

func (c *STSClient) AssumeRole(roleSessionName, roleARN string) (*sts.AssumeRoleResponse, error) {
	assumeRoleReq := sts.CreateAssumeRoleRequest()
	assumeRoleReq.RoleArn = roleARN
	assumeRoleReq.RoleSessionName = roleSessionName
	return c.client.AssumeRole(assumeRoleReq)
}
