/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/app-types';

const TYPES = [
  { value: 'fpe', displayName: 'Format Preserving Encryption (FPE)' },
  { value: 'masking', displayName: 'Masking' },
  { value: 'tokenization', displayName: 'Tokenization' },
];

const TWEAK_SOURCE = [
  { value: 'supplied', displayName: 'supplied' },
  { value: 'generated', displayName: 'generated' },
  { value: 'internal', displayName: 'internal' },
];

type TransformationData = {
  name?: string;
  type?: string;
  tweak_source?: string;
  masking_character?: string;
  template?: string[];
  allowed_roles?: string[];
  deletion_allowed?: boolean;
  convergent?: boolean;
  stores?: string[];
  mapping_mode?: string;
  max_ttl?: string;
  backend?: string;
};

export default class TransformationForm extends Form<TransformationData> {
  formFields = [
    new FormField('name', 'string', {
      editDisabled: true,
      subText: 'The name for your transformation. This cannot be edited later.',
    }),
    new FormField('type', 'string', {
      editDisabled: true,
      possibleValues: TYPES,
      defaultValue: 'fpe',
      label: 'Type',
      subText:
        'Vault provides two types of transformations: Format Preserving Encryption (FPE) is reversible, while Masking is not. This cannot be edited later.',
    }),
    new FormField('deletion_allowed', 'boolean', {
      label: 'Allow deletion',
      subText:
        'If checked, this transform can be deleted otherwise deletion is blocked. Note that deleting the transform deletes the underlying key which makes decoding of tokenized values impossible without restoring from a backup.',
    }),
    new FormField('tweak_source', 'string', {
      possibleValues: TWEAK_SOURCE,
      defaultValue: 'supplied',
      label: 'Tweak source',
      subText:
        'A tweak value is used when performing FPE transformations. This can be supplied, generated, or internal.',
    }),
    new FormField('masking_character', 'string', {
      defaultValue: '*',
      label: 'Masking character',
      subText: 'Specify which character you\u2019d like to mask your data.',
    }),
    new FormField('template', 'array', {
      isSectionHeader: true,
      label: 'Template',
      subText:
        'Templates allow Vault to determine what and how to capture the value to be transformed. Type to use an existing template or create a new one.',
    }),
    new FormField('allowed_roles', 'array', {
      isSectionHeader: true,
      label: 'Allowed roles',
      subText: 'Search for an existing role, type a new role to create it, or use a wildcard (*).',
    }),
    new FormField('mapping_mode', 'string', {
      defaultValue: 'default',
      label: 'Mapping mode',
      subText:
        'Specifies the mapping mode for stored tokenization values. "default" is strongly recommended for highest security, "exportable" allows for all plaintexts to be decoded via the export-decoded endpoint in an emergency.',
    }),
    new FormField('convergent', 'boolean', {
      label: 'Use convergent tokenization',
      subText:
        'If checked, tokenization of the same plaintext more than once results in the same token. Defaults to false as unique tokens are more desirable from a security standpoint.',
    }),
    new FormField('max_ttl', 'string', {
      editType: 'ttl',
      defaultValue: '0',
      label: 'Maximum TTL (time-to-live) of a token',
      subText: 'If \u201c0\u201d or unspecified, tokens may have no expiration.',
    }),
    new FormField('stores', 'array', {
      editType: 'stringArray',
      label: 'Stores',
      subText:
        'The list of tokenization stores to use for tokenization state. Vault\u2019s internal storage is used by default.',
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
    template: [{ type: 'presence', message: 'Template is required.' }],
  };
}
