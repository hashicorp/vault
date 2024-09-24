/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { service } from '@ember/service';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  type: [{ type: 'presence', message: 'Type is required.' }],
  keyType: [{ type: 'presence', message: 'Please select a key type.' }],
  keyName: [
    {
      validator(model) {
        if (model.keyName === 'default') return false;
        return true;
      },
      message: `Key name cannot be the reserved value 'default'`,
    },
  ],
};
const displayFields = ['keyId', 'keyName', 'keyType', 'keyBits'];
const formFieldGroups = [{ default: ['keyName', 'type'] }, { 'Key parameters': ['keyType', 'keyBits'] }];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class PkiKeyModel extends Model {
  @service secretMountPath;

  @attr('string', { detailsLabel: 'Key ID' }) keyId;
  @attr('string', {
    subText: `Optional, human-readable name for this key. The name must be unique across all keys and cannot be 'default'.`,
  })
  keyName;
  @attr('string', {
    noDefault: true,
    possibleValues: ['internal', 'exported'],
    subText:
      'The type of operation. If exported, the private key will be returned in the response; if internal the private key will not be returned and cannot be retrieved later.',
  })
  type;
  @attr('string', {
    noDefault: true,
    possibleValues: ['rsa', 'ec', 'ed25519'],
    subText: 'The type of key that will be generated. Must be rsa, ed25519, or ec. ',
  })
  keyType;
  @attr('string', {
    label: 'Key bits',
    noDefault: true,
    subText: 'Bit length of the key to generate.',
  })
  keyBits; // no possibleValues because dependent on selected key type

  @attr('string') pemBundle;
  @attr('string') privateKey;

  get backend() {
    return this.secretMountPath.currentPath;
  }

  /* CAPABILITIES
   * Default to show UI elements unless we know they can't access the given path
   */

  @lazyCapabilities(apiPath`${'backend'}/key/${'keyId'}`, 'backend', 'keyId') keyPath;
  get canRead() {
    return this.keyPath.get('canRead') !== false;
  }
  get canEdit() {
    return this.keyPath.get('canUpdate') !== false;
  }
  get canDelete() {
    return this.keyPath.get('canDelete') !== false;
  }

  @lazyCapabilities(apiPath`${'backend'}/keys/generate`, 'backend') generatePath;
  get canGenerateKey() {
    return this.generatePath.get('canUpdate') !== false;
  }

  @lazyCapabilities(apiPath`${'backend'}/keys/import`, 'backend') importPath;
  get canImportKey() {
    return this.importPath.get('canUpdate') !== false;
  }
}
