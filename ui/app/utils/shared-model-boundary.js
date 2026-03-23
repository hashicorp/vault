/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function clearModelCache(store, modelType) {
  if (!modelType) {
    return;
  }
  const modelTypes = Array.isArray(modelType) ? modelType : [modelType];
  if (store.isDestroyed || store.isDestroying) {
    // Prevent unload attempt after test teardown, resulting in test failure
    return;
  }
  modelTypes.forEach((type) => {
    store.unloadAll(type);
  });
}
