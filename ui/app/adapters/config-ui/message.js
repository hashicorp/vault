/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class MessageAdapter extends ApplicationAdapter {
  getCustomMessagesUrl(id) {
    let url = `${this.buildURL()}/config/ui/custom-messages`;

    if (id) {
      url = url + '/' + encodePath(id);
    }
    return url;
  }

  query(store, type, query) {
    const { authenticated } = query;

    return this.ajax(this.getCustomMessagesUrl(), 'GET', { data: { authenticated, list: true } });
  }

  deleteRecord(store, type, snapshot) {
    const { id } = snapshot;
    return this.ajax(this.getCustomMessagesUrl(id), 'DELETE');
  }
}
