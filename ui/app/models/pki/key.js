import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class PkiKeyModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('boolean') isDefault;
  @attr('string', { possibleValues: ['internal', 'external'] }) type;
  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string') keyName;
  @attr('string') keyType;
  @attr('string', { detailsLabel: 'Key bit length' }) keyBits; // TODO confirm with crypto team to remove this field from details page

  // TODO refactor when field-to-attrs util is refactored as decorator
  constructor() {
    super(...arguments);
    this.formFields = expandAttributeMeta(this, ['keyId', 'keyName', 'keyType', 'keyBits']);
  }
}
