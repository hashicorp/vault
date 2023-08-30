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
  transformAttrs: computed('type', function () {
    if (this.type === 'masking') {
      return ['name', 'type', 'masking_character', 'template', 'allowed_roles'];
    }
    return ['name', 'type', 'tweak_source', 'template', 'allowed_roles'];
  }),
  transformFieldAttrs: computed('transformAttrs', function () {
    return expandAttributeMeta(this, this.transformAttrs);
  }),

  backend: attr('string', {
    readOnly: true,
  }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/transformation/${'id'}`, 'backend', 'id'),
});
