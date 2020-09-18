import DS from 'ember-data';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

const Model = DS.Model.extend({
  idPrefix: 'alphabet/',
  name: attr('string', {
    fieldValue: 'id',
    readOnly: true,
    subText: 'The alphabet name. Keep in mind that spaces are not allowed and this cannot be edited later.',
  }),
  alphabet: attr('string', {
    label: 'Alphabet',
    subText:
      'Provide the set of valid UTF-8 characters contained within both the input and transformed value. Read more.',
  }),

  attrs: computed(function() {
    let keys = ['name', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }),

  backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/alphabet/${'id'}`,
});
