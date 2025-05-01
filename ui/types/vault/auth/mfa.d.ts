/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface MfaRequirementApiResponse {
  mfa_request_id: string;
  mfa_constraints: MfaConstraints;
}

interface MfaConstraint {
  type: string;
  id: string;
  uses_passcode: boolean;
}

interface MfaConstraints {
  [key: string]: {
    any: MfaConstraint[];
  };
}

export interface ParsedMfaRequirement {
  mfa_requirement: {
    mfa_request_id: string;
    mfa_constraints: MfaConstraint[];
  };
}

interface MfaMethod {
  type: string;
  id: string;
  uses_passcode: boolean;
  label: string;
}

interface MfaConstraint {
  name: string;
  methods: MfaMethod[];
  selectedMethod: MfaMethod;
}
