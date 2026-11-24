/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/app-types';

type TemplateData = {
  name?: string;
  type?: string;
  pattern?: string;
  alphabet?: string[];
  encode_format?: string;
  decode_formats?: Record<string, string>;
  backend?: string;
};

export default class TemplateForm extends Form<TemplateData> {
  idPrefix = 'template/';

  formFields = [
    new FormField('name', 'string', {
      editDisabled: true,
      subText:
        'Templates allow Vault to determine what and how to capture the value to be transformed. This cannot be edited later.',
    }),
    new FormField('pattern', 'string', {
      editType: 'regex',
      subText: "The template's pattern defines the data format. Expressed in regex.",
    }),
    new FormField('encode_format', 'string'),
    new FormField('decode_formats', 'object'),
    new FormField('alphabet', 'array', {
      isSectionHeader: true,
      label: 'Alphabet',
      subText:
        'Alphabet defines a set of characters (UTF-8) that is used for FPE to determine the validity of plaintext and ciphertext values. You can choose a built-in one, or create your own.',
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
    pattern: [{ type: 'presence', message: 'Pattern is required.' }],
  };
}
