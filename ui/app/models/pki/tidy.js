/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { service } from '@ember/service';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

@withExpandedAttributes()
export default class PkiTidyModel extends Model {
  // the backend mount is the model id, only one pki/tidy model will ever persist (the auto-tidy config)
  @service version;

  @attr({
    label: 'Tidy ACME enabled',
    labelDisabled: 'Tidy ACME disabled',
    mapToBoolean: 'tidyAcme',
    helperTextDisabled: 'Tidying of ACME accounts, orders and authorizations is disabled.',
    helperTextEnabled:
      'The amount of time that must pass after creation that an account with no orders is marked revoked, and the amount of time after being marked revoked or deactivated.',
    detailsLabel: 'ACME account safety buffer',
    formatTtl: true,
    defaultValue: '720h',
  })
  acmeAccountSafetyBuffer;

  @attr('boolean', {
    label: 'Tidy ACME',
    defaultValue: false,
  })
  tidyAcme;

  // * auto tidy only fields
  @attr('boolean', {
    label: 'Automatic tidy enabled',
    labelDisabled: 'Automatic tidy disabled',
    helperTextDisabled: 'Automatic tidy operations will not run.',
  })
  enabled; // renders outside FormField loop as a toggle, auto tidy fields only render if enabled

  @attr({
    editType: 'ttl',
    helperTextEnabled:
      'Sets the interval_duration between automatic tidy operations; note that this is from the end of one operation to the start of the next.',
    hideToggle: true,
    formatTtl: true,
  })
  intervalDuration;

  @attr({
    label: 'Minimum startup backoff duration',
    editType: 'ttl',
    helperTextEnabled:
      'Sets the min_startup_backoff_duration field which forces the minimum delay after Vault startup auto-tidy can run.',
    hideToggle: true,
    formatTtl: true,
  })
  minStartupBackoffDuration;

  @attr({
    label: 'Maximum startup backoff duration',
    editType: 'ttl',
    helperTextEnabled:
      'Sets the max_startup_backoff_duration field which forces the maximum delay after Vault startup auto-tidy can run.',
    hideToggle: true,
    formatTtl: true,
  })
  maxStartupBackoffDuration;
  // * end of auto-tidy only fields

  @attr({
    editType: 'ttl',
    helperTextEnabled:
      'Specifies a duration that issuers should be kept for, past their NotAfter validity period. Defaults to 365 days (8760 hours).',
    hideToggle: true,
    formatTtl: true,
  })
  issuerSafetyBuffer;

  @attr('string', {
    editType: 'ttl',
    helperTextEnabled:
      'Specifies the duration to pause between tidying individual certificates. This releases the revocation lock and allows other operations to continue while tidy is running.',
    hideToggle: true,
    formatTtl: true,
  })
  pauseDuration;

  @attr('string', {
    editType: 'ttl',
    helperTextEnabled:
      'Specifies a duration after which cross-cluster revocation requests will be removed as expired.',
    hideToggle: true,
    formatTtl: true,
  })
  revocationQueueSafetyBuffer; // enterprise only

  @attr({
    editType: 'ttl',
    helperTextEnabled:
      'For a certificate to be expunged, the time must be after the expiration time of the certificate (according to the local clock) plus the safety buffer. Defaults to 72 hours.',
    hideToggle: true,
    formatTtl: true,
  })
  safetyBuffer;

  @attr('boolean', { label: 'Tidy the certificate store' })
  tidyCertStore;

  @attr('boolean')
  tidyCertMetadata;

  @attr('boolean', {
    label: 'Tidy cross-cluster revoked certificates',
    subText: 'Remove expired, cross-cluster revocation entries.',
  })
  tidyCrossClusterRevokedCerts; // enterprise only

  @attr('boolean', {
    label: 'Tidy CMPv2 nonce store',
  })
  tidyCmpv2NonceStore; // enterprise only

  @attr('boolean', {
    subText: 'Automatically remove expired issuers after the issuer safety buffer duration has elapsed.',
  })
  tidyExpiredIssuers;

  @attr('boolean', {
    label: 'Tidy legacy CA bundle',
    subText:
      'Backup any legacy CA/issuers bundle (from Vault versions earlier than 1.11) to config/ca_bundle.bak. Migration will only occur after issuer safety buffer has passed.',
  })
  tidyMoveLegacyCaBundle;

  @attr('boolean', {
    label: 'Tidy cross-cluster revocation requests',
  })
  tidyRevocationQueue; // enterprise only

  @attr('boolean', {
    label: 'Tidy revoked certificate issuer associations',
  })
  tidyRevokedCertIssuerAssociations;

  @attr('boolean', {
    label: 'Tidy revoked certificates',
    subText: 'Remove all invalid and expired certificates from storage.',
  })
  tidyRevokedCerts;

  get allGroups() {
    const groups = [{ autoTidy: ['enabled', ...this.autoTidyConfigFields] }, ...this.sharedFields];
    return this._expandGroups(groups);
  }

  // fields that are specific to auto-tidy
  get autoTidyConfigFields() {
    // 'enabled' is not included here because it is responsible for hiding/showing these params and renders separately in the form
    return ['intervalDuration', 'minStartupBackoffDuration', 'maxStartupBackoffDuration'];
  }

  // shared between auto and manual tidy operations
  get sharedFields() {
    const groups = [
      {
        'Universal operations': [
          'tidyCertStore',
          'tidyCertMetadata',
          'tidyRevokedCerts',
          'tidyRevokedCertIssuerAssociations',
          'safetyBuffer',
          'pauseDuration',
        ],
      },
      {
        'ACME operations': ['tidyAcme', 'acmeAccountSafetyBuffer'],
      },
      {
        'Issuer operations': ['tidyExpiredIssuers', 'tidyMoveLegacyCaBundle', 'issuerSafetyBuffer'],
      },
    ];
    if (this.version.isEnterprise) {
      groups.push({
        'Cross-cluster operations': [
          'tidyRevocationQueue',
          'tidyCrossClusterRevokedCerts',
          'tidyCmpv2NonceStore',
          'revocationQueueSafetyBuffer',
        ],
      });
    }
    return groups;
  }

  get formFieldGroups() {
    return this._expandGroups(this.sharedFields);
  }
}
