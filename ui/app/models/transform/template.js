import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const M = Model.extend({
  idPrefix: 'template/',
  idForNav: computed('id', 'idPrefix', function () {
    const modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

  name: attr('string', {
    fieldValue: 'id',
    readOnly: true,
    subText:
      'Templates allow Vault to determine what and how to capture the value to be transformed. This cannot be edited later.',
  }),
  type: attr('string', { defaultValue: 'regex' }),
  pattern: attr('string', {
    editType: 'regex',
    subText: 'The template’s pattern defines the data format. Expressed in regex.',
  }),
  alphabet: attr('array', {
    subText:
      'Alphabet defines a set of characters (UTF-8) that is used for FPE to determine the validity of plaintext and ciphertext values. You can choose a built-in one, or create your own.',
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    label: 'Alphabet',
    models: ['transform/alphabet'],
    selectLimit: 1,
  }),
  encodeFormat: attr('string'),
  decodeFormats: attr(),
  backend: attr('string', { readOnly: true }),

  readAttrs: computed(function () {
    const keys = ['name', 'pattern', 'encodeFormat', 'decodeFormats', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }),
  writeAttrs: computed(function () {
    return expandAttributeMeta(this, ['name', 'pattern', 'alphabet']);
  }),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/template/${'id'}`,
});
