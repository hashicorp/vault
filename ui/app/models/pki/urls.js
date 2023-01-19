import Model, { attr } from '@ember-data/model';
import { service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields()
export default class PkiUrlsModel extends Model {
  @service secretMountPath;
  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/config/urls?help=1`;
  }
  get backend() {
    return this.id;
  }

  @attr({
    label: 'Issuing certificates',
    subText:
      'The URL values for the Issuing Certificate field. These are different URLs for the same resource, and should be added individually, not in a comma-separated list.',
    showHelpText: false,
  })
  issuingCertificates;

  @attr({
    label: 'CRL distribution points',
    subText: 'Specifies the URL values for the CRL Distribution Points field.',
  })
  crlDistributionPoints;

  @lazyCapabilities(apiPath`${'backend'}/config/urls`, 'backend') urlsPath;

  get canSet() {
    return this.urlsPath.get('canCreate') !== false;
  }
}
