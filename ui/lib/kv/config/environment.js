/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */
'use strict';

module.exports = function (environment) {
  const ENV = {
    modulePrefix: 'kv',
    environment,
  };

  return ENV;
};
