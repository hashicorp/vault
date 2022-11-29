import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields(['keyId', 'keyName', 'keyType', 'keyBits'])
export default class PkiKeyModel extends Model {
  @service secretMountPath;

  @attr('boolean') isDefault;
  @attr('string', { possibleValues: ['internal', 'external'] }) type;
  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string') keyName;
  @attr('string') keyType;
  @attr('string') keyBits;

  get backend() {
    return this.secretMountPath.currentPath;
  }
}
