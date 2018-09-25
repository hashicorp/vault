import { computed } from '@ember/object';
import DS from 'ember-data';

import PolicyModel from './rgp';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

let { attr } = DS;

export default PolicyModel.extend({
  paths: attr({
    editType: 'stringArray',
  }),
  additionalAttrs: computed(function() {
    return expandAttributeMeta(this, ['enforcementLevel', 'paths']);
  }),
});
