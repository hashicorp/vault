'use strict';
const recommended = require('ember-template-lint/lib/config/recommended').rules; // octane extends recommended - no additions as of 3.14
const stylistic = require('ember-template-lint/lib/config/stylistic').rules;

const testOverrides = { ...recommended, ...stylistic };
for (const key in testOverrides) {
  testOverrides[key] = false;
}

module.exports = {
  extends: ['octane', 'stylistic'],
  rules: {
    'no-bare-strings': 'off',
    'no-action': 'off',
    'no-duplicate-landmark-elements': 'warn',
    'no-implicit-this': {
      allow: ['supported-auth-backends'],
    },
    'require-input-label': 'off',
    'no-down-event-binding': 'warn',
    'self-closing-void-elements': 'off',
  },
  ignore: ['lib/story-md', 'tests/**'],
  // ember language server vscode extension does not currently respect the ignore field
  // override all rules manually as workround to align with cli
  overrides: [
    {
      files: ['**/*-test.js'],
      rules: testOverrides,
    },
  ],
};
