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
  name: 'kmip',

  lazyLoading: {
    enabled: true,
  },

  included: function (/* app */) {
    this._super.included.apply(this, arguments);

    // this is disabled in test because ember-asset-manifest doesn't work in embroider
    // to make sure engine code is loaded in the host app, we have to
    // disable lazyLoading
    // see https://github.com/embroider-build/embroider/issues/996
    this.options.lazyLoading.enabled = process.env.EMBER_ENV !== 'test';
  },

  isDevelopingAddon() {
    return true;
  },
});
