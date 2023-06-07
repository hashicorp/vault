/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

/* sample response
{
  "data": {
    "cas_required": false,
    "created_time": "2018-03-22T02:24:06.945319214Z",
    "current_version": 3,
    "delete_version_after": "3h25m19s",
    "max_versions": 0,
    "oldest_version": 0,
    "updated_time": "2018-03-22T02:36:43.986212308Z",
    "custom_metadata": {
      "foo": "abc",
      "bar": "123",
      "baz": "5c07d823-3810-48f6-a147-4c06b5219e84"
    },
    "versions": {
      "1": {
        "created_time": "2018-03-22T02:24:06.945319214Z",
        "deletion_time": "",
        "destroyed": false
      },
      "2": {
        "created_time": "2018-03-22T02:36:33.954880664Z",
        "deletion_time": "",
        "destroyed": false
      },
      "3": {
        "created_time": "2018-03-22T02:36:43.986212308Z",
        "deletion_time": "",
        "destroyed": false
      }
    }
  }
}
*/

const validations = {
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};
const formFieldProps = ['path', 'data'];

@withModelValidations(validations)
@withFormFields(formFieldProps)
export default class KvSecretMetadataModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord.

  @attr('number', {
    defaultValue: 0,
    label: 'Maximum number of versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
  })
  maxVersions;

  @attr('number', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText: `Writes will only be allowed if the key's current version matches the version specified in the cas parameter.`,
  })
  casRequired;

  @attr('string', {
    defaultValue: '0s',
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: `A secret's version must be manually deleted.`,
    helperTextEnabled: 'Delete all new versions of this secret after.',
  })
  deleteVersionAfter;

  @attr('object', {
    editType: 'kv',
    subText: 'An optional set of informational key-value pairs that will be stored with all secret versions.',
  })
  customMetadata;

  // Additional Params only returned on the GET response.
  @attr('string') createdTime;
  @attr('number') currentVersion;
  @attr('number') oldestVersion;
  @attr('string') updatedTime;
  @attr('object') versions;
}
