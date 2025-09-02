/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';
import { OidcApiResponse, SamlApiResponse } from 'vault/auth/methods';

export default interface AuthMethodAdapter extends AdapterRegistry {
  exchangeOIDC: (path: string, state: string, code: string) => Promise<OidcApiResponse>;
  pollSAMLToken: (path: string, tokenPollID: string, clientVerifier: string) => Promise<SamlApiResponse>;
}
