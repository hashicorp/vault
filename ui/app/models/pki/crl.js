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

  @attr('boolean') autoRebuild;
  @attr('string', {
    label: 'Auto-rebuild on',
    labelDisabled: 'Auto-rebuild off',
    mapToBoolean: 'autoRebuild',
    isOppositeValue: false,
    helperTextEnabled: 'Vault will rebuild the CRL in the below grace period before expiration',
    helperTextDisabled: 'Vault will not automatically rebuild the CRL',
  })
  autoRebuildGracePeriod;

  @attr('boolean') enableDelta;
  @attr('string', {
    label: 'Delta CRL building on',
    labelDisabled: 'Delta CRL building off',
    mapToBoolean: 'enableDelta',
    isOppositeValue: false,
    helperTextEnabled: 'Vault will rebuild the delta CRL at the interval below:',
    helperTextDisabled: 'Vault will not rebuild the delta CRL at an interval',
  })
  deltaRebuildInterval;

  @attr('boolean') disable;
  @attr('string', {
    label: 'Expiry',
    labelDisabled: 'No expiry',
    mapToBoolean: 'disable',
    isOppositeValue: true,
    helperTextDisabled: 'The CRL will not be built.',
    helperTextEnabled: 'The CRL will expire after:',
  })
  expiry;

  @attr('boolean') ocspDisable;
  @attr('string', {
    label: 'OCSP responder APIs enabled',
    labelDisabled: 'OCSP responder APIs disabled',
    mapToBoolean: 'ocspDisable',
    isOppositeValue: true,
    helperTextEnabled: "Requests about a certificate's status will be valid for:",
    helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
  })
  ocspExpiry;

  // TODO follow-on ticket to add enterprise only attributes:
  /*
  @attr('boolean') crossClusterRevocation;
  @attr('boolean') unifiedCrl;
  @attr('boolean') unifiedCrlOnExistingPaths;
  */

  @lazyCapabilities(apiPath`${'id'}/config/crl`, 'id') crlPath;

  get canSet() {
    return this.crlPath.get('canCreate') !== false;
  }
}
