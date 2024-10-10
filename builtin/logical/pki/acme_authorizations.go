// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"time"
)

type ACMEIdentifierType string

const (
	ACMEDNSIdentifier ACMEIdentifierType = "dns"
	ACMEIPIdentifier  ACMEIdentifierType = "ip"
)

type ACMEIdentifier struct {
	Type          ACMEIdentifierType `json:"type"`
	Value         string             `json:"value"`
	OriginalValue string             `json:"original_value"`
	IsWildcard    bool               `json:"is_wildcard"`
	IsV6IP        bool               `json:"is_v6_ip"`
}

func (ai *ACMEIdentifier) MaybeParseWildcard() (bool, string, error) {
	if ai.Type != ACMEDNSIdentifier || !isWildcardDomain(ai.Value) {
		return false, ai.Value, nil
	}

	// Here on out, technically it is a wildcard.
	ai.IsWildcard = true

	wildcardLabel, reducedName, err := validateWildcardDomain(ai.Value)
	if err != nil {
		return true, "", err
	}

	if wildcardLabel != "*" {
		// Per RFC 8555 Section. 7.1.3. Order Objects:
		//
		// > Any identifier of type "dns" in a newOrder request MAY have a
		// > wildcard domain name as its value.  A wildcard domain name consists
		// > of a single asterisk character followed by a single full stop
		// > character ("*.") followed by a domain name as defined for use in the
		// > Subject Alternate Name Extension by [RFC5280].
		return true, "", fmt.Errorf("wildcard must be entire left-most label")
	}

	if reducedName == "" {
		return true, "", fmt.Errorf("wildcard must not be entire domain name; need at least two domain labels")
	}

	// Parsing was indeed successful, so update our reduced name.
	ai.Value = reducedName

	return true, reducedName, nil
}

func (ai *ACMEIdentifier) NetworkMarshal(useOriginalValue bool) map[string]interface{} {
	value := ai.OriginalValue
	if !useOriginalValue {
		value = ai.Value
	}
	return map[string]interface{}{
		"type":  ai.Type,
		"value": value,
	}
}

type ACMEAuthorizationStatusType string

const (
	ACMEAuthorizationPending     ACMEAuthorizationStatusType = "pending"
	ACMEAuthorizationValid       ACMEAuthorizationStatusType = "valid"
	ACMEAuthorizationInvalid     ACMEAuthorizationStatusType = "invalid"
	ACMEAuthorizationDeactivated ACMEAuthorizationStatusType = "deactivated"
	ACMEAuthorizationExpired     ACMEAuthorizationStatusType = "expired"
	ACMEAuthorizationRevoked     ACMEAuthorizationStatusType = "revoked"
)

type ACMEOrderStatusType string

const (
	ACMEOrderPending    ACMEOrderStatusType = "pending"
	ACMEOrderProcessing ACMEOrderStatusType = "processing"
	ACMEOrderValid      ACMEOrderStatusType = "valid"
	ACMEOrderInvalid    ACMEOrderStatusType = "invalid"
	ACMEOrderReady      ACMEOrderStatusType = "ready"
)

type ACMEChallengeType string

const (
	ACMEHTTPChallenge ACMEChallengeType = "http-01"
	ACMEDNSChallenge  ACMEChallengeType = "dns-01"
	ACMEALPNChallenge ACMEChallengeType = "tls-alpn-01"
)

type ACMEChallengeStatusType string

const (
	ACMEChallengePending    ACMEChallengeStatusType = "pending"
	ACMEChallengeProcessing ACMEChallengeStatusType = "processing"
	ACMEChallengeValid      ACMEChallengeStatusType = "valid"
	ACMEChallengeInvalid    ACMEChallengeStatusType = "invalid"
)

type ACMEChallenge struct {
	Type            ACMEChallengeType       `json:"type"`
	Status          ACMEChallengeStatusType `json:"status"`
	Validated       string                  `json:"validated,optional"`
	Error           map[string]interface{}  `json:"error,optional"`
	ChallengeFields map[string]interface{}  `json:"challenge_fields"`
}

func (ac *ACMEChallenge) NetworkMarshal(acmeCtx *acmeContext, authId string) map[string]interface{} {
	resp := map[string]interface{}{
		"type":   ac.Type,
		"url":    buildChallengeUrl(acmeCtx, authId, string(ac.Type)),
		"status": ac.Status,
	}

	if ac.Validated != "" {
		resp["validated"] = ac.Validated
	}

	if len(ac.Error) > 0 {
		resp["error"] = ac.Error
	}

	for field, value := range ac.ChallengeFields {
		resp[field] = value
	}

	return resp
}

func buildChallengeUrl(acmeCtx *acmeContext, authId, challengeType string) string {
	return acmeCtx.baseUrl.JoinPath("/challenge/", authId, challengeType).String()
}

type ACMEAuthorization struct {
	Id        string `json:"id"`
	AccountId string `json:"account_id"`

	Identifier *ACMEIdentifier             `json:"identifier"`
	Status     ACMEAuthorizationStatusType `json:"status"`

	// Per RFC 8555 Section 7.1.4. Authorization Objects:
	//
	// > This field is REQUIRED for objects with "valid" in the "status"
	// > field.
	Expires string `json:"expires,optional"`

	Challenges []*ACMEChallenge `json:"challenges"`
	Wildcard   bool             `json:"wildcard"`
}

func (aa *ACMEAuthorization) GetExpires() (time.Time, error) {
	if aa.Expires == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, aa.Expires)
}

func (aa *ACMEAuthorization) NetworkMarshal(acmeCtx *acmeContext) map[string]interface{} {
	resp := map[string]interface{}{
		"identifier": aa.Identifier.NetworkMarshal( /* use value, not original value */ false),
		"status":     aa.Status,
		"wildcard":   aa.Wildcard,
	}

	if aa.Expires != "" {
		resp["expires"] = aa.Expires
	}

	if len(aa.Challenges) > 0 {
		challenges := []map[string]interface{}{}
		for _, challenge := range aa.Challenges {
			challenges = append(challenges, challenge.NetworkMarshal(acmeCtx, aa.Id))
		}
		resp["challenges"] = challenges
	}

	return resp
}
