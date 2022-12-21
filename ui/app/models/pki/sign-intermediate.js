import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './certificate/base';

@withFormFields(['csr'])
export default class PkiSignIntermediateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issuer/example/sign-intermediate?help=1`;
  }
}
