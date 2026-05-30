/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

type PolicyFormData = {
  name: string;
  policy: string;
  enforcement_level?: string;
  paths?: string[];
};

export default class PolicyForm extends Form<PolicyFormData> {
  @tracked declare policyType: string;

  formFields = [
    new FormField('name', 'string', {
      label: 'Policy name',
      placeholder: 'Enter the policy name',
    }),
    new FormField('policy', undefined, { editType: 'yield', label: '' }),
  ];
  fieldProps = ['formFields', 'additionalFields'];

  get additionalFields() {
    if (this.policyType === 'acl') {
      return [];
    }

    const fields = [
      new FormField('enforcement_level', 'string', {
        possibleValues: ['hard-mandatory', 'soft-mandatory', 'advisory'],
        label: 'Enforcement level',
      }),
    ];

    if (this.policyType == 'egp') {
      fields.push(new FormField('paths', 'string', { editType: 'stringArray' }));
    }

    return fields;
  }
}
