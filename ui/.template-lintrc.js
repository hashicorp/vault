/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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
    // from bump to ember-template-lint@6.0.0
    'no-builtin-form-components': 'off',
    'no-at-ember-render-modifiers': 'off',
    'no-unnecessary-curly-strings': 'off',
    'no-unnecessary-curly-parens': 'off',
    'no-action-on-submit-button': 'off',
  },
  overrides: [
    {
      files: ['**/*-test.js'],
      rules: {
        prettier: false,
      },
    },
  ],
};
