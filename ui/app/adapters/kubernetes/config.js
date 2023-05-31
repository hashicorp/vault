/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import SecretsEnginePathAdapter from 'vault/adapters/secrets-engine-path';

export default class KubernetesConfigAdapter extends SecretsEnginePathAdapter {
  path = 'config';

  checkConfigVars(backend) {
    return this.ajax(`${this._getURL(backend, 'check')}`, 'GET');
  }
}
