import { computed } from '@ember/object';
import DS from 'ember-data';

const { attr } = DS;
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default DS.Model.extend({
  name: attr('string'),
  attrs: computed(function() {
    return expandAttributeMeta(this, ['name']);
  }),
});
