import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  type: [{ type: 'presence', message: 'Type is required.' }],
};

const displayFields = ['keyId', 'keyName', 'keyType', 'keyBits'];
const formFieldGroups = [{ default: ['keyName', 'type'] }, { 'Key parameters': ['keyType', 'keyBits'] }];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class PkiKeyModel extends Model {
  @service secretMountPath;

  @attr('boolean') isDefault;
  @attr('string', {
    noDefault: true,
    possibleValues: ['internal', 'external'],
    subText:
      'The type of operation. If exported, the private key will be returned in the response; if internal the private key will not be returned and cannot be retrieved later.',
  })
  type;
  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string', { subText: 'Optional, human-readable name for this key.' }) keyName;
  @attr('string', {
    defaultValue: 'rsa',
    possibleValues: ['rsa', 'ec', 'ed25519'],
    subText: 'The type of key that will be generated. Must be rsa, ed25519, or ec. ',
  })
  keyType;
  @attr('string', {
    label: 'Key bits',
    defaultValue: '2048',
    subText: 'Bit length of the key to generate.',
  })
  keyBits; // no possibleValues set here because dependent on selected key type

  get backend() {
    return this.secretMountPath.currentPath;
  }
}
