/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-disable node/no-extraneous-require */
const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: 'pki',
  lazyLoading: {
    enabled: false,
  },
  isDevelopingAddon() {
    return true;
  },
});
