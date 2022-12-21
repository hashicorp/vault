import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './certificate/base';
@withFormFields(null, [
  {
    default: [
      'csr',
      'useCsrValues',
      'commonName',
      'customTtl',
      'notBeforeDuration',
      'format',
      'permittedDnsDomains',
      'maxPathLength',
    ],
  },
  {
    'Key parameters': ['keyId', 'skid'],
  },
  {
    'Subject Alternative Name (SAN) Options': ['altNames', 'ipSans', 'uriSans', 'otherSans'],
  },
  {
    'Other subject data': [],
  },
])
export default class PkiSignIntermediateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issuer/example/sign-intermediate?help=1`;
  }

  @attr({
    label: 'Not valid after',
    detailsLabel: 'Issued certificates expire after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
    editType: 'yield',
  })
  customTtl;

  @attr({
    label: 'Backdate validity',
    detailsLabel: 'Issued certificate backdating',
    helperTextDisabled: 'Vault will use the default value, 30s',
    helperTextEnabled:
      'Also called the not_before_duration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.',
    editType: 'ttl',
    defaultValue: '30s',
  })
  notBeforeDuration;

  @attr({
    hideFormWhen: ['useCsrValues', true],
  })
  commonName;
}
