/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ApiResponse, WrapInfo } from 'vault/auth/api';
import type { POSSIBLE_FIELDS } from 'vault/utils/auth-form-helpers';
import type { MfaRequirementApiResponse } from './mfa';

// ApiResponse includes top-level fields like request_id, etc.
// Some auth methods return login data under the "auth" key,
// while token exchange flows return it under the "data" key.
// The structure of the returned data varies slightly between these cases.
interface SharedAuthResponseData {
  accessor: string;
  entityId: string;
  policies: string[];
  renewable: boolean;
}

// AuthResponseAuthKey defines login data inside the "auth" key
interface AuthResponseAuthKey extends SharedAuthResponseData {
  clientToken: string;
  leaseDuration: number;
  metadata: Record<string, string>;
  mfaRequirement: MfaRequirementApiResponse | null;
  tokenType: 'service' | 'batch';
}

// AuthResponseDataKey defines login data inside the "data" key
interface AuthResponseDataKey extends SharedAuthResponseData {
  displayName: string;
  expireTime: string;
  id: string; // this is the Vault issued token (the equivalent of the clientToken for responses with the "auth" key)
  meta: Record<string, string> | null;
  namespacePath?: string;
  ttl: number;
  type: 'service' | 'batch'; // token type
}

// METHOD SPECIFIC RESPONSES
export interface GithubLoginApiResponse extends ApiResponse {
  auth: AuthResponseAuthKey & {
    metadata: { org: string; username: string };
  };
}

export interface JwtOidcLoginApiResponse extends ApiResponse {
  auth: AuthResponseAuthKey;
}

export interface OidcApiResponse extends ApiResponse {
  auth: AuthResponseData['auth'] & {
    client_token: string;
  };
}

export interface JwtOidcAuthUrlResponse extends ApiResponse {
  data: { authUrl: string };
}

export interface OktaVerifyApiResponse extends ApiResponse {
  data: { correctAnswer: number };
}

export interface SamlLoginApiResponse extends ApiResponse {
  auth: AuthResponseAuthKey;
}

export interface SamlSsoServiceUrlApiResponse extends ApiResponse {
  data: {
    ssoServiceUrl: string;
    tokenPollId: string;
  };
}

export interface TokenLoginApiResponse extends ApiResponse {
  data: AuthResponseDataKey;
}

// auth types: ldap, okta, radius, userpass
export interface UsernameLoginResponse extends ApiResponse {
  auth: AuthResponseAuthKey & {
    metadata: { username: string };
  };
}
