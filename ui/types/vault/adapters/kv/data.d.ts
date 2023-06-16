/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export default interface KvDataAdapter extends AdapterRegistry {
  namespace: string;
  _urlForSecret(backend: string, path: string, version: string | number): string;
}
