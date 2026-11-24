/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { DestinationType } from 'sync/utils/constants';
import CreateDestinationForm from './create-destination';

import type { SystemWriteSyncDestinationsGhNameRequest } from '@hashicorp/vault-client-typescript';

type GhFormData = SystemWriteSyncDestinationsGhNameRequest & {
  name: string;
};

export default class GhForm extends CreateDestinationForm<GhFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      this.commonFields.name,
      new FormField('repository_owner', 'string', {
        subText:
          'Github organization or username that owns the repository. If empty, Vault will use the GITHUB_REPOSITORY_OWNER environment variable if configured.',
        editDisabled: true,
      }),
      new FormField('repository_name', 'string', {
        subText:
          'The name of the Github repository to connect to. If empty, Vault will use the GITHUB_REPOSITORY_NAME environment variable if configured.',
        editDisabled: true,
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('access_token', 'string', {
        subText:
          'Personal access token to authenticate to the GitHub repository. If empty, Vault will use the GITHUB_ACCESS_TOKEN environment variable if configured.',
        sensitive: true,
        noCopy: true,
      }),
    ]),
    new FormFieldGroup('Advanced configuration', [
      this.commonFields.granularity,
      this.commonFields.secretNameTemplate,
    ]),
  ];

  toJSON() {
    const formState = super.toJSON();
    const data = this.getPayload<GhFormData>(DestinationType.Gh, this.data, this.isNew);
    return { ...formState, data };
  }
}
