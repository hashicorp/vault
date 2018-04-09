import Ember from 'ember';
import DS from 'ember-data';
import { queryRecord } from 'ember-computed-query';
import { fragment } from 'ember-data-model-fragments/attributes';

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
  config: attr('object'),
  options: fragment('mount-options'),
  local: attr('boolean'),
  sealWrap: attr('boolean'),

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

  zeroAddressPath: queryRecord(
    'capabilities',
    context => {
      const { id } = context.getProperties('backend', 'id');
      return {
        id: `${id}/config/zeroaddress`,
      };
    },
    'id'
  ),
  canEditZeroAddress: computed.alias('zeroAddressPath.canUpdate'),

  // aws backend attrs
  lease: attr('string'),
  leaseMax: attr('string'),
});
