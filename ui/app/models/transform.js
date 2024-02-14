/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

// these arrays define the order in which the fields will be displayed
// see
// https://developer.hashicorp.com/vault/api-docs/secret/transform#create-update-transformation-deprecated-1-6
const TYPES = [
  {
    value: 'fpe',
    displayName: 'Format Preserving Encryption (FPE)',
  },
  {
    value: 'masking',
    displayName: 'Masking',
  },
  {
    value: 'tokenization',
    displayName: 'Tokenization',
  },
];

const TWEAK_SOURCE = [
  {
    value: 'supplied',
    displayName: 'supplied',
  },
  {
    value: 'generated',
    displayName: 'generated',
  },
  {
    value: 'internal',
    displayName: 'internal',
  },
];

export default Model.extend({
  name: attr('string', {
    // CBS TODO: make this required for making a transformation
    label: 'Name',
    readOnly: true,
    subText: 'The name for your transformation. This cannot be edited later.',
  }),
  type: attr('string', {
    defaultValue: 'fpe',
    label: 'Type',
    possibleValues: TYPES,
    subText:
      'Vault provides two types of transformations: Format Preserving Encryption (FPE) is reversible, while Masking is not. This cannot be edited later.',
  }),
  tweak_source: attr('string', {
    defaultValue: 'supplied',
    label: 'Tweak source',
    possibleValues: TWEAK_SOURCE,
    subText: `A tweak value is used when performing FPE transformations. This can be supplied, generated, or internal.`, // CBS TODO: I do not include the link here.  Need to figure out the best way to approach this.
  }),
  masking_character: attr('string', {
    characterLimit: 1,
    defaultValue: '*',
    label: 'Masking character',
    subText: 'Specify which character you’d like to mask your data.',
  }),
  template: attr('array', {
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    label: 'Template', // CBS TODO: make this required for making a transformation
    models: ['transform/template'],
    selectLimit: 1,
    onlyAllowExisting: true,
    subText:
      'Templates allow Vault to determine what and how to capture the value to be transformed. Type to use an existing template or create a new one.',
  }),
  allowed_roles: attr('array', {
    editType: 'searchSelect',
    isSectionHeader: true,
    label: 'Allowed roles',
    fallbackComponent: 'string-list',
    models: ['transform/role'],
    subText: 'Search for an existing role, type a new role to create it, or use a wildcard (*).',
    wildcardLabel: 'role',
  }),
  deletion_allowed: attr('boolean', {
    label: 'Allow deletion',
    subText:
      'If checked, this transform can be deleted otherwise deletion is blocked. Note that deleting the transform deletes the underlying key which makes decoding of tokenized values impossible without restoring from a backup.',
  }),
  convergent: attr('boolean', {
    label: 'Use convergent tokenization',
    subText:
      "This cannot be edited later. If checked, tokenization of the same plaintext more than once results in the same token. Defaults to false as unique tokens are more desirable from a security standpoint if there isn't a use-case need for convergence.",
  }),
  stores: attr('array', {
    label: 'Stores',
    editType: 'stringArray',
    subText:
      "The list of tokenization stores to use for tokenization state. Vault's internal storage is used by default.",
  }),
  mapping_mode: attr('string', {
    defaultValue: 'default',
    subText:
      'Specifies the mapping mode for stored tokenization values. "default" is strongly recommended for highest security, "exportable" allows for all plaintexts to be decoded via the export-decoded endpoint in an emergency.',
  }),
  max_ttl: attr({
    editType: 'ttl',
    defaultValue: '0',
    label: 'Maximum TTL of a token',
    helperTextDisabled: 'If "0" or unspecified, tokens may have no expiration.',
  }),

  transformAttrs: computed('type', function () {
    // allowed_roles not included so it displays at the bottom of the form
    const baseAttrs = ['name', 'type', 'deletion_allowed'];
    switch (this.type) {
      case 'fpe':
        return [...baseAttrs, 'tweak_source', 'template', 'allowed_roles'];
      case 'masking':
        return [...baseAttrs, 'masking_character', 'template', 'allowed_roles'];
      case 'tokenization':
        return [...baseAttrs, 'mapping_mode', 'convergent', 'max_ttl', 'stores', 'allowed_roles'];
      default:
        return [...baseAttrs];
    }
  }),

  transformFieldAttrs: computed('transformAttrs', function () {
    return expandAttributeMeta(this, this.transformAttrs);
  }),

  backend: attr('string', {
    readOnly: true,
  }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/transformation/${'id'}`, 'backend', 'id'),
});
