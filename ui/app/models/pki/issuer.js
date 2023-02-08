import PkiCertificateBaseModel from './certificate/base';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const issuerUrls = ['issuingCertificates', 'crlDistributionPoints', 'ocspServers'];
@withFormFields(
  ['issuerName', 'leafNotAfterBehavior', 'usage', 'manualChain', ...issuerUrls],
  [
    {
      default: [
        'certificate',
        'caChain',
        'commonName',
        'issuerName',
        'serialNumber',
        'keyId',
        'uriSans',
        'ipSans',
        'notValidBefore',
        'notValidAfter',
      ],
    },
    { 'Issuer URLs': issuerUrls },
  ]
)
export default class PkiIssuerModel extends PkiCertificateBaseModel {
  // there are too many differences between what openAPI returns and the designs for the update form
  // manually defining the attrs instead with the correct meta data
  get useOpenAPI() {
    return false;
  }

  get issuerRef() {
    return this.issuerName || this.issuerId;
  }

  @attr isDefault; // readonly
  @attr('string') issuerId;

  @attr('string', {
    label: 'Default key ID',
  })
  keyId;

  @attr('string') issuerName;

  @attr({
    label: 'Leaf notAfter behavior',
    subText:
      'What happens when a leaf certificate is issued, but its NotAfter field (and therefore its expiry date) exceeds that of this issuer.',
    docLink: '/vault/api-docs/secret/pki#update-issuer',
    editType: 'yield',
    valueOptions: ['err', 'truncate', 'permit'],
  })
  leafNotAfterBehavior;

  @attr({
    label: 'Usage',
    subText: 'Allowed usages for this issuer. It can always be read.',
    editType: 'yield',
    valueOptions: [
      { label: 'Issuing certificates', value: 'issuing-certificates' },
      { label: 'Signing CRLs', value: 'crl-signing' },
      { label: 'Signing OCSPs', value: 'ocsp-signing' },
    ],
  })
  usage;

  @attr('string', {
    label: 'Manual chain',
    subText:
      "An advanced field useful when automatic chain building isn't desired. The first element must be the present issuer's reference.",
  })
  manualChain;

  @attr('string', {
    label: 'Issuing certificates',
    subText:
      'The URL values for the Issuing Certificate field. These are different URLs for the same resource, and should be added individually, not in a comma-separated list.',
    editType: 'stringArray',
  })
  issuingCertificates;

  @attr('string', {
    label: 'CRL distribution points',
    subText: 'Specifies the URL values for the CRL Distribution Points field.',
    editType: 'stringArray',
  })
  crlDistributionPoints;

  @attr('string', {
    label: 'OCSP servers',
    subText: 'Specifies the URL values for the OCSP Servers field.',
    editType: 'stringArray',
  })
  ocspServers;

  @lazyCapabilities(apiPath`${'backend'}/issuer/${'issuerId'}`) issuerPath;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/exported`) rotateExported;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/internal`) rotateInternal;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/existing`) rotateExisting;
  @lazyCapabilities(apiPath`${'backend'}/intermediate/cross-sign`) crossSignPath;
  @lazyCapabilities(apiPath`${'backend'}/issuer/${'issuerId'}/sign-intermediate`) signIntermediate;
  get canRotateIssuer() {
    return (
      this.rotateExported.get('canUpdate') !== false ||
      this.rotateExisting.get('canUpdate') !== false ||
      this.rotateInternal.get('canUpdate') !== false
    );
  }
  get canCrossSign() {
    return this.crossSignPath.get('canUpdate') !== false;
  }
  get canSignIntermediate() {
    return this.signIntermediate.get('canUpdate') !== false;
  }
  get canConfigure() {
    return this.issuerPath.get('canUpdate') !== false;
  }
}
