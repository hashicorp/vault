/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import { assert } from '@ember/debug';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiTidyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(snapshot) {
    const { backend } = snapshot.record;
    const { actionType } = snapshot.adapterOptions;

    if (!backend) {
      throw new Error('Backend missing');
    }

    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;

    switch (actionType) {
      case 'manual-tidy':
        return `${baseUrl}/tidy`;
      case 'auto-tidy':
        return `${baseUrl}/config/auto-tidy`;
      default:
        assert('type must be one of manual-tidy, auto-tidy');
    }
  }

  createRecord(store, type, snapshot) {
    const url = this.urlForCreateRecord(snapshot);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) });
  }
}
