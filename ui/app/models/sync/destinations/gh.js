/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

const displayFields = [
  // connection details
  'name',
  'repositoryOwner',
  'repositoryName',
  'accessToken',
  // vault sync config options
  'granularity',
  'secretNameTemplate',
];
const formFieldGroups = [
  { default: ['name', 'repositoryOwner', 'repositoryName', 'granularity', 'secretNameTemplate'] },
  { Credentials: ['accessToken'] },
];

@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsGithubModel extends SyncDestinationModel {
  @attr('string', {
    subText:
      'Personal access token to authenticate to the GitHub repository. If empty, Vault will use the GITHUB_ACCESS_TOKEN environment variable if configured.',
  })
  accessToken; // obfuscated, never returned by API

  @attr('string', {
    subText:
      'Github organization or username that owns the repository. If empty, Vault will use the GITHUB_REPOSITORY_OWNER environment variable if configured.',
    editDisabled: true,
  })
  repositoryOwner;

  @attr('string', {
    subText:
      'The name of the Github repository to connect to. If empty, Vault will use the GITHUB_REPOSITORY_NAME environment variable if configured.',
    editDisabled: true,
  })
  repositoryName;
}
