/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Adapter from 'ember-data/adapter';
import ModelRegistry from 'ember-data/types/registries/model';

/**
 * Catch-all for ember-data.
 */
export default interface AdapterRegistry {
  [key: keyof ModelRegistry]: Adapter;
}
