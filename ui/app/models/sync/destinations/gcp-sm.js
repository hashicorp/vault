/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
// displayFields are used on the destination details view
const displayFields = [
  // connection details
  'name',
  'projectId',
  'credentials',
  // vault sync config options
  'granularity',
  'secretNameTemplate',
  'customTags',
];
// formFieldGroups are used on the create-edit destination view
const formFieldGroups = [
  { default: ['name', 'projectId'] },
  { Credentials: ['credentials'] },
  { 'Advanced configuration': ['granularity', 'secretNameTemplate', 'customTags'] },
];
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsGoogleCloudSecretManagerModel extends SyncDestinationModel {
  @attr('string', {
    label: 'Project ID',
    subText:
      'The target project to manage secrets in. If set, overrides the project derived from the service account JSON credentials or application default credentials.',
  })
  projectId;

  @attr('string', {
    label: 'JSON credentials',
    subText:
      'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
    editType: 'file',
    docLink: '/vault/docs/secrets/gcp#authentication',
  })
  credentials; // obfuscated, never returned by API. Masking handled by EnableInput component

  @attr('object', {
    subText:
      'An optional set of informational key-value pairs added as additional metadata on secrets synced to this destination. Custom tags are merged with built-in tags.',
    editType: 'kv',
  })
  customTags;
}
