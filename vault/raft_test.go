// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
)

// TestFormatDiscoveredAddr validates that the string returned by formatDiscoveredAddr always respect the format `host:port`.
func TestFormatDiscoveredAddr(t *testing.T) {
	type TestCase struct {
		addr string
		port uint
		res  string
	}
	cases := []TestCase{
		{addr: "127.0.0.1", port: uint(8200), res: "127.0.0.1:8200"},
		{addr: "192.168.137.1:8201", port: uint(8200), res: "192.168.137.1:8201"},
		{addr: "fe80::aa5e:45ff:fe54:c6ce", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8200"},
		{addr: "::1", port: uint(8200), res: "[::1]:8200"},
		{addr: "[::1]", port: uint(8200), res: "[::1]:8200"},
		{addr: "[::1]:8201", port: uint(8200), res: "[::1]:8201"},
		{addr: "[fe80::aa5e:45ff:fe54:c6ce]", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8200"},
		{addr: "[fe80::aa5e:45ff:fe54:c6ce]:8201", port: uint(8200), res: "[fe80::aa5e:45ff:fe54:c6ce]:8201"},
	}
	for i, c := range cases {
		res := formatDiscoveredAddr(c.addr, c.port)
		if res != c.res {
			t.Errorf("case %d result shoud be \"%s\" but is \"%s\"", i, c.res, res)
		}
	}
}
