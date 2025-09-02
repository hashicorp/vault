/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface UnauthMountsByType {
  // key is the auth method type
  // if the value is "null" there is no mount data for that type
  [key: string]: AuthTabMountData[] | null;
}
export interface UnauthMountsResponse {
  // key is the mount path
  [key: string]: { type: string; description?: string; config?: object | null };
}

interface AuthTabMountData {
  path: string;
  type: string;
  description?: string;
  config?: object | null;
}

export type LoginFields = Partial<Record<(typeof POSSIBLE_FIELDS)[number], string | undefined>> & {
  path?: string | undefined;
  namespace?: string | undefined;
};

export interface VisibleAuthMounts {
  // key is auth mount path
  [key: string]: {
    description: string;
    type: string;
    options: null | {};
  };
}

// Auth data returned from each method's login response is
// normalized so each method's information maps to the same key names
export interface NormalizedAuthData extends NormalizeAuthResponseKeys {
  authMethodType: string;
  entityId?: string;
  expireTime?: string;
  renewable: boolean;
  policies: string[];
  mfaRequirement?: MfaRequirementApiResponse | null;
}

// This info is not returned within a consistent key name so each auth method is responsible for
// normalizing it
export interface NormalizeAuthResponseKeys {
  authMountPath: string; // manually added because not a part of the auth response
  displayName?: string; // if not from the "display_name" key, then this may be set from either "meta" or "metadata"
  token: string; // was "client_token" or "id" key for some methods
  ttl: number; // was "lease_duration" key for some methods
}
