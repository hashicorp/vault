/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SecretsEnginePathAdapter from 'vault/adapters/secrets-engine-path';

export default class KubernetesConfigAdapter extends SecretsEnginePathAdapter {
  path = 'config';

  checkConfigVars(backend) {
    return this.ajax(`${this._getURL(backend, 'check')}`, 'GET');
  }
}
