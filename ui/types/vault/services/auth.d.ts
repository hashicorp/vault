/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// temporary interface for auth service until it can be updated to ts
// add properties as needed

import Service from '@ember/service';
import type { MfaRequirementApiResponse, ParsedMfaRequirement } from 'vault/auth/mfa';
import type { NormalizedAuthData, NormalizeAuthResponseKeys } from 'vault/auth/form';
import type { AuthResponseAuthKey, AuthResponseDataKey } from 'vault/auth/methods';

interface AuthData {
  userRootNamespace: string;
  token: string;
  policies: string[];
  renewable: boolean;
  entityId: string;
  displayName?: string;
}

export interface AuthSuccessResponse {
  namespace: string;
  token: string; // the name of the token in local storage, not the actual token
  isRoot: boolean;
}
export interface AuthResponseWithMfa {
  mfa_requirement: MfaRequirementApiResponse;
}

type NormalizedProps = NormalizeAuthResponseKeys & {
  authMethodType: string;
};

export default class AuthService extends Service {
  authData: AuthData;
  currentToken: string;
  mfaErrors: null | Errors[];
  setLastFetch: (time: number) => void;
  authSuccess(clusterId: string, authData: NormalizedAuthData): Promise<AuthSuccessResponse>;
  ajax: (
    url: string,
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    options?: {
      headers?: Record<string, string>;
      namespace?: string;
      data?: Record<string, unknown>;
    }
  ) => Promise<any>;
  getAuthType(): string | undefined;
  parseMfaResponse(mfaResponse: MfaRequirementApiResponse): ParsedMfaRequirement;
  normalizeAuthData(
    authResponse: AuthResponseAuthKey | AuthResponseDataKey,
    normalizedProperties: NormalizedProps
  ): NormalizedAuthData;
  totpValidate({
    clusterId: string,
    mfaRequirement: ParsedMfaRequirement,
    authMethodType: string,
    authMountPath: string,
  }): Promise<AuthSuccessResponse>;
}
