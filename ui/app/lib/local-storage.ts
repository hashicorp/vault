/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default {
  isLocalStorageSupported() {
    try {
      const key = `__storage__test`;
      window.localStorage.setItem(key, '');
      window.localStorage.removeItem(key);
      return true;
    } catch (e) {
      // modify the e object so we can customize the error message.
      // e.message is readOnly.
      Object.defineProperty(e, 'errors', {
        value: [`This is likely due to your browser's cookie settings.`],
        writable: false,
      });

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
