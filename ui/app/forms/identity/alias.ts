/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

interface AliasIdentityFormData {
  name?: string;
  mount_accessor?: string;
}

export default class AliasIdentityForm extends Form<AliasIdentityFormData> {
  identityFormType = 'alias';

  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string'),
      new FormField('mount_accessor', 'string', {
        label: 'Auth backend',
        editType: 'mountAccessor',
      }),
    ]),
  ];
}
