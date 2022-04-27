import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export const KEY_TYPES = [
  'aes256-gcm96',
  'rsa-2048',
  'rsa-3072',
  'rsa-4096',
  'ecdsa-p256',
  'ecdsa-p384',
  'ecdsa-p521',
];
export default class KeymgmtKeyModel extends Model {
  @attr('string') name;
  @attr('string') backend;

  @attr('string', {
    possibleValues: KEY_TYPES,
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
  @attr() provider; // string, or object with permissions error
  @attr() distribution;

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

  get keyTypeOptions() {
    return expandAttributeMeta(this, ['type'])[0];
  }

  get distFields() {
    return [
      {
        name: 'name',
        type: 'string',
        label: 'Distributed name',
        subText: 'The name given to the key by the provider.',
      },
      { name: 'purpose', type: 'string', label: 'Key Purpose' },
      { name: 'protection', type: 'string', subText: 'Where cryptographic operations are performed.' },
    ];
  }
}
