import Ember from 'ember';
import DS from 'ember-data';

const { attr } = DS;
const { computed } = Ember;
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default DS.Model.extend({
  path: attr('string'),
  description: attr('string', {
    editType: 'textarea',
  }),
  fields: computed(function() {
    return expandAttributeMeta(this, ['path']);
  }),
});
