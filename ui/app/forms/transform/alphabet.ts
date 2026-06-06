/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/app-types';

type AlphabetData = {
  name?: string;
  alphabet?: string;
  backend?: string;
};

export default class AlphabetForm extends Form<AlphabetData> {
  // Required by secret-edit-layout.hbs to resolve the correct editComponent via options-for-backend
  idPrefix = 'alphabet/';
  formFields = [
    new FormField('name', 'string', {
      editDisabled: true,
      subText: 'The alphabet name. Keep in mind that spaces are not allowed and this cannot be edited later.',
    }),
    new FormField('alphabet', 'string', {
      label: 'Alphabet',
      subText:
        'Provide the set of valid UTF-8 characters contained within both the input and transformed value.',
      docLink: '/vault/api-docs/secret/transform#create-update-alphabet',
    }),
  ];

  validations: Validations = {
    alphabet: [{ type: 'presence', message: 'Alphabet is required.' }],
  };
}
