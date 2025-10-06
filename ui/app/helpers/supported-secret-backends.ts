/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export enum SupportedSecretBackendsEnum {
  AWS = 'aws',
  AZURE = 'azure',
  CUBBYHOLE = 'cubbyhole',
  DATABASE = 'database',
  GCP = 'gcp',
  GENERIC = 'generic',
  KEYMGMT = 'keymgmt',
  KMIP = 'kmip',
  KUBERNETES = 'kubernetes',
  KV = 'kv',
  LDAP = 'ldap',
  PKI = 'pki',
  SSH = 'ssh',
  TRANSFORM = 'transform',
  TRANSIT = 'transit',
  TOTP = 'totp',
}

export function supportedSecretBackends() {
  return Object.values(SupportedSecretBackendsEnum);
}

export default buildHelper(supportedSecretBackends);
