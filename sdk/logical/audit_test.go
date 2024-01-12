// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"testing"
	"time"

	"github.com/hashicorp/go-sockaddr"

	"github.com/stretchr/testify/require"
)

// TestLogInput_BexprDatum ensures that we can transform a LogInput
// into a LogInputBexpr to be used in audit filtering.
func TestLogInput_BexprDatum(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Request            *Request
		Namespace          string
		ExpectedPath       string
		ExpectedMountPoint string
		ExpectedMountType  string
		ExpectedNamespace  string
		ExpectedOperation  string
	}{
		"nil-no-namespace": {
			Request:            nil,
			Namespace:          "",
			ExpectedPath:       "",
			ExpectedMountPoint: "",
			ExpectedMountType:  "",
			ExpectedNamespace:  "",
			ExpectedOperation:  "",
		},
		"nil-namespace": {
			Request:            nil,
			Namespace:          "juan",
			ExpectedPath:       "",
			ExpectedMountPoint: "",
			ExpectedMountType:  "",
			ExpectedNamespace:  "juan",
			ExpectedOperation:  "",
		},
		"happy-path": {
			Request: &Request{
				MountPoint: "IAmAMountPoint",
				MountType:  "IAmAMountType",
				Operation:  CreateOperation,
				Path:       "IAmAPath",
			},
			Namespace:          "juan",
			ExpectedPath:       "IAmAPath",
			ExpectedMountPoint: "IAmAMountPoint",
			ExpectedMountType:  "IAmAMountType",
			ExpectedNamespace:  "juan",
			ExpectedOperation:  "create",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := &LogInput{Request: tc.Request}

			d := l.BexprDatum(tc.Namespace)

			require.Equal(t, tc.ExpectedPath, d.Path)
			require.Equal(t, tc.ExpectedMountPoint, d.MountPoint)
			require.Equal(t, tc.ExpectedMountType, d.MountType)
			require.Equal(t, tc.ExpectedNamespace, d.Namespace)
			require.Equal(t, tc.ExpectedOperation, d.Operation)
		})
	}
}

func TestLogInput_Clone(t *testing.T) {
	// input := fakeLogInput()
	// TODO: PW:
	// We need to ensure that EVERYTHING is set on the LogInput
	// that means that we always are able to get things from
	// the receivers of exported fields, particularly .Request
	// at the time of writing this seems to be the only
	// recursive ...
}

func fakeAuth() *Auth {
	return &Auth{
		LeaseOptions:              LeaseOptions{},
		InternalData:              map[string]any{},
		DisplayName:               "display-purposes",
		Policies:                  []string{},
		TokenPolicies:             []string{},
		IdentityPolicies:          []string{},
		ExternalNamespacePolicies: map[string][]string{},
		NoDefaultPolicy:           false,
		Metadata:                  map[string]string{},
		ClientToken:               "client-token",
		Accessor:                  "access-granted",
		Period:                    10 * time.Second,
		ExplicitMaxTTL:            60 * time.Second,
		NumUses:                   99,
		EntityID:                  "abc123dohreymi",
		Alias:                     &Alias{},
		GroupAliases:              []*Alias{},
		BoundCIDRs:                []*sockaddr.SockAddrMarshaler{},
		CreationPath:              "",
		TokenType:                 TokenTypeService,
		Orphan:                    false,
		PolicyResults: &PolicyResults{
			Allowed:          false,
			GrantingPolicies: []PolicyInfo{},
		},
		MFARequirement: &MFARequirement{},
		EntityCreated:  false,
	}
}

func fakeRequest() *Request {
	return &Request{
		ID:                 "qwerty-999",
		ReplicationCluster: "other-cluster",
		Operation:          "cruddy",
		Path:               "foo/bar",
		Data: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		Storage:               nil, // TODO: PW: Is that OK?
		Secret:                fakeSecret(),
		Auth:                  nil,
		Headers:               nil,
		Connection:            nil,
		ClientToken:           "",
		ClientTokenAccessor:   "",
		DisplayName:           "foo-secrets",
		MountPoint:            "foo/",
		MountType:             "kv",
		MountAccessor:         "juan-5566",
		mountRunningVersion:   "55.1",
		mountRunningSha256:    "98jf9j38gj390jg904e5",
		mountIsExternalPlugin: false,
		mountClass:            "secret",
		WrapInfo: &RequestWrapInfo{
			TTL:      1 * time.Minute,
			Format:   "json",
			SealWrap: false,
		},
		ClientTokenRemainingUses: 3,
		EntityID:                 "45678ijhgvbnjuyt",
		PolicyOverride:           false,
		Unauthenticated:          false,
		MFACreds:                 MFACreds{},
		tokenEntry: &TokenEntry{
			Type:                     0,
			ID:                       "",
			ExternalID:               "",
			Accessor:                 "",
			Parent:                   "",
			Policies:                 nil,
			InlinePolicy:             "",
			Path:                     "",
			Meta:                     nil,
			InternalMeta:             nil,
			DisplayName:              "",
			NumUses:                  0,
			CreationTime:             0,
			TTL:                      0,
			ExplicitMaxTTL:           0,
			Role:                     "",
			Period:                   0,
			DisplayNameDeprecated:    "",
			NumUsesDeprecated:        0,
			CreationTimeDeprecated:   0,
			ExplicitMaxTTLDeprecated: 0,
			EntityID:                 "",
			NoIdentityPolicies:       false,
			BoundCIDRs:               nil,
			NamespaceID:              "",
			CubbyholeID:              "",
		},
		lastRemoteWAL:     99,
		ControlGroup:      nil,
		ClientTokenSource: 0,
		HTTPRequest:       nil,
		ResponseWriter:    nil,
		requiredState:     nil,
		responseState:     nil,
		ClientID:          "",
		InboundSSCToken:   "",
		ForwardedFrom:     "",
		ChrootNamespace:   "",
	}
}

func fakeResponse() *Response {
	return &Response{
		Secret:    nil,
		Auth:      fakeAuth(),
		Data:      nil,
		Redirect:  "",
		Warnings:  nil,
		WrapInfo:  nil,
		Headers:   nil,
		MountType: "",
	}
}

func fakeSecret() *Secret {
	return &Secret{
		LeaseOptions: LeaseOptions{
			TTL:       3 * time.Second,
			MaxTTL:    30 * time.Second,
			Renewable: false,
			Increment: 0,
			IssueTime: time.Now(),
		},
		InternalData: map[string]any{
			"secretKey1": "secretValue1",
			"secretKey2": "secretValue2",
		},
		LeaseID: "qazxswedcvfrtgbnhy",
	}
}

func fakeLogInput() *LogInput {
	return &LogInput{
		Type:                "my-type",
		Auth:                fakeAuth(),
		Request:             fakeRequest(),
		Response:            fakeResponse(),
		OuterErr:            nil,
		NonHMACReqDataKeys:  nil,
		NonHMACRespDataKeys: nil,
	}
}
