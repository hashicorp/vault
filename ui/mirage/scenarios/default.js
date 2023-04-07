/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ENV from 'vault/config/environment';
const { handler } = ENV['ember-cli-mirage'];
import kubernetesScenario from './kubernetes';

export default function (server) {
  server.create('clients/config');
  server.create('feature', { feature_flags: ['SOME_FLAG', 'VAULT_CLOUD_ADMIN_NAMESPACE'] });

  if (handler === 'kubernetes') {
    kubernetesScenario(server);
  }
}
