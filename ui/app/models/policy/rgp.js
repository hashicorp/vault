import DS from 'ember-data';
import Ember from 'ember';

import PolicyModel from '../policy';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

let { attr } = DS;
let { computed } = Ember;

export default PolicyModel.extend({
  enforcementLevel: attr('string', {
    possibleValues: ['advisory', 'soft-mandatory', 'hard-mandatory'],
    defaultValue: 'hard-mandatory',
  }),

  additionalAttrs: computed(function() {
    return expandAttributeMeta(this, ['enforcementLevel']);
  }),
});
