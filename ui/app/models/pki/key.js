import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class PkiKeyModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('boolean') isDefault;
  @attr('string', { possibleValues: ['internal', 'external'] }) type;
  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string') keyName;
  @attr('string') keyType;
  @attr('string', { detailsLabel: 'Key bit length' }) keyBits;

  // TODO refactor when field-to-attrs util is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['keyId', 'keyName', 'keyType', 'keyBits']);
    }
    return this._attributeMeta;
  }
}
