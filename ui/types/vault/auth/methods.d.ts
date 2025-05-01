/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ApiResponse, WrapInfo } from 'vault/auth/api';
import type { POSSIBLE_FIELDS } from 'vault/utils/supported-login-methods';
import type { MfaRequirementApiResponse } from './mfa';

// ApiResponse has top-level of response with request_id, etc.
// This interface defines the "auth" key
export interface AuthResponseData extends ApiResponse {
  auth: {
    accessor: string;
    policies: string[] | null;
    metadata: Record<string, unknown> | null;
    lease_duration: number;
    renewable: boolean;
    entity_id: string;
    token_type: string;
    orphan: boolean;
    mfa_requirement: MfaRequirementApiResponse | null;
  };
}

// METHOD SPECIFIC RESPONSES
export interface OidcApiResponse extends ApiResponse {
  auth: AuthResponseData['auth'] & {
    client_token: string;
  };
}

export interface SamlApiResponse extends ApiResponse {
  auth: AuthResponseData['auth'] & {
    client_token: string;
    token_policies: string[];
    num_uses: number;
  };
}
