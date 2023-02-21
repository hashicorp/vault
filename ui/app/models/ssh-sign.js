import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const CREATE_FIELDS = [
  'publicKey',
  'keyId',
  'validPrincipals',
  'certType',
  'criticalOptions',
  'extensions',
  'ttl',
];

const DISPLAY_FIELDS = ['signedKey', 'leaseId', 'renewable', 'leaseDuration', 'serialNumber'];

export default Model.extend({
  role: attr('object', {
    readOnly: true,
  }),
  publicKey: attr('string', {
    label: 'Public Key',
    editType: 'textarea',
  }),
  ttl: attr({
    label: 'TTL',
    editType: 'ttl',
  }),
  validPrincipals: attr('string'),
  certType: attr('string', {
    defaultValue: 'user',
    label: 'Certificate Type',
    possibleValues: ['user', 'host'],
  }),
  keyId: attr('string', {
    label: 'Key ID',
  }),
  criticalOptions: attr('object'),
  extensions: attr('object'),

  leaseId: attr('string', {
    label: 'Lease ID',
  }),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  serialNumber: attr('string'),
  signedKey: attr('string'),

  attrs: computed('signedKey', function () {
    const keys = this.signedKey ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),
});
