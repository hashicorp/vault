import Ember from 'ember';
import DS from 'ember-data';

const { attr } = DS;
const { computed } = Ember;
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default DS.Model.extend({
  path: attr('string', {
    validationAttr: 'pathIsValid',
    invalidMessage: 'You have entered and invalid path please only include letters, numbers, -, ., and _.',
  }),
  pathIsValid: computed.match('path', /^[0-9A-Z-._]+$/gi),
  description: attr('string', {
    editType: 'textarea',
  }),
  fields: computed(function() {
    return expandAttributeMeta(this, ['path']);
  }),
});
