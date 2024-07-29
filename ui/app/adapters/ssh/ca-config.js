/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class SshCaConfig extends ApplicationAdapter {
  namespace = 'v1';
  // For now this is only being used on the vault.cluster.secrets.backend.configuration route. This is a read-only route.
  // Eventually, this will be used to create the ca config for the SSH secret backend, replacing the requests located on the secret-engine adapter.
  queryRecord(store, type, query) {
    const { backend } = query;
    return this.ajax(`${this.buildURL()}/${encodePath(backend)}/config/ca`, 'GET').then((resp) => {
      resp.id = backend;
      return resp;
    });
  }
}
