import { attr } from '@ember-data/model';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import PkiCertificateBaseModel from './certificate/base';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    {
      type: 'containsWhiteSpace',
      message: 'Name cannot contain whitespace.',
    },
  ],
};

@withModelValidations(validations)
@withFormFields([
  'certificate',
  'caChain',
  'commonName',
  'issuerName',
  'notValidBefore',
  'serialNumber',
  'keyId',
  'notValidAfter',
  'notValidBefore',
])
export default class PkiIssuerModel extends PkiCertificateBaseModel {
  getHelpUrl(backend) {
    return `/v1/${backend}/issuer/example?help=1`;
  }
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

  @lazyCapabilities(apiPath`${'backend'}/issuer`) issuerPath;
  get canRotateIssuer() {
    return true;
  }

  get canCrossSign() {
    return true;
  }

  get canSignIntermediate() {
    return true;
  }

  get canConfigure() {
    return true;
  }
}
