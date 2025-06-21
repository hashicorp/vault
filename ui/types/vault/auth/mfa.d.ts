/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface MfaRequirementApiResponse {
  mfaRequestId: string;
  mfaConstraints: MfaConstraints;
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
  mfaRequirement: {
    mfaRequestId: string;
    mfaConstraints: MfaConstraint[];
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
