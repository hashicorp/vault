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
  teamId: [
    {
      validator: (model) =>
        !model.isNew && Object.keys(model.changedAttributes()).includes('teamId') ? false : true,
      message: 'Team ID should only be updated if the project was transferred to another account.',
      level: 'warn',
    },
  ],
  // getter/setter for the deploymentEnvironments model attribute
  deploymentEnvironmentsArray: [{ type: 'presence', message: 'At least one environment is required.' }],
};

const displayFields = [
  'name',
  'accessToken',
  'projectId',
  'teamId',
  'deploymentEnvironments',
  'secretNameTemplate',
];
const formFieldGroups = [
  { default: ['name', 'projectId', 'teamId', 'deploymentEnvironments', 'secretNameTemplate'] },
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

  // commaString transforms param from the server's array type
  // to a comma string so changedAttributes() will track changes
  @attr('commaString', {
    subText: 'Deployment environments where the environment variables are available.',
    editType: 'checkboxList',
    possibleValues: ['development', 'preview', 'production'],
    fieldValue: 'deploymentEnvironmentsArray', // getter/setter used to update value
  })
  deploymentEnvironments;

  // Arrays are easier for managing multi-option selection
  // these get/set the deploymentEnvironments attribute via arrays
  get deploymentEnvironmentsArray() {
    // if undefined or an empty string, return empty array
    return !this.deploymentEnvironments ? [] : this.deploymentEnvironments.split(',');
  }

  set deploymentEnvironmentsArray(value) {
    this.deploymentEnvironments = value.join(',');
  }
}
