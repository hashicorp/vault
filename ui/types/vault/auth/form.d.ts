/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface UnauthMountsByType {
  // key is the auth method type
  [key: string]: AuthTabMountData[] | null;
}

export interface AuthTabMountData {
  path: string;
  type: string;
  description?: string;
  config?: object | null;
}

export type LoginFields = Partial<Record<(typeof POSSIBLE_FIELDS)[number], string | undefined>> & {
  path?: string | undefined;
  namespace?: string | undefined;
};
