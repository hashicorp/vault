import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  name: attr('string'),
  policy: attr('string'),
  policyType: computed('constructor.modelName', function () {
    return this.constructor.modelName.split('/')[1];
  }),
  updatePath: lazyCapabilities(apiPath`sys/policies/${'policyType'}/${'id'}`, 'id', 'policyType'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),
  format: computed('policy', function () {
    const policy = this.policy;
    let isJSON;
    try {
      const parsed = JSON.parse(policy);
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
