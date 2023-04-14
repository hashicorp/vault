// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"time"
)

type ACMEIdentifierType string

const (
	ACMEDNSIdentifier ACMEIdentifierType = "dns"
	ACMEIPIdentifier  ACMEIdentifierType = "ip"
)

type ACMEIdentifier struct {
	Type  ACMEIdentifierType `json:"type"`
	Value string             `json:"value"`
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
	URL             string                  `json:"url"`
	Status          ACMEChallengeStatusType `json:"status"`
	Validated       string                  `json:"validated,optional"`
	Error           map[string]interface{}  `json:"error,optional"`
	ChallengeFields map[string]interface{}  `json:"challenge_fields"`
}

func (ac *ACMEChallenge) NetworkMarshal() map[string]interface{} {
	resp := map[string]interface{}{
		"type":   ac.Type,
		"url":    ac.URL,
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

func (aa *ACMEAuthorization) NetworkMarshal() map[string]interface{} {
	resp := map[string]interface{}{
		"identifier": aa.Identifier,
		"status":     aa.Status,
		"wildcard":   aa.Wildcard,
	}

	if aa.Expires != "" {
		resp["expires"] = aa.Expires
	}

	if len(aa.Challenges) > 0 {
		challenges := []map[string]interface{}{}
		for _, challenge := range aa.Challenges {
			challenges = append(challenges, challenge.NetworkMarshal())
		}
		resp["challenges"] = challenges
	}

	return resp
}
