/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

const cache: { [key: string]: string } = {};

export default {
  getItem(key: string) {
    const item = cache[key];
    return item && JSON.parse(item);
  },

  setItem(key: string, val: unknown) {
    cache[key] = JSON.stringify(val);
  },

  removeItem(key: string) {
    delete cache[key];
  },

  keys() {
    return Object.keys(cache);
  },
};
