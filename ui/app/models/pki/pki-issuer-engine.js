import Model, { attr } from '@ember-data/model';
// import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
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
export default class PkiIssuersEngineModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string', {
    label: 'Issuer name',
    fieldValue: 'id',
  })
  name;

  // get useOpenAPI() {
  //   return true;
  // }
  // getHelpUrl(backend) {
  //   return `/v1/${backend}/issuers/example?help=1`;
  // }

  @attr('object') keyInfo;

  // Form Fields not hidden in toggle options
  // _attributeMeta = null;
  // get formFields() {
  //   if (!this._attributeMeta) {
  //     this._attributeMeta = expandAttributeMeta(this, [
  //       'name',
  //       'leafNotAfterBehavior',
  //       'usage',
  //       'manualChain',
  //       'issuingCertifications',
  //       'crlDistributionPoints',
  //       'ocspServers',
  //     ]);
  //   }
  //   return this._attributeMeta;
  // }
}
