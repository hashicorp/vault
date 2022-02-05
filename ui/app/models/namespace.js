import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  path: attr('string', {
    validationAttr: 'pathIsValid',
    invalidMessage: 'You have entered and invalid path please only include letters, numbers, -, ., and _.',
  }),
  pathIsValid: computed('path', function () {
    return this.path && this.path.match(/^[\w\d-.]+$/g);
  }),
  description: attr('string', {
    editType: 'textarea',
  }),
  fields: computed(function () {
    return expandAttributeMeta(this, ['path']);
  }),
});
