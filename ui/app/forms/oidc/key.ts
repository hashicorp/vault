/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/app-types';
import type { OidcWriteKeyRequest } from '@hashicorp/vault-client-typescript';

type OidcKeyFormData = OidcWriteKeyRequest & {
  name: string;
};

export default class OidcKeyForm extends Form<OidcKeyFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string', { editDisabled: true }),
      new FormField('algorithm', 'string', {
        possibleValues: ['RS256', 'RS384', 'RS512', 'ES256', 'ES384', 'ES512', 'EdDSA'],
      }),
      new FormField('rotation_period', undefined, {
        editType: 'ttl',
      }),
      new FormField('verification_ttl', undefined, {
        label: 'Verification TTL',
        editType: 'ttl',
      }),
    ]),
  ];

  validations: Validations = {
    name: [
      { type: 'presence', message: 'Name is required.' },
      {
        type: 'containsWhiteSpace',
        message: 'Name cannot contain whitespace.',
      },
    ],
  };
}
