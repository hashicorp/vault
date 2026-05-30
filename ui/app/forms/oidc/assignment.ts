/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import isPresent from '@ember/utils/lib/is_present';

import type { Validations } from 'vault/app-types';
import type { OidcWriteAssignmentRequest } from '@hashicorp/vault-client-typescript';

type OidcAssignmentFormData = OidcWriteAssignmentRequest & {
  name: string;
};

export default class OidcAssignmentForm extends Form<OidcAssignmentFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string', { editDisabled: true }),
      // SearchSelect within the FormField component works in conjunction with Ember Data Models
      // we can still use the component since it supports passing in an array of objects as options for the select
      // yield out the fields so scopes can be fetched in the route and passed directly to SearchSelect
      new FormField('entity_ids', undefined, {
        label: 'Entities',
        editType: 'yield',
      }),
      new FormField('group_ids', undefined, {
        label: 'Groups',
        editType: 'yield',
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
    targets: [
      {
        validator(data: OidcAssignmentFormData) {
          return isPresent(data.entity_ids) || isPresent(data.group_ids);
        },
        message: 'At least one entity or group is required.',
      },
    ],
  };
}
