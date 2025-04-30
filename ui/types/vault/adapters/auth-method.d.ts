/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';
import type { AuthData } from 'vault/services/auth';

export default interface AuthMethodAdapter extends AdapterRegistry {
  exchangeOIDC: (
    path: string,
    state: string,
    code: string
  ) => Promise<{
    auth: AuthData;
  }>;
  pollSAMLToken: (
    path: string,
    tokenPollID: string,
    clientVerifier: string
  ) => Promise<{
    auth: AuthData;
  }>;
}
