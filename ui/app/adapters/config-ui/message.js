/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';

export default class MessageAdapter extends ApplicationAdapter {
  pathForType() {
    return 'config/ui/custom-messages';
  }

  query(store, type, query) {
    const { authenticated } = query;
    return super.query(store, type, { authenticated, list: true });
  }

  queryRecord(store, type, id) {
    return this.ajax(`${this.buildURL(type)}/${id}`, 'GET');
  }

  updateRecord(store, type, snapshot) {
    return this.ajax(`${this.buildURL(type)}/${snapshot.record.id}`, 'POST', {
      data: this.serialize(snapshot.record),
    });
  }
}
