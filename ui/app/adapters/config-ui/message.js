/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';

export default class MessageAdapter extends ApplicationAdapter {
  async query(store, type, query) {
    const { authenticated } = query;
    const url = `/v1/sys/config/ui/custom-messages`;

    try {
      return await this.ajax(url, 'GET', { data: { authenticated } });
    } catch (e) {
      if (e.httpStatus === 404) {
        return [];
      }
      throw e;
    }
  }
}
