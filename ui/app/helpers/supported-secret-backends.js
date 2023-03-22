/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

const SUPPORTED_SECRET_BACKENDS = [
  'aws',
  'database',
  'cubbyhole',
  'generic',
  'kv',
  'pki',
  'ssh',
  'transit',
  'kmip',
  'transform',
  'keymgmt',
  'kubernetes',
];

export function supportedSecretBackends() {
  return SUPPORTED_SECRET_BACKENDS;
}

export default buildHelper(supportedSecretBackends);
