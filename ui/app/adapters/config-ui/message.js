/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';

export default class MessageAdapter extends ApplicationAdapter {
  pathForType() {
    return 'config/ui/custom-messages';
  }

  query(store, type, query, recordArray, adapterOptions) {
    const { authenticated } = query;
    return super.query(store, type, { authenticated, list: true }, recordArray, adapterOptions);
  }
}
