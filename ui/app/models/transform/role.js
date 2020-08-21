import { computed } from '@ember/object';
import DS from 'ember-data';
import { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import attachCapabilities from 'vault/lib/attach-capabilities';

const Model = DS.Model.extend({
  idPrefix: 'role/',

  name: DS.attr('string', {
    // TODO: make this required for making a transformation
    label: 'Name',
    fieldValue: 'id',
    readOnly: true,
    subText: 'The name for your role. This cannot be edited later.',
  }),
  transformations: DS.attr('string', {
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    label: 'Transformations',
    models: ['transform'],
    subLabel: 'Transformations',
    subText: 'Select which transformations this role will have access to. It must already exist.',
  }),

  attrs: computed('transformations', function() {
    let keys = ['name', 'transformations'];
    return expandAttributeMeta(this, keys);
  }),
});

export default attachCapabilities(Model, {
  // TODO: Update to dynamic backend name
  updatePath: apiPath`transform/role/${'id'}`,
});
