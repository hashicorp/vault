import { assert } from '@ember/debug';
import { computed } from '@ember/object';
import DS from 'ember-data';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default DS.Model.extend({
  formFields: computed(function() {
    return assert('formFields should be overridden', false);
  }),

  fields: computed('formFields', 'formFields.[]', function() {
    return expandAttributeMeta(this, this.get('formFields'));
  }),

  identityType: computed(function() {
    let modelType = this.constructor.modelName.split('/')[1];
    return modelType;
  }),
});
