/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields(['expiry', 'autoRebuildGracePeriod', 'deltaRebuildInterval', 'ocspExpiry'])
export default class PkiCrlModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backendPath) {
    return `/v1/${backendPath}/config/crl?help=1`;
  }

  @attr('boolean', { defaultValue: true })
  autoRebuild;

  @attr('string', {
    label: 'Auto-rebuild on',
    labelDisabled: 'Auto-rebuild off',
    defaultValue: '12h',
    editType: 'ttl',
    booleanBuddy: 'autoRebuild',
    helperTextEnabled: 'Vault will rebuild the CRL in the below grace period before expiration',
    helperTextDisabled: 'Vault will not automatically rebuild the CRL',
  })
  autoRebuildGracePeriod;

  @attr('boolean', { defaultValue: true })
  enableDelta;

  @attr('string', {
    label: 'Expiry',
    labelDisabled: 'No expiry',
    defaultValue: '72h',
    editType: 'ttl',
    booleanBuddy: 'disable',
    helperTextEnabled: 'The CRL will expire after:',
    helperTextDisabled: 'The CRL will not be built.',
  })
  expiry;

  @attr('string', {
    label: 'Delta CRL building on',
    labelDisabled: 'Delta CRL building off',
    defaultValue: '15mh',
    editType: 'ttl',
    booleanBuddy: 'enableDelta',
    helperTextEnabled: 'Vault will rebuild the delta CRL at the interval below:',
    helperTextDisabled: 'Vault will not rebuild the delta CRL at an interval',
  })
  deltaRebuildInterval;

  @attr('boolean', { defaultValue: true })
  disable;

  @attr('string', {
    label: 'OCSP responder APIs enabled',
    labelDisabled: 'OCSP responder APIs disabled',
    defaultValue: '1h',
    booleanBuddy: 'ocspDisable',
    helperTextEnabled: "Requests about a certificate's status will be valid for:",
    helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
  })
  ocspExpiry;

  @attr('boolean', { label: 'OCSP disable', defaultValue: true })
  ocspDisable;

  // TODO missing from designs, enterprise only - add?
  /*
  "cross_cluster_revocation": true,
  "unified_crl": true,
  "unified_crl_on_existing_paths": true
  */

  @lazyCapabilities(apiPath`${'id'}/config/crl`, 'id') crlPath;

  get canSet() {
    return this.crlPath.get('canCreate') !== false;
  }
}
