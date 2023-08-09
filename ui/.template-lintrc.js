/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

'use strict';

const fs = require('fs');
let testOverrides = {};
try {
  // ember-template-lint no longer exports anything so we cannot access the rule definitions conventionally
  // read file, convert to json string and parse
  const toJSON = (str) => {
    return JSON.parse(
      str
        .slice(str.indexOf(':') + 2) // get rid of export statement
        .slice(0, -(str.length - str.lastIndexOf(','))) // remove trailing brackets from export
        .replace(/:.*,/g, `: ${false},`) // convert values to false
        .replace(/,([^,]*)$/, '$1') // remove last comma
        .replace(/'/g, '"') // convert to double quotes
        .replace(/(\w[^"].*[^"]):/g, '"$1":') // wrap quotes around single word keys
        .trim()
    );
  };
  const recommended = toJSON(
    fs.readFileSync('node_modules/ember-template-lint/lib/config/recommended.js').toString()
  );
  const stylistic = toJSON(
    fs.readFileSync('node_modules/ember-template-lint/lib/config/stylistic.js').toString()
  );
  testOverrides = {
    ...recommended,
    ...stylistic,
    prettier: false,
  };
} catch (error) {
  console.log(error); // eslint-disable-line
}

module.exports = {
  plugins: ['ember-template-lint-plugin-prettier'],
  extends: ['recommended', 'ember-template-lint-plugin-prettier:recommended'],
  rules: {
    'no-action': 'off',
    'no-implicit-this': {
      allow: ['supported-auth-backends'],
    },
    'require-input-label': 'off',
    'no-array-prototype-extensions': 'off',
  },
  ignore: ['lib/story-md', 'tests/**'],
  // ember language server vscode extension does not currently respect the ignore field
  // override all rules manually as workaround to align with cli
  overrides: [
    {
      files: ['**/*-test.js'],
      rules: testOverrides,
    },
  ],
};
