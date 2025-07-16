// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_encodeIdentifierForHTTP01Challenge validates the encoding behaviors of our identifiers
// for the HTTP01 challenge. Basically properly encode the identifier for an HTTP request.
func Test_encodeIdentifierForHTTP01Challenge(t *testing.T) {
	tests := []struct {
		name string
		arg  *ACMEIdentifier
		want string
	}{
		{
			name: "dns",
			arg:  &ACMEIdentifier{Type: ACMEDNSIdentifier, Value: "www.dadgarcorp.com"},
			want: "www.dadgarcorp.com",
		},
		{
			name: "ipv4",
			arg:  &ACMEIdentifier{Type: ACMEIPIdentifier, Value: "192.168.1.1"},
			want: "192.168.1.1",
		},
		{
			name: "ipv6",
			arg:  &ACMEIdentifier{Type: ACMEIPIdentifier, Value: "2001:0db8:0000:0000:0000:0000:0000:0068", IsV6IP: true},
			want: "[2001:0db8:0000:0000:0000:0000:0000:0068]",
		},
		{
			name: "ipv6-zoned",
			arg:  &ACMEIdentifier{Type: ACMEIPIdentifier, Value: "fe80::1cc0:3e8c:119f:c2e1%ens18", IsV6IP: true},
			want: "[fe80::1cc0:3e8c:119f:c2e1%25ens18]",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, encodeIdentifierForHTTP01Challenge(tt.arg), "encodeIdentifierForHTTP01Challenge(%v)", tt.arg)
		})
	}
}
