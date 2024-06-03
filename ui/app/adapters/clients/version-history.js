/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';

export default class VersionHistoryAdapter extends ApplicationAdapter {
  findAll() {
    return this.ajax(this.buildURL() + '/version-history', 'GET', {
      data: {
        list: true,
      },
    }).then((resp) => {
      return resp;
    });
  }
}
