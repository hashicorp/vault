/* eslint-disable no-undef */

'use strict';

module.exports = {
  parser: 'babel-eslint',
  root: true,
  parserOptions: {
    ecmaVersion: 2018,
    sourceType: 'module',
    ecmaFeatures: {
      legacyDecorators: true,
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
      plugins: ['node'],
      extends: ['plugin:node/recommended'],
      rules: {
        // this can be removed once the following is fixed
        // https://github.com/mysticatea/eslint-plugin-node/issues/77
        'node/no-unpublished-require': 'off',
      },
    },
    {
      // test files
      files: ['tests/**/*-test.{js,ts}'],
      extends: ['plugin:qunit/recommended'],
    },
    {
      files: ['**/*.ts'],
      extends: ['plugin:@typescript-eslint/recommended'],
    },
  ],
};
