/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/vault/app-types';

export type PkiKeyFormData = {
  key_name?: string;
  key_id?: string;
  type: 'internal' | 'exported';
  key_type: 'rsa' | 'ec' | 'ed25519';
  key_bits: number;
};

export default class PkiKeyForm extends Form<PkiKeyFormData> {
  validations: Validations = {
    type: [{ type: 'presence', message: 'Type is required.' }],
    key_type: [{ type: 'presence', message: 'Please select a key type.' }],
    key_name: [
      {
        type: 'isNot',
        options: { value: 'default' },
        message: `Key name cannot be the reserved value 'default'`,
      },
    ],
  };

  keyNameField = new FormField('key_name', 'string', {
    subText: `Optional, human-readable name for this key. The name must be unique across all keys and cannot be 'default'.`,
  });
  keyTypeField = new FormField('key_type', 'string', {
    noDefault: true,
    possibleValues: ['rsa', 'ec', 'ed25519'],
    subText: 'The type of key that will be generated. Must be rsa, ed25519, or ec. ',
  });

  formFields = [this.keyNameField, this.keyTypeField];

  formFieldGroups = [
    new FormFieldGroup('default', [
      this.keyNameField,
      new FormField('type', 'string', {
        noDefault: true,
        possibleValues: ['internal', 'exported'],
        subText:
          'The type of operation. If exported, the private key will be returned in the response; if internal the private key will not be returned and cannot be retrieved later.',
      }),
    ]),
    new FormFieldGroup('Key parameters', [
      this.keyTypeField,
      new FormField('key_bits', 'number', {
        label: 'Key bits',
        noDefault: true,
        subText: 'Bit length of the key to generate.',
      }),
    ]),
  ];

  toJSON() {
    const formState = super.toJSON();
    if (!this.isNew) {
      // the only editable property is key_name which is optional so the form will always be valid
      return { ...formState, isValid: true, state: {}, invalidFormMessage: '' };
    }
    return formState;
  }
}
