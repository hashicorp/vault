/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/app-types';

type RoleData = {
  name?: string;
  transformations?: string[];
  backend?: string;
};

export default class RoleForm extends Form<RoleData> {
  idPrefix = 'role/';

  formFields = [
    new FormField('name', 'string', {
      editDisabled: true,
      subText: 'The name for your role. This cannot be edited later.',
    }),
    new FormField('transformations', 'array', {
      isSectionHeader: true,
      label: 'Transformations',
      subText: 'Select which transformations this role will have access to. It must already exist.',
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
    transformations: [{ type: 'presence', message: 'At least one transformation is required.' }],
  };
}
