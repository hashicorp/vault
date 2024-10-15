/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { keyIsFolder } from 'core/utils/key-utils';
import { isDeleted } from 'kv/utils/kv-deleted';

const validations = {
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};
const formFieldProps = ['customMetadata', 'maxVersions', 'casRequired', 'deleteVersionAfter'];

@withModelValidations(validations)
@withFormFields(formFieldProps)
export default class KvSecretMetadataModel extends Model {
  @attr('string') backend;
  @attr('string') path;
  @attr('string') fullSecretPath;

  @attr('number', {
    defaultValue: 0,
    label: 'Maximum number of versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
  })
  maxVersions;

  @attr('boolean', {
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
    helperTextEnabled: 'Delete all new versions of this secret after:',
  })
  deleteVersionAfter;

  // the API returns custom_metadata: null if empty but because the attr is an 'object' ember data transforms it to an empty object.
  // this is important because we rely on the empty object as a truthy value in template conditionals
  @attr('object', {
    editType: 'kv',
    isSectionHeader: true,
    subText: 'An optional set of informational key-value pairs that will be stored with all secret versions.',
  })
  customMetadata;

  // Additional Params only returned on the GET response.
  @attr('string') createdTime;
  @attr('number') currentVersion;
  @attr('number') oldestVersion;
  @attr('string') updatedTime;
  @attr('object') versions;

  // used for KV list and list-directory view
  get pathIsDirectory() {
    // ex: beep/
    return keyIsFolder(this.path);
  }

  // turns version object into an array for version dropdown menu
  get sortedVersions() {
    const array = [];
    for (const key in this.versions) {
      this.versions[key].isSecretDeleted = isDeleted(this.versions[key].deletion_time);
      array.push({ version: key, ...this.versions[key] });
    }
    // version keys are in order created with 1 being the oldest, we want newest first
    return array.reverse();
  }

  // helps in long logic statements for state of a currentVersion
  get currentSecret() {
    if (!this.versions || !this.currentVersion) return false;
    const data = this.versions[this.currentVersion];
    const state = data.destroyed ? 'destroyed' : isDeleted(data.deletion_time) ? 'deleted' : 'created';
    return {
      state,
      isDeactivated: state !== 'created',
      deletionTime: data.deletion_time,
    };
  }

  get permissionsPath() {
    return this.fullSecretPath || this.path;
  }

  // permissions needed for the list view where kv/data has not yet been called. Allows us to conditionally show action items in the LinkedBlock popups.
  @lazyCapabilities(apiPath`${'backend'}/data/${'permissionsPath'}`, 'backend', 'permissionsPath') dataPath;
  @lazyCapabilities(apiPath`${'backend'}/metadata/${'permissionsPath'}`, 'backend', 'permissionsPath')
  metadataPath;

  get canDeleteMetadata() {
    return this.metadataPath.get('canDelete') !== false;
  }
  get canReadMetadata() {
    return this.metadataPath.get('canRead') !== false;
  }
  get canUpdateMetadata() {
    return this.metadataPath.get('canUpdate') !== false;
  }
  get canCreateVersionData() {
    return this.dataPath.get('canUpdate') !== false;
  }
}
