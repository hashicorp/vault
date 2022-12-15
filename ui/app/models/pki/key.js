import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  type: [{ type: 'presence', message: 'Type is required.' }],
  keyType: [{ type: 'presence', message: 'Please select a key type.' }],
};
const displayFields = ['keyId', 'keyName', 'keyType', 'keyBits'];
const formFieldGroups = [{ default: ['keyName', 'type'] }, { 'Key parameters': ['keyType', 'keyBits'] }];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class PkiKeyModel extends Model {
  @service secretMountPath;

  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string', { subText: 'Optional, human-readable name for this key.' }) keyName;
  @attr('string') privateKey;
  @attr('string', {
    noDefault: true,
    possibleValues: ['internal', 'exported'],
    subText:
      'The type of operation. If exported, the private key will be returned in the response; if internal the private key will not be returned and cannot be retrieved later.',
  })
  type;
  @attr('string', {
    noDefault: true,
    possibleValues: ['rsa', 'ec', 'ed25519'],
    subText: 'The type of key that will be generated. Must be rsa, ed25519, or ec. ',
  })
  keyType;
  @attr('string', {
    label: 'Key bits',
    noDefault: true,
    subText: 'Bit length of the key to generate.',
  })
  keyBits; // no possibleValues because dependent on selected key type

  get backend() {
    return this.secretMountPath.currentPath;
  }
}
