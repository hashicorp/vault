import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class KeymgmtKeyModel extends Model {
  @attr('string') name;
  @attr('string') backend;

  @attr('string', {
    possibleValues: ['aes256-gcm96', 'rsa-2048', 'rsa-3072', 'rsa-4096'],
  })
  type;

  @attr('boolean', {
    defaultValue: false,
  })
  deletionAllowed;

  @attr('number', {
    label: 'Current version',
  })
  latestVersion;

  @attr('number', {
    defaultValue: 0,
  })
  minEnabledVersion;

  @attr('array')
  keys;

  @attr('date', {
    defaultShown: 'Not yet rotated',
  })
  lastRotated;

  // The following are from endpoints other than the main read one
  @attr('string') provider;

  get hasKeys() {
    return this.keys.length > 0;
  }

  get fields() {
    // Create and update have different fields. On create, must create and then update
    // const updateFields = ['minEnabledVersion', 'deletionAllowed'];
    const createFields = ['name', 'type', 'deletionAllowed'];
    // const showFields = ['name', 'type', 'deletionAllowed', 'latestVersion', 'minEnabledVersion', 'keys'];
    return expandAttributeMeta(this, createFields);
  }

  get showFields() {
    return expandAttributeMeta(this, [
      'name',
      'created',
      'type',
      'deletionAllowed',
      'latestVersion',
      'minEnabledVersion',
      'lastRotated',
    ]);
  }
}
