import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default DS.Model.extend({
  idPrefix: 'cert/',

  backend: attr('string', {
    readOnly: true,
  }),
  //the id prefixed with `cert/` so we can use it as the *secret param for the secret show route
  idForNav: attr('string', {
    readOnly: true,
  }),
  DISPLAY_FIELDS: computed(function() {
    return [
      'certificate',
      'issuingCa',
      'caChain',
      'privateKey',
      'privateKeyType',
      'serialNumber',
      'revocationTime',
    ];
  }),
  role: attr('object', {
    readOnly: true,
  }),

  revocationTime: attr('number'),
  commonName: attr('string', {
    label: 'Common Name',
  }),

  altNames: attr('string', {
    label: 'DNS/Email Subject Alternative Names (SANs)',
  }),

  ipSans: attr('string', {
    label: 'IP Subject Alternative Names (SANs)',
  }),
  otherSans: attr({
    editType: 'stringArray',
    label: 'Other SANs',
    helpText:
      'The format is the same as OpenSSL: <oid>;<type>:<value> where the only current valid type is UTF8',
  }),

  ttl: attr({
    label: 'TTL',
    editType: 'ttl',
  }),

  format: attr('string', {
    defaultValue: 'pem',
    possibleValues: ['pem', 'der', 'pem_bundle'],
  }),

  excludeCnFromSans: attr('boolean', {
    label: 'Exclude Common Name from Subject Alternative Names (SANs)',
    defaultValue: false,
  }),

  certificate: attr('string'),
  issuingCa: attr('string', {
    label: 'Issuing CA',
  }),
  caChain: attr('string', {
    label: 'CA chain',
  }),
  privateKey: attr('string'),
  privateKeyType: attr('string'),
  serialNumber: attr('string'),

  fieldsToAttrs(fieldGroups) {
    return fieldToAttrs(this, fieldGroups);
  },

  fieldDefinition: computed(function() {
    const groups = [
      { default: ['commonName', 'format'] },
      { Options: ['altNames', 'ipSans', 'ttl', 'excludeCnFromSans', 'otherSans'] },
    ];
    return groups;
  }),

  fieldGroups: computed('fieldDefinition', function() {
    return this.fieldsToAttrs(this.get('fieldDefinition'));
  }),

  attrs: computed('certificate', 'csr', function() {
    let keys = this.get('certificate') || this.get('csr') ? this.get('DISPLAY_FIELDS').slice(0) : [];
    return expandAttributeMeta(this, keys);
  }),

  toCreds: computed(
    'certificate',
    'issuingCa',
    'caChain',
    'privateKey',
    'privateKeyType',
    'revocationTime',
    'serialNumber',
    function() {
      const props = this.getProperties(
        'certificate',
        'issuingCa',
        'caChain',
        'privateKey',
        'privateKeyType',
        'revocationTime',
        'serialNumber'
      );
      const propsWithVals = Object.keys(props).reduce((ret, prop) => {
        if (props[prop]) {
          ret[prop] = props[prop];
          return ret;
        }
        return ret;
      }, {});
      return JSON.stringify(propsWithVals, null, 2);
    }
  ),

  revokePath: lazyCapabilities(apiPath`${'backend'}/revoke`, 'backend'),
  canRevoke: alias('revokePath.canUpdate'),
});
