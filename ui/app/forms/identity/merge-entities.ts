/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

interface MergeEntitiesFormData {
  to_entity_id?: string;
  from_entity_ids?: string[];
  force?: boolean;
}

export default class MergeEntitiesForm extends Form<MergeEntitiesFormData> {
  identityFormType = 'merge-entities';

  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('to_entity_id', 'string', {
        label: 'Entity ID to merge to',
      }),
      new FormField('from_entity_ids', 'string', {
        label: 'Entity ID to merge from',
        editType: 'stringArray',
      }),
      new FormField('force', 'boolean', {
        label: 'Keep MFA secrets from the "to" entity if there are merge conflicts',
      }),
    ]),
  ];
}
