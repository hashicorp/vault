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

  @attr('boolean') autoRebuild;
  @attr('string', {
    label: 'Auto-rebuild on',
    labelDisabled: 'Auto-rebuild off',
    editType: 'ttl',
    mapToBoolean: 'autoRebuild',
    helperTextEnabled: 'Vault will rebuild the CRL in the below grace period before expiration',
    helperTextDisabled: 'Vault will not automatically rebuild the CRL',
  })
  autoRebuildGracePeriod;

  @attr('boolean') enableDelta; // add validations auto_rebuild must be enabled
  @attr('string', {
    label: 'Delta CRL building on',
    labelDisabled: 'Delta CRL building off',
    editType: 'ttl',
    mapToBoolean: 'enableDelta',
    helperTextEnabled: 'Vault will rebuild the delta CRL at the interval below:',
    helperTextDisabled: 'Vault will not rebuild the delta CRL at an interval',
  })
  deltaRebuildInterval;

  @attr('boolean') disable;
  @attr('string', {
    label: 'Expiry',
    labelDisabled: 'No expiry',
    editType: 'ttl',
    mapToBoolean: 'disable',
    helperTextEnabled: 'The CRL will expire after:',
    helperTextDisabled: 'The CRL will not be built.',
  })
  expiry;

  @attr('boolean', { label: 'OCSP disable' }) ocspDisable;
  @attr('string', {
    label: 'OCSP responder APIs enabled',
    labelDisabled: 'OCSP responder APIs disabled',
    mapToBoolean: 'ocspDisable',
    helperTextEnabled: "Requests about a certificate's status will be valid for:",
    helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
  })
  ocspExpiry;

  // TODO missing from designs, enterprise only - add?
  /*
  to set cross_cluster_revocation=true or unified_crl=true
  need to have auto_rebuild=true and
  you need to have Vault Ent and this must not be a local-only mount.
  
  "cross_cluster_revocation": true,
  "unified_crl": true,
  "unified_crl_on_existing_paths": true
  */

  @lazyCapabilities(apiPath`${'id'}/config/crl`, 'id') crlPath;

  get canSet() {
    return this.crlPath.get('canCreate') !== false;
  }
}
