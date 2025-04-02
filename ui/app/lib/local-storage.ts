/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ApiError } from 'vault/api';

export default {
  isLocalStorageSupported() {
    try {
      const key = `__storage__test`;
      window.localStorage.setItem(key, '');
      window.localStorage.removeItem(key);
      return true;
    } catch (e) {
      const error = e as ApiError;
      // modify the e object so we can customize the error message.
      // e.message is readOnly.
      error.errors = [`This is likely due to your browser's cookie settings.`];
      throw e;
    }
  },

  getItem(key: string) {
    const item = window.localStorage.getItem(key);
    return item && JSON.parse(item);
  },

  setItem(key: string, val: unknown) {
    window.localStorage.setItem(key, JSON.stringify(val));
  },

  removeItem(key: string) {
    return window.localStorage.removeItem(key);
  },

  keys() {
    return Object.keys(window.localStorage);
  },

  cleanupStorage(string: string, keyToKeep: string) {
    if (!string) return;
    const relevantKeys = this.keys().filter((str) => str.startsWith(string));
    relevantKeys?.forEach((key) => {
      if (key !== keyToKeep) {
        localStorage.removeItem(key);
      }
    });
  },
};
