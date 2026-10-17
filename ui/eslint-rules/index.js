/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

'use strict';

/**
 * eslint-plugin-vault
 *
 * Local ESLint plugin for Vault UI custom rules.
 * Rules live in ./rules/ — add new rules there and register them here.
 */

module.exports = {
  rules: {
    'require-version-guard': require('./require-version-guard'),
  },
};
