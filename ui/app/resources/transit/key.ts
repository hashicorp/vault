/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';
import clamp from 'vault/utils/clamp';

export interface TransitKeyData {
  backend: string;
  id: string;
  name?: string;
  type?: string;
  keys?: Record<string, number | { creation_time: string; public_key?: string }>;
  latest_version?: number;
  min_decryption_version?: number;
  min_encryption_version?: number;
  auto_rotate_period?: number | string;
  deletion_allowed?: boolean;
  derived?: boolean;
  convergent_encryption?: boolean;
  exportable?: boolean;
  supports_encryption?: boolean;
  supports_decryption?: boolean;
  supports_signing?: boolean;
}

const ACTION_VALUES = {
  encrypt: {
    isSupported: 'supports_encryption',
    description: 'Looks up wrapping properties for the given token.',
    glyph: 'lock-fill',
  },
  decrypt: {
    isSupported: 'supports_decryption',
    description: 'Decrypts the provided ciphertext using this key.',
    glyph: 'mail-open',
  },
  datakey: {
    isSupported: 'supports_encryption',
    description: 'Generates a new key and value encrypted with this key.',
    glyph: 'key',
  },
  rewrap: {
    isSupported: 'supports_encryption',
    description: 'Rewraps the ciphertext using the latest version of the named key.',
    glyph: 'reload',
  },
  sign: {
    isSupported: 'supports_signing',
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
} as const;

/** Read-only Transit key data and the action/version choices derived from the API response. */
export default class TransitKeyResource extends baseResourceFactory<TransitKeyData>() {
  get validKeyVersions() {
    return Object.keys(this.keys || {});
  }

  get keyVersions() {
    const maxVersion = Math.max(0, ...this.validKeyVersions.map(Number));
    return Array.from({ length: maxVersion }, (_, index) => index + 1);
  }

  get encryptionKeyVersions() {
    return this.keyVersions.filter((version) => version >= (this.min_decryption_version || 0)).reverse();
  }

  get keysForEncryption() {
    const latestVersion = this.latest_version || 0;
    const minimum = clamp((this.min_encryption_version || 0) - 1, 0, latestVersion);
    return Array.from({ length: latestVersion - minimum }, (_, index) => latestVersion - index);
  }

  get exportKeyTypes() {
    return [
      ...(this.supports_encryption ? ['encryption'] : []),
      ...(this.supports_signing ? ['signing'] : []),
      'hmac',
    ];
  }

  get supportedActions() {
    return Object.entries(ACTION_VALUES)
      .filter(([, { isSupported }]) => isSupported === true || this[isSupported])
      .map(([name, { description, glyph }]) => ({ name, description, glyph }));
  }
}
