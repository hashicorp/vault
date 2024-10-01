/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

module.exports = {
  ...require('@hashicorp/platform-cli/config/.eslintrc'),
  /* Specify overrides here */
  ignorePatterns: ['public/']
}
