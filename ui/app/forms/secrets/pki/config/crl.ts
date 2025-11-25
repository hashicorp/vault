/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { PkiConfigureCrlRequest } from '@hashicorp/vault-client-typescript';

export default class PkiConfigCrlForm extends Form<PkiConfigureCrlRequest> {
  crlFields = [
    new FormField('expiry', 'string', {
      label: 'Expiry',
      labelDisabled: 'No expiry',
      mapToBoolean: 'disable',
      isOppositeValue: true,
      editType: 'ttl',
      helperTextDisabled: 'The CRL will not be built.',
      helperTextEnabled: 'The CRL will expire after:',
    }),
    new FormField('auto_rebuild_grace_period', 'string', {
      label: 'Auto-rebuild on',
      labelDisabled: 'Auto-rebuild off',
      mapToBoolean: 'auto_rebuild',
      isOppositeValue: false,
      editType: 'ttl',
      helperTextEnabled: 'Vault will rebuild the CRL in the below grace period before expiration',
      helperTextDisabled: 'Vault will not automatically rebuild the CRL',
    }),
    new FormField('delta_rebuild_interval', 'string', {
      label: 'Delta CRL building on',
      labelDisabled: 'Delta CRL building off',
      mapToBoolean: 'enable_delta',
      isOppositeValue: false,
      editType: 'ttl',
      helperTextEnabled: 'Vault will rebuild the delta CRL at the interval below:',
      helperTextDisabled: 'Vault will not rebuild the delta CRL at an interval',
    }),
  ];

  ocspFields = [
    new FormField('ocsp_expiry', 'string', {
      label: 'OCSP responder APIs enabled',
      labelDisabled: 'OCSP responder APIs disabled',
      mapToBoolean: 'ocsp_disable',
      isOppositeValue: true,
      editType: 'ttl',
      helperTextEnabled: "Requests about a certificate's status will be valid for:",
      helperTextDisabled: 'Requests cannot be made to check if an individual certificate is valid.',
    }),
  ];

  revocationFields = [
    new FormField('cross_cluster_revocation', 'boolean', {
      label: 'Cross-cluster revocation',
      helpText:
        'Enables cross-cluster revocation request queues. When a serial not issued on this local cluster is passed to the /revoke endpoint, it is replicated across clusters and revoked by the issuing cluster if it is online.',
    }),
    new FormField('unified_crl', 'boolean', {
      label: 'Unified CRL',
      helpText:
        'Enables unified CRL and OCSP building. This synchronizes all revocations between clusters; a single, unified CRL will be built on the active node of the primary performance replication (PR) cluster.',
    }),
    new FormField('unified_crl_on_existing_paths', 'boolean', {
      label: 'Unified CRL on existing paths',
      helpText:
        'If enabled, existing CRL and OCSP paths will return the unified CRL instead of a response based on cluster-local data.',
    }),
  ];

  formFields = [
    new FormField('auto_rebuild', 'boolean'),
    new FormField('enable_delta', 'boolean'),
    new FormField('disable', 'boolean'),
    new FormField('ocsp_disable', 'boolean'),
    ...this.crlFields,
    ...this.ocspFields,
    ...this.revocationFields,
  ];

  formFieldGroups = [
    new FormFieldGroup('Certificate Revocation List (CRL)', this.crlFields),
    new FormFieldGroup('Online Certificate Status Protocol (OCSP)', this.ocspFields),
    new FormFieldGroup('Unified Revocation', this.revocationFields),
  ];
}
