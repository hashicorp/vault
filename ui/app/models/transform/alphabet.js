import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const M = Model.extend({
  idPrefix: 'alphabet/',
  idForNav: computed('id', 'idPrefix', function () {
    const modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

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

  attrs: computed(function () {
    const keys = ['name', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }),

  backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/alphabet/${'id'}`,
});
