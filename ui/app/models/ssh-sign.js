import DS from 'ember-data';
import Ember from 'ember';
const { attr } = DS;
const { computed, get } = Ember;
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
  publicKey: attr('string'),
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
    get(this.constructor, 'attributes').forEach((meta, name) => {
      const index = keys.indexOf(name);
      if (index === -1) {
        return;
      }
      keys.replace(index, 1, {
        type: meta.type,
        name,
        options: meta.options,
      });
    });
    return keys;
  }),
});
