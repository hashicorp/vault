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
  entity_id: string;
  policies: string[];
  renewable: boolean;
}

// AuthResponseAuthKey defines login data inside the "auth" key
interface AuthResponseAuthKey extends SharedAuthResponseData {
  client_token: string;
  lease_duration: number;
  metadata: Record<string, string>;
  mfa_requirement: MfaRequirementApiResponse | null;
  token_type: 'service' | 'batch';
}

// AuthResponseDataKey defines login data inside the "data" key
interface AuthResponseDataKey extends SharedAuthResponseData {
  display_name: string;
  expire_time: string;
  id: string; // this is the Vault issued token (the equivalent of the clientToken for responses with the "auth" key)
  meta: Record<string, string> | null;
  namespace_path?: string;
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
  data: { auth_url: string };
}

export interface OktaVerifyApiResponse extends ApiResponse {
  data: { correct_answer: number };
}

export interface SamlLoginApiResponse extends ApiResponse {
  auth: AuthResponseAuthKey;
}

export interface SamlSsoServiceUrlApiResponse extends ApiResponse {
  data: {
    sso_service_url: string;
    token_poll_id: string;
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
