/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class AwsRootConfig extends ApplicationAdapter {
  namespace = 'v1';
  // For now this is only being used on the vault.cluster.secrets.backend.configuration route. This is a read-only route.
  // Eventually, this will be used to create the root config for the AWS secret backend, replacing the requests located on the secret-engine adapter.
  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(`/v1/${encodePath(backend)}/config/root`, 'GET').then((resp) => {
      resp.id = backend;
      return resp;
    });
  }
}
