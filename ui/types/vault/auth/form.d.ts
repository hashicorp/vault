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
