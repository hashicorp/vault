/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { commonFields, getPayload } from './shared';

import type { SystemWriteSyncDestinationsVercelProjectNameRequest } from '@hashicorp/vault-client-typescript';

type VercelProjectFormData = SystemWriteSyncDestinationsVercelProjectNameRequest & {
  name: string;
};

export default class VercelProjectForm extends Form<VercelProjectFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      commonFields.name,
      new FormField('project_id', 'string', {
        label: 'Project ID',
        subText: 'Project ID where to manage environment variables.',
        editDisabled: true,
      }),
      new FormField('team_id', 'string', {
        label: 'Team ID',
        subText: 'Team ID the project belongs to. Optional.',
      }),
      new FormField('deployment_environments', 'string', {
        subText: 'Deployment environments where the environment variables are available.',
        editType: 'checkboxList',
        possibleValues: ['development', 'preview', 'production'],
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('access_token', 'string', {
        subText: 'Vercel API access token with the permissions to manage environment variables.',
        sensitive: true,
        noCopy: true,
      }),
    ]),
    new FormFieldGroup('Advanced configuration', [commonFields.granularity, commonFields.secretNameTemplate]),
  ];

  toJSON() {
    const formState = super.toJSON();
    const data = getPayload<VercelProjectFormData>('vercel-project', this.data, this.isNew);
    return { ...formState, data };
  }
}
