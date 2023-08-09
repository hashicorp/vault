/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { set, get, computed } from '@ember/object';
import clamp from 'vault/utils/clamp';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const ACTION_VALUES = {
  encrypt: {
    isSupported: 'supportsEncryption',
    description: 'Looks up wrapping properties for the given token',
    glyph: 'lock-fill',
  },
  decrypt: {
    isSupported: 'supportsDecryption',
    description: 'Decrypts the provided ciphertext using this key',
    glyph: 'mail-open',
  },
  datakey: {
    isSupported: 'supportsEncryption',
    description: 'Generates a new key and value encrypted with this key',
    glyph: 'key',
  },
  rewrap: {
    isSupported: 'supportsEncryption',
    description: 'Rewraps the ciphertext using the latest version of the named key',
    glyph: 'reload',
  },
  sign: {
    isSupported: 'supportsSigning',
    description: 'Get the cryptographic signature of the given data',
    glyph: 'pencil-tool',
  },
  hmac: {
    isSupported: true,
    description: 'Generate a data digest using a hash algorithm',
    glyph: 'shuffle',
  },
  verify: {
    isSupported: true,
    description: 'Validate the provided signature for the given data',
    glyph: 'check-circle',
  },
  export: {
    isSupported: 'exportable',
    description: 'Get the named key',
    glyph: 'external-link',
  },
};

export default Model.extend({
  type: attr('string', {
    defaultValue: 'aes256-gcm96',
  }),
  name: attr('string', {
    label: 'Name',
    readOnly: true,
  }),
  autoRotatePeriod: attr({
    defaultValue: '0',
    defaultShown: 'Key is not automatically rotated',
    editType: 'ttl',
    label: 'Auto-rotation period',
  }),
  deletionAllowed: attr('boolean'),
  derived: attr('boolean'),
  exportable: attr('boolean'),
  minDecryptionVersion: attr('number', {
    defaultValue: 1,
  }),
  minEncryptionVersion: attr('number', {
    defaultValue: 0,
  }),
  latestVersion: attr('number'),
  keys: attr('object'),
  convergentEncryption: attr('boolean'),
  convergentEncryptionVersion: attr('number'),

  supportsSigning: attr('boolean'),
  supportsEncryption: attr('boolean'),
  supportsDecryption: attr('boolean'),
  supportsDerivation: attr('boolean'),

  setConvergentEncryption(val) {
    if (val === true) {
      set(this, 'derived', val);
    }
    set(this, 'convergentEncryption', val);
  },

  setDerived(val) {
    if (val === false) {
      set(this, 'convergentEncryption', val);
    }
    set(this, 'derived', val);
  },

  supportedActions: computed('type', function () {
    return Object.keys(ACTION_VALUES)
      .filter((name) => {
        const { isSupported } = ACTION_VALUES[name];
        return typeof isSupported === 'boolean' || get(this, isSupported);
      })
      .map((name) => {
        const { description, glyph } = ACTION_VALUES[name];
        return { name, description, glyph };
      });
  }),

  canDelete: computed('deletionAllowed', 'lastLoadTS', function () {
    const deleteAttrChanged = Boolean(this.changedAttributes().deletionAllowed);
    return this.deletionAllowed && deleteAttrChanged === false;
  }),

  keyVersions: computed('validKeyVersions', function () {
    let maxVersion = Math.max(...this.validKeyVersions);
    const versions = [];
    while (maxVersion > 0) {
      versions.unshift(maxVersion);
      maxVersion--;
    }
    return versions;
  }),

  encryptionKeyVersions: computed(
    'keyVerisons',
    'keyVersions',
    'latestVersion',
    'minDecryptionVersion',
    function () {
      const { keyVersions, minDecryptionVersion } = this;

      return keyVersions
        .filter((version) => {
          return version >= minDecryptionVersion;
        })
        .reverse();
    }
  ),

  keysForEncryption: computed('minEncryptionVersion', 'latestVersion', function () {
    let { minEncryptionVersion, latestVersion } = this;
    const minVersion = clamp(minEncryptionVersion - 1, 0, latestVersion);
    const versions = [];
    while (latestVersion > minVersion) {
      versions.push(latestVersion);
      latestVersion--;
    }
    return versions;
  }),

  validKeyVersions: computed('keys', function () {
    return Object.keys(this.keys);
  }),

  exportKeyTypes: computed('exportable', 'supportsEncryption', 'supportsSigning', 'type', function () {
    const types = ['hmac'];
    if (this.supportsSigning) {
      types.unshift('signing');
    }
    if (this.supportsEncryption) {
      types.unshift('encryption');
    }
    return types;
  }),

  backend: attr('string'),

  rotatePath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}/rotate`, 'backend', 'id'),
  canRotate: alias('rotatePath.canUpdate'),
  secretPath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id'),
  canRead: alias('secretPath.canUpdate'),
  canEdit: alias('secretPath.canUpdate'),
});
