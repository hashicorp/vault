/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

/* sample response
{
  "data": {
    "data": {
      "foo": "bar"
    },
    "metadata": {
      "created_time": "2018-03-22T02:24:06.945319214Z",
      "custom_metadata": {
        "owner": "jdoe",
        "mission_critical": "false"
      },
      "deletion_time": "",
      "destroyed": false,
      "version": 2
    }
  }
}
*/

const validations = {
  path: [
    { type: 'presence', message: `Path can't be blank.` },
    { type: 'endsInSlash', message: `Path can't end in forward slash '/'.` },
    {
      type: 'containsWhiteSpace',
      message:
        "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.",
      level: 'warn',
    },
  ],
  secretData: [
    {
      validator: (model) =>
        model.secretData !== undefined && typeof model.secretData !== 'object' ? false : true,
      message: 'Vault expects data to be formatted as an JSON object.',
    },
  ],
};
@withModelValidations(validations)
@withFormFields()
export default class KvSecretDataModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord.
  @attr('string', { label: 'Path for this secret' }) path;
  @attr('object') secretData; // { key: value } data of the secret version

  // Params returned on the GET response.
  @attr('string') createdTime;
  @attr('object') customMetadata;
  @attr('string') deletionTime;
  @attr('boolean') destroyed;
  @attr('number') version;
  // Set in adapter if read failed
  @attr('number') failReadErrorCode;

  // if creating a new version this value is set in the edit route's
  // model hook from metadata or secret version, pending permissions
  // if the value is not a number, don't send options.cas on payload
  @attr('number')
  casVersion;

  get state() {
    if (this.destroyed) return 'destroyed';
    if (this.deletionTime) return 'deleted';
    if (this.createdTime) return 'created';
    return '';
  }

  // Permissions
  @lazyCapabilities(apiPath`${'backend'}/data/${'path'}`, 'backend', 'path') dataPath;
  @lazyCapabilities(apiPath`${'backend'}/metadata/${'path'}`, 'backend', 'path') metadataPath;
  @lazyCapabilities(apiPath`${'backend'}/delete/${'path'}`, 'backend', 'path') deletePath;
  @lazyCapabilities(apiPath`${'backend'}/destroy/${'path'}`, 'backend', 'path') destroyPath;
  @lazyCapabilities(apiPath`${'backend'}/undelete/${'path'}`, 'backend', 'path') undeletePath;

  get canDeleteLatestVersion() {
    return this.dataPath.get('canDelete') !== false;
  }
  get canDeleteVersion() {
    return this.deletePath.get('canUpdate') !== false;
  }
  get canUndelete() {
    return this.undeletePath.get('canUpdate') !== false;
  }
  get canDestroyVersion() {
    return this.destroyPath.get('canUpdate') !== false;
  }
  get canEditData() {
    return this.dataPath.get('canUpdate') !== false;
  }
  get canReadData() {
    return this.dataPath.get('canRead') !== false;
  }
  get canReadMetadata() {
    return this.metadataPath.get('canRead') !== false;
  }
  get canUpdateMetadata() {
    return this.metadataPath.get('canUpdate') !== false;
  }
  get canListMetadata() {
    return this.metadataPath.get('canList') !== false;
  }
  get canDeleteMetadata() {
    return this.metadataPath.get('canDelete') !== false;
  }
}
