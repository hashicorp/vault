/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable */

module.exports = {
  name: require('./package').name,

  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },

  isDevelopingAddon() {
    return true;
  },
};
