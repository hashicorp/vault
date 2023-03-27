/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Adapter from 'ember-data/adapter';
import ModelRegistry from 'ember-data/types/registries/model';
import PkiIssuerAdapter from 'vault/adapters/pki/issuer';

/**
 * Catch-all for ember-data.
 */
export default interface AdapterRegistry {
  'pki/issuer': PkiIssuerAdapter;
  [key: keyof ModelRegistry]: Adapter;
}
