// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const ec2DescribePolicy = `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": ["ec2:DescribeInstances"], "Resource": "*"}]}`

// ec2AllPolicy also uses a string instead of a list for the Action
const ec2AllPolicy = `{"Version": "2012-10-17","Statement": [{"Effect": "Allow", "Action": "ec2:*", "Resource": "*"}]}`

// ec2SingleStatement is an example of the Statement portion containing a single statement that's not a list
const ec2SingleStatement = `{"Version": "2012-10-17", "Statement": {"Effect": "Allow", "Action": ["ec2:DescribeInstances"], "Resource": "*"}}`

type mockGroupIAMClient struct {
	iamiface.IAMAPI
	ListAttachedGroupPoliciesResp iam.ListAttachedGroupPoliciesOutput
	ListGroupPoliciesResp         iam.ListGroupPoliciesOutput
	GetGroupPolicyResp            iam.GetGroupPolicyOutput
}

func (m mockGroupIAMClient) ListAttachedGroupPoliciesWithContext(_ aws.Context, in *iam.ListAttachedGroupPoliciesInput, _ ...request.Option) (*iam.ListAttachedGroupPoliciesOutput, error) {
	return &m.ListAttachedGroupPoliciesResp, nil
}

func (m mockGroupIAMClient) ListGroupPoliciesWithContext(_ aws.Context, in *iam.ListGroupPoliciesInput, _ ...request.Option) (*iam.ListGroupPoliciesOutput, error) {
	return &m.ListGroupPoliciesResp, nil
}

func (m mockGroupIAMClient) GetGroupPolicyWithContext(_ aws.Context, in *iam.GetGroupPolicyInput, _ ...request.Option) (*iam.GetGroupPolicyOutput, error) {
	return &m.GetGroupPolicyResp, nil
}

