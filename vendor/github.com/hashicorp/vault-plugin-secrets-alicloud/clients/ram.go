// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
)

func NewRAMClient(sdkConfig *sdk.Config, key, secret string) (*RAMClient, error) {
	creds, err := chainedCreds(key, secret)
	if err != nil {
		return nil, err
	}
	// We hard-code a region here because there's only one RAM endpoint regardless of the
	// region you're in.
	client, err := ram.NewClientWithOptions("us-east-1", sdkConfig, creds)
	if err != nil {
		return nil, err
	}
	return &RAMClient{client: client}, nil
}

type RAMClient struct {
	client *ram.Client
}

func (c *RAMClient) CreateAccessKey(userName string) (*ram.CreateAccessKeyResponse, error) {
	accessKeyReq := ram.CreateCreateAccessKeyRequest()
	accessKeyReq.UserName = userName
	return c.client.CreateAccessKey(accessKeyReq)
}

func (c *RAMClient) DeleteAccessKey(userName, accessKeyID string) error {
	req := ram.CreateDeleteAccessKeyRequest()
	req.UserAccessKeyId = accessKeyID
	req.UserName = userName
	_, err := c.client.DeleteAccessKey(req)
	return err
}

func (c *RAMClient) CreatePolicy(policyName, policyDoc string) (*ram.CreatePolicyResponse, error) {
	createPolicyReq := ram.CreateCreatePolicyRequest()
	createPolicyReq.PolicyName = policyName
	createPolicyReq.Description = "Created by Vault."
	createPolicyReq.PolicyDocument = policyDoc
	return c.client.CreatePolicy(createPolicyReq)
}

func (c *RAMClient) DeletePolicy(policyName string) error {
	req := ram.CreateDeletePolicyRequest()
	req.PolicyName = policyName
	_, err := c.client.DeletePolicy(req)
	return err
}

func (c *RAMClient) AttachPolicy(userName, policyName, policyType string) error {
	attachPolReq := ram.CreateAttachPolicyToUserRequest()
	attachPolReq.UserName = userName
	attachPolReq.PolicyName = policyName
	attachPolReq.PolicyType = policyType
	_, err := c.client.AttachPolicyToUser(attachPolReq)
	return err
}

func (c *RAMClient) DetachPolicy(userName, policyName, policyType string) error {
	req := ram.CreateDetachPolicyFromUserRequest()
	req.UserName = userName
	req.PolicyName = policyName
	req.PolicyType = policyType
	_, err := c.client.DetachPolicyFromUser(req)
	return err
}

func (c *RAMClient) CreateUser(userName string) (*ram.CreateUserResponse, error) {
	createUserReq := ram.CreateCreateUserRequest()
	createUserReq.UserName = userName
	createUserReq.DisplayName = userName
	return c.client.CreateUser(createUserReq)
}

// Note: deleteUser will fail if the user is presently associated with anything
// in Alibaba.
func (c *RAMClient) DeleteUser(userName string) error {
	req := ram.CreateDeleteUserRequest()
	req.UserName = userName
	_, err := c.client.DeleteUser(req)
	return err
}
