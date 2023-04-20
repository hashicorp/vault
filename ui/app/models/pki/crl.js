/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields(['crlExpiryData', 'autoRebuildData', 'deltaCrlBuildingData', 'ocspExpiryData'])
export default class PkiCrlModel extends Model {
  // This model uses the backend value as the model ID

  @attr('object', {
    label: 'Auto-rebuild on',
    labelDisabled: 'Auto-rebuild off',
    editType: 'ttl',
    defaultValue() {
      return { enabled: false, duration: '12h' };
    },
    helperTextEnabled: 'Vault will rebuild the CRL in the below grace period before expiration',
    helperTextDisabled: 'Vault will not automatically rebuild the CRL',
  })
  autoRebuildData; // sets auto_rebuild (boolean), auto_rebuild_grace_period (duration)

  @attr('object', {
    label: 'Delta CRL building on',
    labelDisabled: 'Delta CRL building off',
    editType: 'ttl',
    defaultValue() {
      return { enabled: false, duration: '15m' };
    },
    helperTextEnabled: 'Vault will rebuild the delta CRL at the interval below:',
    helperTextDisabled: 'Vault will not rebuild the delta CRL at an interval',
  })
  deltaCrlBuildingData; // sets enable_delta (boolean), delta_rebuild_interval (duration)

  @attr('object', {
    label: 'Expiry',
    labelDisabled: 'No expiry',
    editType: 'ttl',
    defaultValue() {
      return { enabled: true, duration: '72h' };
    },
    helperTextEnabled: 'The CRL will expire after:',
    helperTextDisabled: 'The CRL will not be built.',
  })
  crlExpiryData; // sets disable (boolean), expiry (duration)

  @attr('object', {
    label: 'OCSP responder APIs enabled',
    labelDisabled: 'OCSP responder APIs disabled',
    defaultValue() {
      return { enabled: true, duration: '12h' };
    },
    helperTextEnabled: "Requests about a certificate's status will be valid for:",
    helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
  })
  ocspExpiryData; // sets ocsp_disable (boolean), ocsp_expiry (duration)

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
