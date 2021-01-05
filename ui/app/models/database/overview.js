import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
// import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const M = Model.extend({
  idPrefix: 'overview/',
  idForNav: computed('id', 'idPrefix', function() {
    let modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

  name: attr('string', {
    fieldValue: 'id',
    readOnly: true,
    subText:
      'Templates allow Vault to determine what and how to capture the value to be transformed. This cannot be edited later.',
  }),
  //   type: attr('string', { defaultValue: 'regex' }),
  //   pattern: attr('string', {
  //     subText: 'The template’s pattern defines the data format. Expressed in regex.',
  //   }),
  //   alphabet: attr('array', {
  //     subText:
  //       'Alphabet defines a set of characters (UTF-8) that is used for FPE to determine the validity of plaintext and ciphertext values. You can choose a built-in one, or create your own.',
  //     editType: 'searchSelect',
  //     fallbackComponent: 'string-list',
  //     label: 'Alphabet',
  //     models: ['transform/alphabet'],
  //     selectLimit: 1,
  //   }),

  //   attrs: computed(function() {
  //     let keys = ['name', 'pattern', 'alphabet'];
  //     return expandAttributeMeta(this, keys);
  //   }),

  //   backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(M, {
  // ARG TODO figure out the apiPath for the Overview
  updatePath: apiPath`${'backend'}/template/${'id'}`,
});
