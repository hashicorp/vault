/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

interface EntityIdentityFormData {
  name?: string;
  disabled?: boolean;
  policies?: string[];
  metadata?: Record<string, string>;
}

export default class EntityIdentityForm extends Form<EntityIdentityFormData> {
  identityFormType = 'entity';

  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string'),
      new FormField('disabled', 'boolean', {
        label: 'Disable entity',
        helpText: 'All associated tokens cannot be used, but are not revoked.',
      }),
      new FormField('policies', undefined, {
        editType: 'yield',
        isSectionHeader: true,
      }),
      new FormField('metadata', 'object', {
        editType: 'kv',
        isSectionHeader: true,
      }),
    ]),
  ];
}
