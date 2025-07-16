/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable n/no-extraneous-require */
const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: 'config-ui',
  lazyLoading: {
    enabled: false,
  },
  isDevelopingAddon() {
    return true;
  },
});
