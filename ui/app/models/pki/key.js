import Model, { attr } from '@ember-data/model';

export default class PkiKeyModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('boolean') isDefault;
  @attr('string') keyRef; // reference to an existing key: either, vault generate identifier, literal string 'default', or the name assigned to the key. Part of the request URL.
  @attr('string') keyId;
  @attr('string') keyName;
  @attr('string') keyType;
}
