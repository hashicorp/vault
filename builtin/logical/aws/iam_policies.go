package aws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

// PolicyDocument represents an IAM policy document
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry represents a statement in an IAM policy document. Note that
// Action can be either a list or a string, so we represent it here as
// interface{}.
type StatementEntry struct {
	Effect   string
	Action   interface{}
	Resource string
}

// getGroupPolicies takes a list of IAM Group names and returns a list of the
// inline policy documents, and a list of the attached managed policy ARNs
func (b *backend) getGroupPolicies(iamGroups []string) (groupPolicies []string, groupPolicyARNs []string, err error) {
	var agp *iam.ListAttachedGroupPoliciesOutput
	var inlinePolicies *iam.ListGroupPoliciesOutput
	var inlinePolicyDoc *iam.GetGroupPolicyOutput

	for _, g := range iamGroups {
		// Collect managed policy ARNs from configured IAM Groups
		agp, err = b.iamClient.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return
		}
		for _, p := range agp.AttachedPolicies {
			groupPolicyARNs = append(groupPolicyARNs, *p.PolicyArn)
		}

		// Collect inline policy names from configured IAM Groups
		inlinePolicies, err = b.iamClient.ListGroupPolicies(&iam.ListGroupPoliciesInput{
			GroupName: aws.String(g),
		})
		if err != nil {
			return
		}
		for _, iP := range inlinePolicies.PolicyNames {
			inlinePolicyDoc, err = b.iamClient.GetGroupPolicy(&iam.GetGroupPolicyInput{
				GroupName:  &g,
				PolicyName: iP,
			})
			if err != nil {
				return
			}
			if inlinePolicyDoc != nil && inlinePolicyDoc.PolicyDocument != nil {
				groupPolicies = append(groupPolicies, *inlinePolicyDoc.PolicyDocument)
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
