/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export default interface PkiRoleAdapter extends AdapterRegistry {
  namespace: string;
  _urlForRole(backend: string, id: string): string;
  _optionsForQuery(id: string): { data: unknown };
}
