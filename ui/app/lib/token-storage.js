/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
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
