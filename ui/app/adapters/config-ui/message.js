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

  async query(store, type, query) {
    const { authenticated } = query;

    try {
      return await this.ajax(this.getCustomMessagesUrl(), 'GET', { data: { authenticated } });
    } catch (e) {
      if (e.httpStatus === 404) {
        return [];
      }
      throw e;
    }
  }
}
