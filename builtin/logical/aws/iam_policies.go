// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/vault/sdk/logical"
)

// PolicyDocument represents an IAM policy document
type PolicyDocument struct {
	Version    string           `json:"Version"`
	Statements StatementEntries `json:"Statement"`
}

// StatementEntries is a slice of statements that make up a PolicyDocument
type StatementEntries []interface{}

// UnmarshalJSON is defined here for StatementEntries because the Statement
// portion of an IAM Policy can either be a list or a single element, so if it's
// a single element this wraps it in a []interface{} so that it's easy to
// combine with other policy statements:
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_statement.html
func (se *StatementEntries) UnmarshalJSON(b []byte) error {
	var out StatementEntries

	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	switch t := data.(type) {
	case []interface{}:
		out = t
	case interface{}:
		out = []interface{}{t}
	default:
		return fmt.Errorf("unsupported data type %T for StatementEntries", t)
	}
	*se = out
	return nil
}

// getGroupPolicies takes a list of IAM Group names and returns a list of their
// inline policy documents, and a list of the attached managed policy ARNs
func (b *backend) getGroupPolicies(ctx context.Context, s logical.Storage, iamGroups []string) ([]string, []string, error) {
	var groupPolicies []string
	var groupPolicyARNs []string
	var err error
	var agp *iam.ListAttachedGroupPoliciesOutput
	var inlinePolicies *iam.ListGroupPoliciesOutput
	var inlinePolicyDoc *iam.GetGroupPolicyOutput
	var iamClient iamiface.IAMAPI

	// Return early if there are no groups, to avoid creating an IAM client
	// needlessly
	if len(iamGroups) == 0 {
		return nil, nil, nil
	}

	iamClient, err = b.clientIAM(ctx, s, nil)
	if err != nil {
		return nil, nil, err
	}

	for _, g := range iamGroups {
		// Collect managed policy ARNs from the IAM Group
		agp, err = iamClient.ListAttachedGroupPoliciesWithContext(ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return nil, nil, err
		}
		for _, p := range agp.AttachedPolicies {
			groupPolicyARNs = append(groupPolicyARNs, *p.PolicyArn)
		}

		// Collect inline policy names from the IAM Group
		inlinePolicies, err = iamClient.ListGroupPoliciesWithContext(ctx, &iam.ListGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return nil, nil, err
		}
		for _, iP := range inlinePolicies.PolicyNames {
			inlinePolicyDoc, err = iamClient.GetGroupPolicyWithContext(ctx, &iam.GetGroupPolicyInput{
				GroupName:  &g,
				PolicyName: iP,
			})
			if err != nil {
				return nil, nil, err
			}
			if inlinePolicyDoc != nil && inlinePolicyDoc.PolicyDocument != nil {
				var policyStr string
				if policyStr, err = url.QueryUnescape(*inlinePolicyDoc.PolicyDocument); err != nil {
					return nil, nil, err
				}
				groupPolicies = append(groupPolicies, policyStr)
			}
		}
	}
	return groupPolicies, groupPolicyARNs, nil
}

// combinePolicyDocuments takes policy strings as input, and combines them into
// a single policy document string
func combinePolicyDocuments(policies ...string) (string, error) {
	var policy string
	var err error
	var policyBytes []byte
	newPolicy := PolicyDocument{
		// 2012-10-17 is the current version of the AWS policy language:
		// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html
		Version: "2012-10-17",
	}
	newPolicy.Statements = make(StatementEntries, 0, len(policies))

	for _, p := range policies {
		if len(p) == 0 {
			continue
		}
		var tmpDoc PolicyDocument
		err = json.Unmarshal([]byte(p), &tmpDoc)
		if err != nil {
			return "", err
		}
		newPolicy.Statements = append(newPolicy.Statements, tmpDoc.Statements...)
	}

	policyBytes, err = json.Marshal(&newPolicy)
	if err != nil {
		return "", err
	}
	policy = string(policyBytes)
	return policy, nil
}
