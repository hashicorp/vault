import { attr } from '@ember-data/model';
import { computed } from '@ember/object';

import PolicyModel from '../policy';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default PolicyModel.extend({
  enforcementLevel: attr('string', {
    possibleValues: ['advisory', 'soft-mandatory', 'hard-mandatory'],
    defaultValue: 'hard-mandatory',
  }),

  additionalAttrs: computed(function () {
    return expandAttributeMeta(this, ['enforcementLevel']);
  }),
});
