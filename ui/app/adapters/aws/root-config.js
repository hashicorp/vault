/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class AwsRootConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord(store, type, query) {
    return this.ajax(`/v1/${encodePath(query.backend)}/config/root`, 'GET').then((resp) => {
      resp.id = query.backend;
      return resp;
    });
  }
}
