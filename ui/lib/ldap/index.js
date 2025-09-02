/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable n/no-extraneous-require */
'use strict';

const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: 'ldap',

  lazyLoading: Object.freeze({
    enabled: false,
  }),

  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },

  isDevelopingAddon() {
    return true;
  },
});
