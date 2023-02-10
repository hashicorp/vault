import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';
import PkiCertificateBaseModel from './base';

const generateFromRole = [
  {
    default: ['commonName'],
  },
  {
    'Subject Alternative Name (SAN) Options': [
      'altNames',
      'ipSans',
      'uriSans',
      'otherSans',
      'ttl',
      'format',
      'privateKeyFormat',
      'excludeCnFromSans',
      'notAfter',
    ],
  },
];
const validations = {
  commonName: [{ type: 'presence', message: 'Common name is required.' }],
};

@withModelValidations(validations)
@withFormFields(null, generateFromRole)
export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issue/example?help=1`;
  }
  @attr('string') role;
}
