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
    defaultShown: 'All versions enabled',
  })
  minEnabledVersion;

  @attr('array')
  versions;

  // The following are calculated in serializer
  @attr('date')
  created;

  @attr('date', {
    defaultShown: 'Not yet rotated',
  })
  lastRotated;

  // The following are from endpoints other than the main read one
  @attr('string') provider;

  icon = 'key';

  get hasVersions() {
    return this.versions.length > 1;
  }

  get createFields() {
    const createFields = ['name', 'type', 'deletionAllowed'];
    return expandAttributeMeta(this, createFields);
  }

  get updateFields() {
    return expandAttributeMeta(this, ['minEnabledVersion', 'deletionAllowed']);
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
