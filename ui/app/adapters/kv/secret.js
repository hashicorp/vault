/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KvSecretAdapter extends ApplicationAdapter {
  namespace = 'v1';

  getURL(backend, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/metadata/`;
    return name ? `${base}/${name}` : base;
  }

  query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.getURL(backend), 'GET', { data: { list: true } });
  }
}
