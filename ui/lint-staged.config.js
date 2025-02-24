/* eslint-env node */

module.exports = {
  '*.{js,ts}': ['prettier --config .prettierrc.js --write', 'eslint --quiet', () => 'tsc --noEmit'],
  '*.hbs': ['prettier --config .prettierrc.js --write', 'ember-template-lint --quiet'],
  '*.scss': ['prettier --write'],
};
