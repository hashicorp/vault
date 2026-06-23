/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

interface GroupIdentityFormData {
  name?: string;
  type?: string;
  policies?: string[];
  metadata?: Record<string, string>;
  member_group_ids?: string[];
  member_entity_ids?: string[];
}

export default class GroupIdentityForm extends Form<GroupIdentityFormData> {
  identityFormType = 'group';

  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string'),
      new FormField('type', 'string', {
        possibleValues: ['internal', 'external'],
      }),
      new FormField('policies', undefined, {
        editType: 'yield',
        isSectionHeader: true,
      }),
      new FormField('metadata', 'object', {
        editType: 'kv',
        isSectionHeader: true,
      }),
      new FormField('member_group_ids', undefined, {
        editType: 'yield',
        isSectionHeader: true,
      }),
      new FormField('member_entity_ids', undefined, {
        editType: 'yield',
        isSectionHeader: true,
      }),
    ]),
  ];
}
