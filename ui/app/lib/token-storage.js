/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import localStorageWrapper from './local-storage';
import memoryStorage from './memory-storage';

export default function (type) {
  if (type === 'memory') {
    return memoryStorage;
  }
  let storage;
  try {
    window.localStorage.getItem('test');
    storage = localStorageWrapper;
  } catch (e) {
    storage = memoryStorage;
  }
  return storage;
}
