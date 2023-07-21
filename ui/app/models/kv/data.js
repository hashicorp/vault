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
  path: [{ type: 'presence', message: `Path can't be blank.` }],
};
const formFieldProps = ['path', 'data'];

@withModelValidations(validations)
@withFormFields(formFieldProps)
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

  @attr('number', {
    defaultValue: 0, // version 0 only occurs on creating a secret
  })
  casVersion;

  // Permissions
  @lazyCapabilities(apiPath`${'backend'}/data/${'path'}`, 'backend', 'path') dataPath;
  @lazyCapabilities(apiPath`${'backend'}/metadata/${'path'}`, 'backend', 'path') metadataPath;
  @lazyCapabilities(apiPath`${'backend'}/delete/${'path'}`, 'backend', 'path') deletePath;
  @lazyCapabilities(apiPath`${'backend'}/destroy/${'path'}`, 'backend', 'path') destroyPath;
  @lazyCapabilities(apiPath`${'backend'}/undelete/${'path'}`, 'backend', 'path') undeletePath;

  get canDeleteLatestVersion() {
    return this.dataPath.get('canDelete');
  }
  get canDeleteVersion() {
    return this.deletePath.get('canUpdate');
  }
  get canUndelete() {
    return this.undeletePath.get('canUpdate');
  }
  get canDestroyVersion() {
    return this.destroyPath.get('canUpdate');
  }
  get canEditData() {
    return this.dataPath.get('canUpdate');
  }
  get canReadData() {
    return this.dataPath.get('canRead');
  }
  get canReadMetadata() {
    return this.metadataPath.get('canRead');
  }
  get canUpdateMetadata() {
    return this.metadataPath.get('canUpdate');
  }
  get canListMetadata() {
    return this.metadataPath.get('canList');
  }
}
