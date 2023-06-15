/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export default interface PkiRoleAdapter extends AdapterRegistry {
  namespace: string;
  _urlForMetadata(backend: string, path: string): string;
}
