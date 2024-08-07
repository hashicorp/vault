/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
export default class KvConfigAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForFindRecord(id) {
    return `${this.buildURL()}/${id}/config`;
  }
}
