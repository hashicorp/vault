/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

export default interface PkiTidyAdapter extends AdapterRegistry {
  namespace: string;
  cancelTidy(backend: string);
}
