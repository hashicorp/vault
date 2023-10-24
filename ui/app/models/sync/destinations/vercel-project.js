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
  deploymentEnvironments: [{ type: 'presence', message: 'Deployment environments are required.' }],
};
const fields = ['name', 'accessToken', 'projectId', 'teamId', 'deploymentEnvironments'];

@withModelValidations(validations)
@withFormFields(fields)
export default class SyncDestinationsVercelProjectModel extends SyncDestinationModel {
  @attr('string', {
    subText: 'Vercel API access token with the permissions to manage environment variables.',
  })
  accessToken;

  @attr('string', {
    label: 'Project ID',
    subText: 'Project ID where to manage environment variables.',
  })
  projectId;

  @attr('string', {
    label: 'Team ID',
    subText: 'Team ID the project belongs to. Optional.',
  })
  teamId;

  @attr({
    subText: 'Deployment environments where the environment variables are available.',
    editType: 'yield',
    checkboxOptions: ['deployment', 'preview', 'production'],
    // TODO can also be a string, return and test how this works with live API
    defaultValue: () => [],
  })
  deploymentEnvironments;
}
