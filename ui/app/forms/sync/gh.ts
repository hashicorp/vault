/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { commonFields, getPayload } from './shared';

import type { SystemWriteSyncDestinationsGhNameRequest } from '@hashicorp/vault-client-typescript';

type GhFormData = SystemWriteSyncDestinationsGhNameRequest & {
  name: string;
};

export default class GcpSmForm extends Form<GhFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      commonFields.name,
      new FormField('repositoryOwner', 'string', {
        subText:
          'Github organization or username that owns the repository. If empty, Vault will use the GITHUB_REPOSITORY_OWNER environment variable if configured.',
        editDisabled: true,
      }),
      new FormField('repositoryName', 'string', {
        subText:
          'The name of the Github repository to connect to. If empty, Vault will use the GITHUB_REPOSITORY_NAME environment variable if configured.',
        editDisabled: true,
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('accessToken', 'string', {
        subText:
          'Personal access token to authenticate to the GitHub repository. If empty, Vault will use the GITHUB_ACCESS_TOKEN environment variable if configured.',
        sensitive: true,
        noCopy: true,
      }),
    ]),
    new FormFieldGroup('Advanced configuration', [commonFields.granularity, commonFields.secretNameTemplate]),
  ];

  toJSON() {
    const formState = super.toJSON();
    const data = getPayload<GhFormData>('gh', this.data, this.isNew);
    return { ...formState, data };
  }
}
