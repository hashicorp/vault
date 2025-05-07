/* eslint-env node */
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// defining config here rather than in package.json to run tsc on all .ts files, not just the staged changes
// this is accomplished by using function syntax rather than string

module.exports = {
  '*.{js,ts}': ['prettier --config .prettierrc.js --write', 'eslint --quiet', () => 'tsc --noEmit'],
  '*.hbs': ['prettier --config .prettierrc.js --write', 'ember-template-lint --quiet'],
  '*.scss': ['prettier --write'],
};
