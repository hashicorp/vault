/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ENV from 'vault/config/environment';
const { handler } = ENV['ember-cli-mirage'];
import scenarios from './index';

export default function (server) {
  server.create('clients/config');
  server.create('feature', { feature_flags: ['SOME_FLAG', 'VAULT_CLOUD_ADMIN_NAMESPACE'] });

  if (handler in scenarios) {
    scenarios[handler](server);
  }
}
