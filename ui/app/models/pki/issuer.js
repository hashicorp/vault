import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

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
export default class PkiIssuerModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string', {
    label: 'Issuer name',
    fieldValue: 'id',
  })
  name;

  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/issuer/example?help=1`;
  }

  @attr('boolean') isDefault;
  @attr('string') issuerName;

  // Form Fields not hidden in toggle options
  _attributeMeta = null;
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, [
        'name',
        'leafNotAfterBehavior',
        'usage',
        'manualChain',
        'issuingCertifications',
        'crlDistributionPoints',
        'ocspServers',
        'deltaCrlUrls', // new endpoint, mentioned in RFC, but need to confirm it's there.
      ]);
    }
    return this._attributeMeta;
  }
}
