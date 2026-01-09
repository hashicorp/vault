/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable n/no-extraneous-require */
const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: 'kubernetes',
  lazyLoading: {
    enabled: false,
  },
  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },
  'ember-cli-babel': {
    enableTypeScriptTransform: true,
  },
  isDevelopingAddon() {
    return true;
  },
});
