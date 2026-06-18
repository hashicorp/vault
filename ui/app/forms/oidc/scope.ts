/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/app-types';
import type { OidcWriteScopeRequest } from '@hashicorp/vault-client-typescript';

type OidcScopeFormData = OidcWriteScopeRequest & {
  name: string;
};

export default class OidcScopeForm extends Form<OidcScopeFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string', { editDisabled: true }),
      new FormField('description', 'string', { editType: 'textarea' }),
      new FormField('template', 'string', { label: 'JSON Template', editType: 'json', mode: 'ruby' }),
    ]),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
  };
}
