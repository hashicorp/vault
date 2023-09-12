/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export default interface LdapRoleAdapter extends AdapterRegistry {
  fetchCredentials(backend: string, type: string, name: string);
  rotateStaticPassword(backend: string, name: string);
}
