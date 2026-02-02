/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable */

module.exports = {
  name: require('./package').name,

  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },

  options: {
    'ember-cli-babel': {
      enableTypeScriptTransform: true,
    },
  },

  isDevelopingAddon() {
    return true;
  },
};
