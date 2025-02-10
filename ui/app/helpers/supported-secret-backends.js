/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

const SUPPORTED_SECRET_BACKENDS = [
  'aws',
  'azure',
  'cubbyhole',
  'database',
  'gcp',
  'generic',
  'keymgmt',
  'kmip',
  'kubernetes',
  'kv',
  'ldap',
  'pki',
  'ssh',
  'transform',
  'transit',
];

export function supportedSecretBackends() {
  return SUPPORTED_SECRET_BACKENDS;
}

export default buildHelper(supportedSecretBackends);
