/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

'use strict';

module.exports = {
  singleQuote: true,
  trailingComma: 'es5',
  printWidth: 110,
  overrides: [
    {
      files: '*.hbs',
      options: {
        singleQuote: false,
        printWidth: 125,
      },
    },
    {
      files: '*.{js,ts}',
      options: {
        singleQuote: true,
      },
    },
  ],
};
