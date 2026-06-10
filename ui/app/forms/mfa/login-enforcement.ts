/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/app-types';
import type { MfaWriteLoginEnforcementRequest } from '@hashicorp/vault-client-typescript';

type MfaLoginEnforcementFormData = MfaWriteLoginEnforcementRequest & {
  name: string;
};

export default class MfaLoginEnforcementForm extends Form<MfaLoginEnforcementFormData> {
  formFieldGroups = [
    new FormFieldGroup('name', [
      new FormField('name', 'string', {
        label: 'Name',
        subText:
          'The name for this enforcement. Giving it a name means that you can refer to it again later. This name will not be editable later.',
        editDisabled: !this.isNew,
      }),
    ]),
    new FormFieldGroup('mfa_methods', [
      new FormField('mfa_methods', 'array', {
        label: 'MFA methods',
        subText: 'The MFA method(s) that this enforcement will apply to.',
        editType: 'yield',
      }),
    ]),
    new FormFieldGroup('targets', [
      new FormField('auth_method_accessors', 'array', {
        editType: 'yield',
      }),
      new FormField('auth_method_types', 'array', {
        editType: 'yield',
      }),
      new FormField('identity_entity_ids', 'array', {
        editType: 'yield',
      }),
      new FormField('identity_group_ids', 'array', {
        editType: 'yield',
      }),
    ]),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required' }],
    mfa_methods: [{ type: 'presence', message: 'At least one MFA method is required' }],
    targets: [
      {
        validator(model: MfaLoginEnforcementFormData) {
          // Check if at least one target is specified
          return !!(
            (model.auth_method_accessors && model.auth_method_accessors.length > 0) ||
            (model.auth_method_types && model.auth_method_types.length > 0) ||
            (model.identity_entity_ids && model.identity_entity_ids.length > 0) ||
            (model.identity_group_ids && model.identity_group_ids.length > 0)
          );
        },
        message:
          "At least one target is required. If you've selected one, click 'Add' to make sure it's added to this enforcement.",
      },
    ],
  };
}
