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
  'region',
  'accessKeyId',
  'secretAccessKey',
  'roleArn',
  'externalId',
  // sync config options
  'granularity',
  'secretNameTemplate',
  'customTags',
];
// formFieldGroups are used on the create-edit destination view
const formFieldGroups = [
  {
    default: ['name', 'region', 'roleArn', 'externalId'],
  },
  { Credentials: ['accessKeyId', 'secretAccessKey'] },
  { 'Advanced configuration': ['granularity', 'secretNameTemplate', 'customTags'] },
];
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsAwsSecretsManagerModel extends SyncDestinationModel {
  @attr('string', {
    label: 'Access key ID',
    subText:
      'Access key ID to authenticate against the secrets manager. If empty, Vault will use the AWS_ACCESS_KEY_ID environment variable if configured.',
    sensitive: true,
    noCopy: true,
  })
  accessKeyId; // obfuscated, never returned by API

  @attr('string', {
    label: 'Secret access key',
    subText:
      'Secret access key to authenticate against the secrets manager. If empty, Vault will use the AWS_SECRET_ACCESS_KEY environment variable if configured.',
    sensitive: true,
    noCopy: true,
  })
  secretAccessKey; // obfuscated, never returned by API

  @attr('string', {
    subText:
      'For AWS secrets manager, the name of the region must be supplied, something like “us-west-1.” If empty, Vault will use the AWS_REGION environment variable if configured.',
    editDisabled: true,
  })
  region;

  @attr('object', {
    subText:
      'An optional set of informational key-value pairs added as additional metadata on secrets synced to this destination. Custom tags are merged with built-in tags.',
    editType: 'kv',
  })
  customTags;

  @attr('string', {
    label: 'Role ARN',
    subText:
      'Specifies a role to assume when connecting to AWS. When assuming a role, Vault uses temporary STS credentials to authenticate.',
  })
  roleArn;

  @attr('string', {
    label: 'External ID',
    subText:
      'Optional extra protection that must match the trust policy granting access to the AWS IAM role ARN. We recommend using a different random UUID per destination.',
  })
  externalId;
}
