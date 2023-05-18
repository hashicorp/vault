import Model, { attr } from '@ember-data/model';
import { service } from '@ember/service';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

@withExpandedAttributes()
export default class PkiTidyModel extends Model {
  // the backend mount is the model id, only one pki/tidy model will ever persist (the auto-tidy config)
  @service version;

  @attr('boolean', {
    label: 'Automatic tidy enabled',
    labelDisabled: 'Automatic tidy disabled',
    defaultValue: false,
  })
  enabled; // auto-tidy only

  @attr({
    label: 'Automatic tidy enabled',
    labelDisabled: 'Automatic tidy disabled',
    mapToBoolean: 'enabled',
    helperTextEnabled:
      'Sets the interval_duration between automatic tidy operations; note that this is from the end of one operation to the start of the next.',
    helperTextDisabled: 'Automatic tidy operations will not run.',
    detailsLabel: 'Automatic tidy duration',
    formatTtl: true,
  })
  intervalDuration; // auto-tidy only

  @attr('string', {
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

  @attr('string', {
    editType: 'ttl',
    helperTextEnabled:
      'For a certificate to be expunged, the time must be after the expiration time of the certificate (according to the local clock) plus the safety buffer. Defaults to 72 hours.',
    hideToggle: true,
    formatTtl: true,
  })
  safetyBuffer;

  @attr('boolean', { label: 'Tidy the certificate store' })
  tidyCertStore;

  @attr('boolean', {
    label: 'Tidy cross-cluster revoked certificates',
    subText: 'Remove expired, cross-cluster revocation entries.',
  })
  tidyCrossClusterRevokedCerts; // enterprise only

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
    subText:
      'Check to associate revoked certificates with their corresponding issuers; this improves the performance of OCSP and CRL building, by shifting work to a tidy operation instead. It is suggested to run this tidy when removing or importing new issuers and on the first upgrade to a post-1.11 Vault version, but otherwise not to run it during automatic tidy operations.',
  })
  tidyRevokedCertIssuerAssociations;

  @attr('boolean', {
    label: 'Tidy revoked certificates',
    subText: 'Remove all invalid and expired certificates from storage.',
  })
  tidyRevokedCerts;

  /*
  NOT IN DOCS - check with crypto
  @attr('string') acme_account_safety_buffer;
  @attr('boolean', { defaultValue: false }) maintain_stored_certificate_counts;
  @attr('boolean', { defaultValue: false }) publish_stored_certificate_count_metrics;
  @attr('boolean', { defaultValue: false }) tidy_acme;
  */

  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/config/auto-tidy?help=1`;
  }

  get displayGroups() {
    const groups = [
      { default: ['enabled', 'intervalDuration'] },
      {
        'Universal operations': ['tidyCertStore', 'tidyRevokedCerts', 'safetyBuffer', 'pauseDuration'],
      },
      {
        'Issuer operations': [
          'tidyExpiredIssuers',
          'tidyMoveLegacyCaBundle',
          'tidyRevokedCertIssuerAssociations',
          'issuerSafetyBuffer',
        ],
      },
    ];
    if (this.version.isEnterprise) {
      groups.push({
        'Cross-cluster operations': [
          'tidyRevocationQueue',
          'tidyCrossClusterRevokedCerts',
          'revocationQueueSafetyBuffer',
        ],
      });
    }
    return this._expandGroups(groups);
  }

  get sharedFields() {
    const groups = [
      {
        'Universal operations': ['tidyCertStore', 'tidyRevokedCerts', 'safetyBuffer', 'pauseDuration'],
      },
      {
        'Issuer operations': [
          'tidyExpiredIssuers',
          'tidyMoveLegacyCaBundle',
          'tidyRevokedCertIssuerAssociations',
          'issuerSafetyBuffer',
        ],
      },
    ];
    if (this.version.isEnterprise) {
      groups.push({
        'Cross-cluster operations': [
          'tidyRevocationQueue',
          'tidyCrossClusterRevokedCerts',
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
