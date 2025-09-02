/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
/* eslint-disable n/no-extraneous-require */
'use strict';

const EngineAddon = require('ember-engines/lib/engine-addon');

module.exports = EngineAddon.extend({
  name: 'replication',

  lazyLoading: {
    enabled: true,
  },

  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },

  isDevelopingAddon() {
    return true;
  },
});
