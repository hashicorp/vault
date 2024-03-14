/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { set, get } from '@ember/object';
import clamp from 'vault/utils/clamp';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const ACTION_VALUES = {
  encrypt: {
    isSupported: 'supportsEncryption',
    description: 'Looks up wrapping properties for the given token.',
    glyph: 'lock-fill',
  },
  decrypt: {
    isSupported: 'supportsDecryption',
    description: 'Decrypts the provided ciphertext using this key.',
    glyph: 'mail-open',
  },
  datakey: {
    isSupported: 'supportsEncryption',
    description: 'Generates a new key and value encrypted with this key.',
    glyph: 'key',
  },
  rewrap: {
    isSupported: 'supportsEncryption',
    description: 'Rewraps the ciphertext using the latest version of the named key.',
    glyph: 'reload',
  },
  sign: {
    isSupported: 'supportsSigning',
    description: 'Get the cryptographic signature of the given data.',
    glyph: 'pencil-tool',
  },
  hmac: {
    isSupported: true,
    description: 'Generate a data digest using a hash algorithm.',
    glyph: 'shuffle',
  },
  verify: {
    isSupported: true,
    description: 'Validate the provided signature for the given data.',
    glyph: 'check-circle',
  },
  export: {
    isSupported: 'exportable',
    description: 'Get the named key.',
    glyph: 'external-link',
  },
};

export default class TransitKeyModel extends Model {
  @attr('string') backend;
  @attr('string', {
    defaultValue: 'aes256-gcm96',
  })
  type;

  @attr('string', {
    label: 'Name',
    readOnly: true,
  })
  name;

  @attr({
    defaultValue: '0',
    defaultShown: 'Key is not automatically rotated',
    editType: 'ttl',
    label: 'Auto-rotation period',
  })
  autoRotatePeriod;

  @attr('boolean') deletionAllowed;
  @attr('boolean') derived;
  @attr('boolean') exportable;

  @attr('number', {
    defaultValue: 1,
  })
  minDecryptionVersion;

  @attr('number', {
    defaultValue: 0,
  })
  minEncryptionVersion;

  @attr('number') latestVersion;
  @attr('object') keys;
  @attr('boolean') convergentEncryption;
  @attr('number') convergentEncryptionVersion;

  @attr('boolean') supportsSigning;
  @attr('boolean') supportsEncryption;
  @attr('boolean') supportsDecryption;
  @attr('boolean') supportsDerivation;

  setConvergentEncryption(val) {
    if (val === true) {
      set(this, 'derived', val);
    }
    set(this, 'convergentEncryption', val);
  }

  setDerived(val) {
    if (val === false) {
      set(this, 'convergentEncryption', val);
    }
    set(this, 'derived', val);
  }

  get supportedActions() {
    return Object.keys(ACTION_VALUES)
      .filter((name) => {
        const { isSupported } = ACTION_VALUES[name];
        return typeof isSupported === 'boolean' || get(this, isSupported);
      })
      .map((name) => {
        const { description, glyph } = ACTION_VALUES[name];
        return { name, description, glyph };
      });
  }

  get canDelete() {
    const deleteAttrChanged = Boolean(this.changedAttributes().deletionAllowed);
    return this.deletionAllowed && deleteAttrChanged === false;
  }

  get keyVersions() {
    let maxVersion = Math.max(...this.validKeyVersions);
    const versions = [];
    while (maxVersion > 0) {
      versions.unshift(maxVersion);
      maxVersion--;
    }
    return versions;
  }

  get encryptionKeyVersions() {
    const { keyVersions, minDecryptionVersion } = this;

    return keyVersions
      .filter((version) => {
        return version >= minDecryptionVersion;
      })
      .reverse();
  }

  get keysForEncryption() {
    let { minEncryptionVersion, latestVersion } = this;
    const minVersion = clamp(minEncryptionVersion - 1, 0, latestVersion);
    const versions = [];
    while (latestVersion > minVersion) {
      versions.push(latestVersion);
      latestVersion--;
    }
    return versions;
  }

  get validKeyVersions() {
    return Object.keys(this.keys);
  }

  get exportKeyTypes() {
    const types = ['hmac'];
    if (this.supportsSigning) {
      types.unshift('signing');
    }
    if (this.supportsEncryption) {
      types.unshift('encryption');
    }
    return types;
  }
  @lazyCapabilities(apiPath`${'backend'}/keys/${'id'}/rotate`, 'backend', 'id') rotatePath;
  @lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id') secretPath;

  get canRotate() {
    return this.rotatePath.get('canUpdate') !== false;
  }
  get canRead() {
    return this.secretPath.get('canUpdate') !== false;
  }
  get canEdit() {
    return this.secretPath.get('canUpdate') !== false;
  }
}
