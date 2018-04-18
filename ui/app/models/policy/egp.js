import DS from 'ember-data';
import Ember from 'ember';

import PolicyModel from './rgp';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

let { attr } = DS;
let { computed } = Ember;

export default PolicyModel.extend({
  paths: attr({
    editType: 'stringArray',
  }),
  additionalAttrs: computed(function() {
    return expandAttributeMeta(this, ['enforcementLevel', 'paths']);
  }),
});
