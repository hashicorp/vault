/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const formFieldGroups = [
  {
    'Certificate Revocation List (CRL)': ['expiry', 'autoRebuildGracePeriod', 'deltaRebuildInterval'],
  },
  {
    'Online Certificate Status Protocol (OCSP)': ['ocspExpiry'],
  },
  { 'Unified Revocation': ['crossClusterRevocation', 'unifiedCrl', 'unifiedCrlOnExistingPaths'] },
];
@withFormFields(null, formFieldGroups)
export default class PkiConfigCrlModel extends Model {
  // This model uses the backend value as the model ID

  @attr('boolean') autoRebuild;
  @attr('string', {
    label: 'Auto-rebuild on',
    labelDisabled: 'Auto-rebuild off',
    mapToBoolean: 'autoRebuild',
    isOppositeValue: false,
    editType: 'ttl',
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
    editType: 'ttl',
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
    editType: 'ttl',
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
    editType: 'ttl',
    helperTextEnabled: "Requests about a certificate's status will be valid for:",
    helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
  })
  ocspExpiry;

  // enterprise only params
  @attr('boolean', {
    label: 'Cross-cluster revocation',
    helpText:
      'Enables cross-cluster revocation request queues. When a serial not issued on this local cluster is passed to the /revoke endpoint, it is replicated across clusters and revoked by the issuing cluster if it is online.',
  })
  crossClusterRevocation;

  @attr('boolean', {
    label: 'Unified CRL',
    helpText:
      'Enables unified CRL and OCSP building. This synchronizes all revocations between clusters; a single, unified CRL will be built on the active node of the primary performance replication (PR) cluster.',
  })
  unifiedCrl;

  @attr('boolean', {
    label: 'Unified CRL on existing paths',
    helpText:
      'If enabled, existing CRL and OCSP paths will return the unified CRL instead of a response based on cluster-local data.',
  })
  unifiedCrlOnExistingPaths;

  @lazyCapabilities(apiPath`${'id'}/config/crl`, 'id') crlPath;

  get canSet() {
    return this.crlPath.get('canUpdate') !== false;
  }
}
