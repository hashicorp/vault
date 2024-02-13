/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// temporary interface for auth service until it can be updated to ts
// add properties as needed

import Service from '@ember/service';

export interface AuthData {
  userRootNamespace: string;
  token: string;
  policies: string[];
  renewable: boolean;
  entity_id: string;
  displayName?: string;
}

export default class AuthService extends Service {
  authData: AuthData;
}
