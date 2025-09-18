/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface MfaRequirementApiResponse {
  mfa_request_id: string;
  mfa_constraints: MfaConstraints;
}
interface MfaTotpSelfEnrollApiResponse {
  data: SelfEnrollmentData;
}

interface SelfEnrollmentData {
  barcode: string;
  url: string;
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

interface ParsedMfaRequirement {
  mfa_request_id: string;
  mfa_constraints: ParsedMfaConstraint[];
}

interface ParsedMfaConstraint {
  name: string;
  methods: ParsedMfaMethod[];
  selectedMethod: ParsedMfaMethod | null;
  passcode?: string; // DUMB
}
interface ParsedMfaMethod {
  type: string;
  id: string;
  uses_passcode: boolean;
  label: string;
  self_enrollment_enabled?: boolean;
}

interface MfaAuthData {
  mfaRequirement: ParsedMfaRequirement;
  authMethodType: string;
  authMountPath: string;
}

interface MfaConstraintState {
  methods: ParsedMfaMethod[];
  name: string;
  passcode: string;
  qrCode?: string;
  selectedMethod?: ParsedMfaMethod;
  selfEnrollMethod: ParsedMfaMethod | null;
  isSatisfied: boolean;
  setPasscode(value: string): void;
  setSelectedMethod(value: string): void;
}
