// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"testing"
)

func TestTokenPolicies(t *testing.T) {
	var s *Secret

	// Verify some of the short-circuit paths in the function
	if policies, err := s.TokenPolicies(); policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	s = &Secret{}

	if policies, err := s.TokenPolicies(); policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	s.Auth = &SecretAuth{}

	if policies, err := s.TokenPolicies(); policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	s.Auth.Policies = []string{}

	if policies, err := s.TokenPolicies(); policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	s.Auth.Policies = []string{"test"}

	if policies, err := s.TokenPolicies(); policies == nil {
		t.Error("policies was nil")
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	s.Auth = nil
	s.Data = make(map[string]interface{})

	if policies, err := s.TokenPolicies(); policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	} else if err != nil {
		t.Errorf("err was not nil, got %v", err)
	}

	// Verify that s.Data["policies"] are properly processed
	{
		policyList := make([]string, 0)
		s.Data["policies"] = policyList

		if policies, err := s.TokenPolicies(); len(policies) != len(policyList) {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}

		policyList = append(policyList, "policy1", "policy2")
		s.Data["policies"] = policyList

		if policies, err := s.TokenPolicies(); len(policyList) != 2 {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}
	}

	// Do it again but with an interface{} slice
	{
		s.Auth = nil
		policyList := make([]interface{}, 0)
		s.Data["policies"] = policyList

		if policies, err := s.TokenPolicies(); len(policies) != len(policyList) {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}

		policyItems := make([]interface{}, 2)
		policyItems[0] = "policy1"
		policyItems[1] = "policy2"

		policyList = append(policyList, policyItems...)
		s.Data["policies"] = policyList

		if policies, err := s.TokenPolicies(); len(policies) != 2 {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}

		s.Auth = nil
		s.Data["policies"] = 7.0

		if policies, err := s.TokenPolicies(); err == nil {
			t.Error("err was nil")
		} else if policies != nil {
			t.Errorf("policies was not nil, got %v", policies)
		}

		s.Auth = nil
		s.Data["policies"] = []int{2, 3, 5, 8, 13}

		if policies, err := s.TokenPolicies(); err == nil {
			t.Error("err was nil")
		} else if policies != nil {
			t.Errorf("policies was not nil, got %v", policies)
		}
	}

	s.Auth = nil
	s.Data["policies"] = nil

	if policies, err := s.TokenPolicies(); err != nil {
		t.Errorf("err was not nil, got %v", err)
	} else if policies != nil {
		t.Errorf("policies was not nil, got %v", policies)
	}

	// Verify that logic that merges s.Data["policies"] and s.Data["identity_policies"] works
	{
		policyList := []string{"policy1", "policy2", "policy3"}
		s.Data["policies"] = policyList[:1]
		s.Data["identity_policies"] = "not_a_slice"
		s.Auth = nil

		if policies, err := s.TokenPolicies(); err == nil {
			t.Error("err was nil")
		} else if policies != nil {
			t.Errorf("policies was not nil, got %v", policies)
		}

		s.Data["identity_policies"] = policyList[1:]

		if policies, err := s.TokenPolicies(); len(policyList) != len(policies) {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}
	}

	// Do it again but with an interface{} slice
	{
		policyList := []interface{}{"policy1", "policy2", "policy3"}
		s.Data["policies"] = policyList[:1]
		s.Data["identity_policies"] = "not_a_slice"
		s.Auth = nil

		if policies, err := s.TokenPolicies(); err == nil {
			t.Error("err was nil")
		} else if policies != nil {
			t.Errorf("policies was not nil, got %v", policies)
		}

		s.Data["identity_policies"] = policyList[1:]

		if policies, err := s.TokenPolicies(); len(policyList) != len(policies) {
			t.Errorf("expecting policies length %d, got %d", len(policyList), len(policies))
		} else if err != nil {
			t.Errorf("err was not nil, got %v", err)
		} else if s.Auth == nil {
			t.Error("Auth field is still nil")
		}

		s.Auth = nil
		s.Data["identity_policies"] = []int{2, 3, 5, 8, 13}

		if policies, err := s.TokenPolicies(); err == nil {
			t.Error("err was nil")
		} else if policies != nil {
			t.Errorf("policies was not nil, got %v", policies)
		}
	}

	s.Auth = nil
	s.Data["policies"] = []string{"policy1"}
	s.Data["identity_policies"] = nil

	if policies, err := s.TokenPolicies(); err != nil {
		t.Errorf("err was not nil, got %v", err)
	} else if len(policies) != 1 {
		t.Errorf("expecting policies length %d, got %d", 1, len(policies))
	} else if s.Auth == nil {
		t.Error("Auth field is still nil")
	}
}
