import DS from 'ember-data';
import Ember from 'ember';
import { queryRecord } from 'ember-computed-query';

let { attr } = DS;
let { computed } = Ember;

export default DS.Model.extend({
  name: attr('string'),
  policy: attr('string'),
  policyType: computed(function() {
    return this.constructor.modelName.split('/')[1];
  }),

  updatePath: queryRecord(
    'capabilities',
    context => {
      const { policyType, id } = context.getProperties('policyType', 'id');
      if (!policyType && id) {
        return;
      }
      return {
        id: `sys/${policyType}/policies/${id}`,
      };
    },
    'id',
    'policyType'
  ),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  canRead: computed.alias('updatePath.canRead'),
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
