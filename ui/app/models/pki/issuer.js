import PkiCertificateBaseModel from './certificate/base';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields(null, [
  {
    default: [
      'certificate',
      'caChain',
      'commonName',
      'issuerName',
      'notValidBefore',
      'serialNumber',
      'keyId',
      'uriSans',
      'notValidAfter',
    ],
  },
  { 'Issuer URLs': ['issuingCertificates', 'crlDistributionPoints', 'ocspServers', 'deltaCrlUrls'] },
])
export default class PkiIssuerModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issuer/example?help=1`;
  }

  @attr('string') issuerId;
  @attr('string', { displayType: 'masked' }) certificate;
  @attr('string', { displayType: 'masked', label: 'CA Chain' }) caChain;
  @attr('date', {
    label: 'Issue date',
  })
  notValidBefore;

  @attr('string', {
    label: 'Default key ID',
  })
  keyId;

  @attr({
    label: 'Subject Alternative Names',
  })
  uriSans;

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
