/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import ENV from 'vault/config/environment';

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
  PKI_EXTERNAL = 'pki-external-ca',
  SSH = 'ssh',
  TRANSFORM = 'transform',
  TRANSIT = 'transit',
  TOTP = 'totp',
}

export function supportedSecretBackends() {
  if (ENV.environment === 'production') {
    return Object.values(SupportedSecretBackendsEnum).filter(
      (v) => v !== SupportedSecretBackendsEnum.PKI_EXTERNAL
    );
  }
  return Object.values(SupportedSecretBackendsEnum);
}

export default buildHelper(supportedSecretBackends);
