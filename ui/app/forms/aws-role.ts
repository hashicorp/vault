/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { Validations } from 'vault/app-types';

interface RoleAwsData {
  name?: string;
  credential_type?: string;
  role_arns?: string[];
  policy_arns?: string[];
  policy_document?: string;
  backend?: string;
}

const CREDENTIAL_TYPE_OPTIONS = [
  { value: 'iam_user', displayName: 'IAM User' },
  { value: 'assumed_role', displayName: 'Assumed Role' },
  { value: 'federation_token', displayName: 'Federation Token' },
  { value: 'session_token', displayName: 'Session Token' },
];

export default class RoleAwsForm extends Form<RoleAwsData> {
  get id() {
    return this.data.name;
  }

  get backend() {
    return this.data.backend;
  }

  get formFields(): FormField[] {
    const credentialType = this.data.credential_type || 'iam_user';

    const baseFields = [
      new FormField('name', 'string', {
        label: 'Role name',
        editDisabled: !this.isNew,
      }),
      new FormField('credential_type', 'string', {
        label: 'Credential type',
        possibleValues: CREDENTIAL_TYPE_OPTIONS,
        defaultValue: 'iam_user',
      }),
    ];

    const conditionalFields: FormField[] = [];

    if (credentialType === 'iam_user') {
      conditionalFields.push(
        new FormField('policy_arns', undefined, {
          editType: 'stringArray',
          label: 'Policy ARNs',
        }),
        new FormField('policy_document', 'string', {
          editType: 'json',
          label: 'Policy Document',
          helpText:
            'A policy is an object in AWS that, when associated with an identity or resource, defines their permissions.',
        })
      );
    } else if (credentialType === 'assumed_role') {
      conditionalFields.push(
        new FormField('role_arns', undefined, {
          editType: 'stringArray',
          label: 'Role ARNs',
        }),
        new FormField('policy_document', 'string', {
          editType: 'json',
          label: 'Policy Document',
          helpText:
            'A policy is an object in AWS that, when associated with an identity or resource, defines their permissions.',
        })
      );
    } else if (credentialType === 'federation_token') {
      conditionalFields.push(
        new FormField('policy_document', 'string', {
          editType: 'json',
          label: 'Policy Document',
          helpText:
            'A policy is an object in AWS that, when associated with an identity or resource, defines their permissions.',
        })
      );
    }

    return [...baseFields, ...conditionalFields];
  }

  validations: Validations = {
    name: [{ type: 'presence', message: "Role name can't be blank." }],
    credential_type: [{ type: 'presence', message: "Credential type can't be blank." }],
  };

  toJSON() {
    const { credential_type, policy_document } = this.data;
    const role_arns = credential_type === 'assumed_role' ? this.data.role_arns || [] : [];
    const policy_arns = credential_type === 'iam_user' ? this.data.policy_arns || [] : [];
    const data = {
      credential_type,
      role_arns,
      policy_arns,
      policy_document: policy_document === '{}' ? '' : policy_document,
    };
    return { ...super.toJSON(), data };
  }
}
