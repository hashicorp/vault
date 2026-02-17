// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package stepwise

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/vault/api"
)

// NewAssertAuthPoliciesFunc returns a function that checks that the auth's
// policies in the response matches the expected policies.
func NewAssertAuthPoliciesFunc(policies []string) func(resp *api.Secret, err error) error {
	return func(resp *api.Secret, err error) error {
		if err != nil {
			return err
		}

		if resp == nil || resp.Auth == nil {
			return fmt.Errorf("no auth in response")
		}
		expected := make([]string, len(policies))
		copy(expected, policies)
		sort.Strings(expected)
		ret := make([]string, len(resp.Auth.Policies))
		copy(ret, resp.Auth.Policies)
		sort.Strings(ret)
		if !reflect.DeepEqual(ret, expected) {
			return fmt.Errorf("invalid policies: expected %#v, got %#v", expected, ret)
		}

		return nil
	}
}
