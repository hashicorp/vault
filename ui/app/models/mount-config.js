import attr from 'ember-data/attr';
import Fragment from 'ember-data-model-fragments/fragment';

export default Fragment.extend({
  defaultLeaseTtl: attr({
    label: 'Default Lease TTL',
    editType: 'ttl',
  }),
  maxLeaseTtl: attr({
    label: 'Max Lease TTL',
    editType: 'ttl',
  }),
});
