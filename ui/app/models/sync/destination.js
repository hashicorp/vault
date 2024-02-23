/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { findDestination } from 'vault/helpers/sync-destinations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';

// Base model for all secret sync destination types
const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    { type: 'containsWhiteSpace', message: 'Name cannot contain whitespace.' },
  ],
};

@withModelValidations(validations)
export default class SyncDestinationModel extends Model {
  @attr('string', { subText: 'Specifies the name for this destination.', editDisabled: true })
  name;

  @attr type;

  @attr('string', {
    subText:
      'Go-template string that indicates how to format the secret name at the destination. The default template varies by destination type but is generally in the form of "vault-<accessor_id>-<secret_path>" e.g. "vault-kv-1234-my-secret-1".',
  })
  secretNameTemplate;

  @attr('string', {
    editType: 'radio',
    label: 'Secret sync granularity',
    possibleValues: [
      {
        label: 'Secret path',
        subText: 'Sync entire secret contents as a single entry at the destination.',
        value: 'secret-path',
      },
      {
        label: 'Secret key',
        subText: 'Sync each key-value pair of secret data as a distinct entry at the destination.',
        helpText:
          'Only top-level keys will be synced and any nested or complex values will be encoded as a JSON string.',
        value: 'secret-key',
      },
    ],
  })
  granularity; // default value depends on type and is set in create route

  // only present if delete action has been initiated
  @attr('string') purgeInitiatedAt;
  @attr('string') purgeError;

  // findDestination returns static attributes for each destination type
  get icon() {
    return findDestination(this.type)?.icon;
  }

  get typeDisplayName() {
    return findDestination(this.type)?.name;
  }

  get maskedParams() {
    return findDestination(this.type)?.maskedParams;
  }

  @lazyCapabilities(apiPath`sys/sync/destinations/${'type'}/${'name'}`, 'type', 'name') destinationPath;
  @lazyCapabilities(apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/set`, 'type', 'name')
  setAssociationPath;

  get canCreate() {
    return this.destinationPath.get('canCreate') !== false;
  }
  get canDelete() {
    return this.destinationPath.get('canDelete') !== false;
  }
  get canEdit() {
    return this.destinationPath.get('canUpdate') !== false && !this.purgeInitiatedAt;
  }
  get canRead() {
    return this.destinationPath.get('canRead') !== false;
  }
  get canSync() {
    return this.setAssociationPath.get('canUpdate') !== false && !this.purgeInitiatedAt;
  }
}
