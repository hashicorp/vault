/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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

  isDevelopingAddon() {
    return true;
  },
});
