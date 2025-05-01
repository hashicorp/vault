/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// temporary interface for auth service until it can be updated to ts
// add properties as needed

import Service from '@ember/service';
import type { MfaRequirementApiResponse, ParsedMfaRequirement } from 'vault/auth/mfa';

export interface AuthData {
  userRootNamespace: string;
  token: string;
  policies: string[];
  renewable: boolean;
  entity_id: string;
  displayName?: string;
  mfa_requirement?: MfaRequirementApiResponse;
}

export default class AuthService extends Service {
  authData: AuthData;
  currentToken: string;
  mfaErrors: null | Errors[];
  setLastFetch: (time: number) => void;
  handleError: (error: Error | string) => string | error[] | [error];
  authenticate(params: {
    clusterId: string;
    backend: string;
    data: object;
    selectedAuth: string;
  }): Promise<any>;
  ajax: (
    url: string,
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    options?: {
      headers?: Record<string, string>;
      namespace?: string;
      data?: Record<string, unknown>;
    }
  ) => Promise<any>;
  _parseMfaResponse(mfaResponse: MfaRequirementApiResponse): ParsedMfaRequirement;
}
