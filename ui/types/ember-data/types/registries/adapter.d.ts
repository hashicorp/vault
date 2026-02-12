/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from 'vault/adapters/application';
import Adapter from 'ember-data/adapter';
import ModelRegistry from 'ember-data/types/registries/model';
/**
 * Catch-all for ember-data.
 */
export default interface AdapterRegistry {
  application: Application;
  [key: keyof ModelRegistry]: Adapter;
}
