import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

const { attr } = DS;

export default DS.Model.extend({
  idPrefix: 'cert/',

  backend: attr('string', {
    readOnly: true,
  }),
  useOpenAPI: true,
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
  otherSans: attr({
    helpText:
      'The format is the same as OpenSSL: <oid>;<type>:<value> where the only current valid type is UTF8',
  }),

  fieldsToAttrs(fieldGroups) {
    return fieldToAttrs(this, fieldGroups);
  },

  fieldDefinition: computed('newFields', function() {
    let groups = [
      { default: ['commonName', 'format'] },
      { Options: ['altNames', 'ipSans', 'ttl', 'excludeCnFromSans', 'otherSans'] },
    ];
    if (this.newFields.length) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
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
