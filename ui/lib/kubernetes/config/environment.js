/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
'use strict';

module.exports = function (environment) {
  const ENV = {
    modulePrefix: 'kubernetes',
    environment,
  };

  return ENV;
};
