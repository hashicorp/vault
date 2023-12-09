/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
  accessToken: [{ type: 'presence', message: 'Access token is required.' }],
  projectId: [{ type: 'presence', message: 'Project ID is required.' }],
  deploymentEnvironments: [{ type: 'presence', message: 'At least one environment is required.' }],
};
const displayFields = ['name', 'accessToken', 'projectId', 'teamId', 'deploymentEnvironments'];
const formFieldGroups = [
  { default: ['name', 'projectId', 'teamId', 'deploymentEnvironments'] },
  { Credentials: ['accessToken'] },
];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsVercelProjectModel extends SyncDestinationModel {
  @attr('string', {
    subText: 'Vercel API access token with the permissions to manage environment variables.',
  })
  accessToken; // obfuscated, never returned by API

  @attr('string', {
    label: 'Project ID',
    subText: 'Project ID where to manage environment variables.',
    editDisabled: true,
  })
  projectId;

  @attr('string', {
    label: 'Team ID',
    subText: 'Team ID the project belongs to. Optional.',
  })
  teamId;

  // TODO can also be a string, return and test how this works with live API
  @attr('array', {
    subText: 'Deployment environments where the environment variables are available.',
    editType: 'checkboxList',
    possibleValues: ['development', 'preview', 'production'],
    defaultValue: () => [],
  })
  deploymentEnvironments;
}
