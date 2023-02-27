import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['commonName', 'format', 'privateKeyFormat', 'customTtl'],
  },
  {
    'Subject Alternative Name (SAN) Options': [
      'altNames',
      'ipSans',
      'uriSans',
      'otherSans',
      'excludeCnFromSans',
    ],
  },
];
@withFormFields(null, generateFromRole)
export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issue/example?help=1`;
  }
  @attr('string') role;
}
