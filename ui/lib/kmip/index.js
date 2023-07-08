/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */
/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
/* eslint-disable node/no-extraneous-require */
'use strict';

const EngineAddon = require('ember-engines/lib/engine-addon');

module.exports = EngineAddon.extend({
  name: 'kmip',

  lazyLoading: {
    enabled: true,
  },

  isDevelopingAddon() {
    return true;
  },
});
