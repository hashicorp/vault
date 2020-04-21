package aws

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/vault/sdk/logical"
)

// PolicyDocument represents an IAM policy document
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry represents a statement in an IAM policy document
type StatementEntry interface{}

// getGroupPolicies takes a list of IAM Group names and returns a list of the
// inline policy documents, and a list of the attached managed policy ARNs
func (b *backend) getGroupPolicies(ctx context.Context, s logical.Storage, iamGroups []string) (groupPolicies []string, groupPolicyARNs []string, err error) {
	var agp *iam.ListAttachedGroupPoliciesOutput
	var inlinePolicies *iam.ListGroupPoliciesOutput
	var inlinePolicyDoc *iam.GetGroupPolicyOutput
	var iamClient iamiface.IAMAPI

	iamClient, err = b.clientIAM(ctx, s)
	if err != nil {
		return
	}

	for _, g := range iamGroups {
		// Collect managed policy ARNs from configured IAM Groups
		agp, err = iamClient.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return
		}
		for _, p := range agp.AttachedPolicies {
			groupPolicyARNs = append(groupPolicyARNs, *p.PolicyArn)
		}

		// Collect inline policy names from configured IAM Groups
		inlinePolicies, err = iamClient.ListGroupPolicies(&iam.ListGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return
		}
		for _, iP := range inlinePolicies.PolicyNames {
			inlinePolicyDoc, err = iamClient.GetGroupPolicy(&iam.GetGroupPolicyInput{
				GroupName:  &g,
				PolicyName: iP,
			})
			if err != nil {
				return
			}
			if inlinePolicyDoc != nil && inlinePolicyDoc.PolicyDocument != nil {
				var policyStr string
				if policyStr, err = url.QueryUnescape(*inlinePolicyDoc.PolicyDocument); err != nil {
					return
				}
				groupPolicies = append(groupPolicies, policyStr)
			}
		}
	}
	return
}

// combinePolicyDocuments takes policy strings as input, and combines them into
// a single policy document string
func combinePolicyDocuments(policies ...string) (policy string, err error) {
	var policyBytes []byte
	var newPolicy PolicyDocument = PolicyDocument{
		Version: "2012-10-17",
	}
	newPolicy.Statement = make([]StatementEntry, 0)

	for _, p := range policies {
		if len(p) == 0 {
			continue
		}
		var tmpDoc PolicyDocument
		err = json.Unmarshal([]byte(p), &tmpDoc)
		if err != nil {
			return
		}
		newPolicy.Statement = append(newPolicy.Statement, tmpDoc.Statement...)
	}

	policyBytes, err = json.Marshal(&newPolicy)
	if err != nil {
		return
	}
	policy = string(policyBytes)
	return
}
