import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

let { attr } = DS;

export default DS.Model.extend({
  name: attr('string'),
  policy: attr('string'),
  policyType: computed(function() {
    return this.constructor.modelName.split('/')[1];
  }),

  updatePath: lazyCapabilities(apiPath`sys/policies/${'policyType'}/${'id'}`, 'id', 'policyType'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),
  format: computed('policy', function() {
    let policy = this.get('policy');
    let isJSON;
    try {
      let parsed = JSON.parse(policy);
      if (parsed) {
        isJSON = true;
      }
    } catch (e) {
      // can't parse JSON
      isJSON = false;
    }
    return isJSON ? 'json' : 'hcl';
  }),
});
