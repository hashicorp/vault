/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import SecretsEnginePathAdapter from 'vault/adapters/secrets-engine-path';

export default class LdapConfigAdapter extends SecretsEnginePathAdapter {
  path = 'config';
}
