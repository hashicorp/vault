import { computed } from '@ember/object';
import DS from 'ember-data';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const { attr } = DS;
const CREATE_FIELDS = [
  'publicKey',
  'keyId',
  'validPrincipals',
  'certType',
  'criticalOptions',
  'extension',
  'ttl',
];

const DISPLAY_FIELDS = ['signedKey', 'leaseId', 'renewable', 'leaseDuration', 'serialNumber'];

export default DS.Model.extend({
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
  extension: attr('object'),

  leaseId: attr('string', {
    label: 'Lease ID',
  }),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  serialNumber: attr('string'),
  signedKey: attr('string'),

  attrs: computed('signedKey', function() {
    let keys = this.get('signedKey') ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),
});
