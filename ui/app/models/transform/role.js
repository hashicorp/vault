import { computed } from '@ember/object';
import DS from 'ember-data';
import { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr } = DS;

const Model = DS.Model.extend({
  // used for getting appropriate options for backend
  idPrefix: 'role/',
  // the id prefixed with `role/` so we can use it as the *secret param for the secret show route
  idForNav: computed('id', 'idPrefix', function() {
    let modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

  name: attr('string', {
    // TODO: make this required for making a transformation
    label: 'Name',
    fieldValue: 'id',
    readOnly: true,
    subText: 'The name for your role. This cannot be edited later.',
  }),
  transformations: attr('string', {
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    label: 'Transformations',
    models: ['transform'],
    subLabel: 'Transformations',
    subText: 'Select which transformations this role will have access to. It must already exist.',
    onlyAllowExisting: true,
  }),

  attrs: computed('transformations', function() {
    let keys = ['name', 'transformations'];
    return expandAttributeMeta(this, keys);
  }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/role/${'id'}`,
});
