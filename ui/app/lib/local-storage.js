/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default {
  isLocalStorageSupported() {
    try {
      const key = `__storage__test`;
      window.localStorage.setItem(key, null);
      window.localStorage.removeItem(key);
      return true;
    } catch (e) {
      // modify the e object so we can customize the error message.
      // e.message is readOnly.
      e.errors = [`This is likely due to your browser's cookie settings.`];
      throw e;
    }
  },

  getItem(key) {
    const item = window.localStorage.getItem(key);
    return item && JSON.parse(item);
  },

  setItem(key, val) {
    window.localStorage.setItem(key, JSON.stringify(val));
  },

  removeItem(key) {
    return window.localStorage.removeItem(key);
  },

  keys() {
    return Object.keys(window.localStorage);
  },

  cleanupStorage(string, keyToKeep) {
    if (!string) return;
    const relevantKeys = this.keys().filter((str) => str.startsWith(string));
    relevantKeys?.forEach((key) => {
      if (key !== keyToKeep) {
        localStorage.removeItem(key);
      }
    });
  },
};
