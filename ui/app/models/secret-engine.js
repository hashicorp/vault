import Ember from 'ember';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { fragment } from 'ember-data-model-fragments/attributes';

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

//identity will be managed separately and the inclusion
//of the system backend is an implementation detail
const LIST_EXCLUDED_BACKENDS = ['system', 'identity'];

export default DS.Model.extend({
  path: attr('string'),
  accessor: attr('string'),
  name: attr('string'),
  type: attr('string'),
  description: attr('string'),
  config: fragment('mount-config', { defaultValue: {} }),
  options: fragment('mount-options', { defaultValue: {} }),
  local: attr('boolean'),
  sealWrap: attr('boolean'),

  modelTypeForKV: computed('type', 'options.version', function() {
    let type = this.get('type');
    let version = this.get('options.version');
    let modelType = 'secret';
    if ((type === 'kv' || type === 'generic') && version === 2) {
      modelType = 'secret-v2';
    }
    return modelType;
  }),

  formFields: [
    'type',
    'path',
    'description',
    'accessor',
    'local',
    'sealWrap',
    'config.{defaultLeaseTtl,maxLeaseTtl}',
    'options.{version}',
  ],

  attrs: computed('formFields', function() {
    return expandAttributeMeta(this, this.get('formFields'));
  }),

  shouldIncludeInList: computed('type', function() {
    return !LIST_EXCLUDED_BACKENDS.includes(this.get('type'));
  }),

  localDisplay: Ember.computed('local', function() {
    return this.get('local') ? 'local' : 'replicated';
  }),

  // ssh specific ones
  privateKey: attr('string'),
  publicKey: attr('string'),
  generateSigningKey: attr('boolean', {
    defaultValue: true,
  }),

  saveCA(options) {
    if (this.get('type') !== 'ssh') {
      return;
    }
    if (options.isDelete) {
      this.setProperties({
        privateKey: null,
        publicKey: null,
        generateSigningKey: false,
      });
    }
    return this.save({
      adapterOptions: {
        options: options,
        apiPath: 'config/ca',
        attrsToSend: ['privateKey', 'publicKey', 'generateSigningKey'],
      },
    });
  },

  saveZeroAddressConfig() {
    return this.save({
      adapterOptions: {
        adapterMethod: 'saveZeroAddressConfig',
      },
    });
  },

  zeroAddressPath: lazyCapabilities(apiPath`${'id'}/config/zeroaddress`, 'id'),
  canEditZeroAddress: computed.alias('zeroAddressPath.canUpdate'),

  // aws backend attrs
  lease: attr('string'),
  leaseMax: attr('string'),
});
