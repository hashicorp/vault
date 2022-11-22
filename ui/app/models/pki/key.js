import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class PkiKeyModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('boolean') isDefault;
  @attr('string') keyRef; // reference to an existing key: either, vault generate identifier, literal string 'default', or the name assigned to the key. Part of the request URL.
  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string') keyName;
  @attr('string') keyType;
  @attr('string', { detailsLabel: 'Key bit length' }) keyBit;

  // TODO refactor when field-to-attrs util is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['keyId', 'keyName', 'keyType', 'keyBit']);
    }
    return this._attributeMeta;
  }
}
