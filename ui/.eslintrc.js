/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable no-undef */

'use strict';

module.exports = {
  parser: '@babel/eslint-parser',
  root: true,
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    requireConfigFile: false,
    babelOptions: {
      plugins: [['@babel/plugin-proposal-decorators', { decoratorsBeforeExport: true }]],
    },
  },
  plugins: ['ember'],
  extends: [
    'eslint:recommended',
    'plugin:ember/recommended',
    'plugin:prettier/recommended',
    'plugin:compat/recommended',
  ],
  env: {
    browser: true,
  },
  rules: {
    'no-console': 'error',
    'prefer-const': ['error', { destructuring: 'all' }],
    'ember/no-mixins': 'warn',
    'ember/no-new-mixins': 'off', // should be warn but then every line of the mixin is green
    // need to be fully glimmerized before these rules can be turned on
    'ember/no-classic-classes': 'off',
    'ember/no-classic-components': 'off',
    'ember/no-actions-hash': 'off',
    'ember/require-tagless-components': 'off',
    'ember/no-component-lifecycle-hooks': 'off',
  },
  overrides: [
    // node files
    {
      files: [
        './.eslintrc.js',
        './.prettierrc.js',
        './.stylelintrc.js',
        './.template-lintrc.js',
        './ember-cli-build.js',
        './testem.js',
        './blueprints/*/index.js',
        './config/**/*.js',
        './lib/*/index.js',
        './server/**/*.js',
      ],
      parserOptions: {
        sourceType: 'script',
      },
      env: {
        browser: false,
        node: true,
      },
      extends: ['plugin:n/recommended'],
    },
    {
      // test files
      files: ['tests/**/*-test.{js,ts}'],
      extends: ['plugin:qunit/recommended'],
      rules: {
        'qunit/require-expect': 'off',
      },
    },
    {
      files: ['**/*.ts'],
      extends: ['plugin:@typescript-eslint/recommended'],
    },
  ],
};
