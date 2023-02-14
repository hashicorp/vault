import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields()
export default class PkiUrlsModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backendPath) {
    return `/v1/${backendPath}/config/urls?help=1`;
  }

  @attr({
    label: 'Issuing certificates',
    subText:
      'The URL values for the Issuing Certificate field. These are different URLs for the same resource, and should be added individually, not in a comma-separated list.',
    showHelpText: false,
    editType: 'stringArray',
  })
  issuingCertificates;

  @attr({
    label: 'CRL distribution points',
    subText: 'Specifies the URL values for the CRL Distribution Points field.',
    showHelpText: false,
    editType: 'stringArray',
  })
  crlDistributionPoints;

  @attr({
    label: 'OSCP Servers',
    subText: 'Specifies the URL values for the OCSP Servers field.',
    showHelpText: false,
    editType: 'stringArray',
  })
  ocspServers;

  @lazyCapabilities(apiPath`${'id'}/config/urls`, 'id') urlsPath;

  get canSet() {
    return this.urlsPath.get('canCreate') !== false;
  }
}
