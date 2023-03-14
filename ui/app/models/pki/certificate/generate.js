import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['commonName', 'customTtl', 'format', 'privateKeyFormat'],
  },
  {
    'Subject Alternative Name (SAN) Options': [
      'excludeCnFromSans',
      'altNames',
      'ipSans',
      'uriSans',
      'otherSans',
    ],
  },
];
@withFormFields(null, generateFromRole)
export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issue/example?help=1`;
  }
  @attr('string') role; // role name to issue certificate against for request URL
}