func Test_getGroupPolicies(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description         string
		listAGPResp         iam.ListAttachedGroupPoliciesOutput
		listGPResp          iam.ListGroupPoliciesOutput
		getGPResp           iam.GetGroupPolicyOutput
		iamGroupArg         []string
		wantGroupPolicies   []string
		wantGroupPolicyARNs []string
		wantErr             bool
	}{
		{
			description: "All IAM calls respond with data",
			listAGPResp: iam.ListAttachedGroupPoliciesOutput{
				AttachedPolicies: []*iam.AttachedPolicy{
					{
						PolicyArn:  aws.String("abcdefghijklmnopqrst"),
						PolicyName: aws.String("test policy"),
					},
				},
			},
			listGPResp: iam.ListGroupPoliciesOutput{
				PolicyNames: []*string{
					aws.String("inline policy"),
				},
			},
			getGPResp: iam.GetGroupPolicyOutput{
				GroupName:      aws.String("inline policy"),
				PolicyDocument: aws.String(ec2DescribePolicy),
				PolicyName:     aws.String("ec2 describe"),
			},
			iamGroupArg:         []string{"testgroup1"},
			wantGroupPolicies:   []string{ec2DescribePolicy},
			wantGroupPolicyARNs: []string{"abcdefghijklmnopqrst"},
			wantErr:             false,
		},
		{
			description: "No managed policies",
			listAGPResp: iam.ListAttachedGroupPoliciesOutput{},
			listGPResp: iam.ListGroupPoliciesOutput{
				PolicyNames: []*string{
					aws.String("inline policy"),
				},
			},
			getGPResp: iam.GetGroupPolicyOutput{
				GroupName:      aws.String("inline policy"),
				PolicyDocument: aws.String(ec2DescribePolicy),
				PolicyName:     aws.String("ec2 describe"),
			},
			iamGroupArg:         []string{"testgroup1", "testgroup2"},
			wantGroupPolicies:   []string{ec2DescribePolicy, ec2DescribePolicy},
			wantGroupPolicyARNs: []string(nil),
			wantErr:             false,
		},
		{
			description: "No inline policies",
			listAGPResp: iam.ListAttachedGroupPoliciesOutput{
				AttachedPolicies: []*iam.AttachedPolicy{
					{
						PolicyArn:  aws.String("abcdefghijklmnopqrst"),
						PolicyName: aws.String("test policy"),
					},
				},
			},
			listGPResp:          iam.ListGroupPoliciesOutput{},
			getGPResp:           iam.GetGroupPolicyOutput{},
			iamGroupArg:         []string{"testgroup1"},
			wantGroupPolicies:   []string(nil),
			wantGroupPolicyARNs: []string{"abcdefghijklmnopqrst"},
			wantErr:             false,
		},
		{
			description:         "No policies",
			listAGPResp:         iam.ListAttachedGroupPoliciesOutput{},
			listGPResp:          iam.ListGroupPoliciesOutput{},
			getGPResp:           iam.GetGroupPolicyOutput{},
			iamGroupArg:         []string{"testgroup1"},
			wantGroupPolicies:   []string(nil),
			wantGroupPolicyARNs: []string(nil),
			wantErr:             false,
		},
		{
			description:         "empty iam_groups arg",
			listAGPResp:         iam.ListAttachedGroupPoliciesOutput{},
			listGPResp:          iam.ListGroupPoliciesOutput{},
			getGPResp:           iam.GetGroupPolicyOutput{},
			iamGroupArg:         []string{},
			wantGroupPolicies:   []string(nil),
			wantGroupPolicyARNs: []string(nil),
			wantErr:             false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// configure backend and iam client
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}

			b := Backend()
			if err := b.Setup(context.Background(), config); err != nil {
				t.Fatal(err)
			}
			b.iamClient = &mockGroupIAMClient{
				ListAttachedGroupPoliciesResp: tc.listAGPResp,
				ListGroupPoliciesResp:         tc.listGPResp,
				GetGroupPolicyResp:            tc.getGPResp,
			}

			// run the test and compare results
			groupPolicies, groupPolicyARNs, err := b.getGroupPolicies(context.TODO(), config.StorageView, tc.iamGroupArg)
			assert.Equal(t, tc.wantGroupPolicies, groupPolicies)
			assert.Equal(t, tc.wantGroupPolicyARNs, groupPolicyARNs)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_combinePolicyDocuments(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description    string
		input          []string
		expectedOutput string
		expectedErr    bool
	}{
		{
			description: "one policy",
			input: []string{
				ec2AllPolicy,
			},
			expectedOutput: `{"Version":"2012-10-17","Statement":[{"Action":"ec2:*","Effect":"Allow","Resource":"*"}]}`,
			expectedErr:    false,
		},
		{
			description: "two policies",
			input: []string{
				ec2AllPolicy,
				ec2DescribePolicy,
			},
			expectedOutput: `{"Version": "2012-10-17", "Statement":[
				{"Action": "ec2:*", "Effect": "Allow", "Resource": "*"},
				{"Action": ["ec2:DescribeInstances"], "Effect": "Allow", "Resource": "*"}]}`,
			expectedErr: false,
		},
		{
			description: "two policies, one with empty statement",
			input: []string{
				ec2AllPolicy,
				`{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedOutput: `{"Version": "2012-10-17", "Statement": [{"Action": "ec2:*", "Effect": "Allow", "Resource": "*"}]}`,
			expectedErr:    false,
		},
		{
			description: "malformed json",
			input: []string{
				`"Version": "2012-10-17","Statement": [{"Effect": "Allow", "Action": "ec2:*", "Resource": "*"}]}`,
				`{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedOutput: ``,
			expectedErr:    true,
		},
		{
			description: "not action",
			input: []string{
				`{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "NotAction": "ec2:DescribeAvailabilityZones", "Resource": "*"}]}`,
			},
			expectedOutput: `{"Version": "2012-10-17","Statement":[{"Effect": "Allow","NotAction": "ec2:DescribeAvailabilityZones",	"Resource": "*"}]}`,
			expectedErr:    false,
		},
		{
			description: "one blank policy",
			input: []string{
				"",
				`{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedOutput: `{"Version": "2012-10-17", "Statement": []}`,
			expectedErr:    false,
		},
		{
			description: "when statement is not a list",
			input: []string{
				ec2SingleStatement,
			},
			expectedOutput: `{"Version": "2012-10-17", "Statement": [{"Action": ["ec2:DescribeInstances"], "Effect": "Allow", "Resource": "*"}]}`,
			expectedErr:    false,
		},
		{
			description: "statement is malformed json",
			input: []string{
				`{"Version": "2012-10-17", "Statement": {true}`,
			},
			expectedOutput: "",
			expectedErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			policyOut, err := combinePolicyDocuments(tc.input...)
			if (err != nil) != tc.expectedErr {
				t.Fatalf("got unexpected error: %s", err)
			}
			if (err != nil) != tc.expectedErr {
				t.Fatalf("got unexpected error: %s", err)
			}
			// remove whitespace
			if tc.expectedOutput != "" {
				tc.expectedOutput, err = compactJSON(tc.expectedOutput)
				if err != nil {
					t.Fatalf("error compacting JSON: %s", err)
				}
			}
			if policyOut != tc.expectedOutput {
				t.Fatalf("did not receive expected output: want %s, got %s", tc.expectedOutput, policyOut)
			}
		})
	}
}
