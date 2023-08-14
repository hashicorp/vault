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
  ignore: ['tests/**'],
};
