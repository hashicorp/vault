module.exports = {
  root: true,
  parserOptions: {
    ecmaVersion: 2017,
    sourceType: 'module',
    ecmaFeatures: {
      experimentalObjectRestSpread: true,
    },
  },
  plugins: ['ember'],
  extends: ['eslint:recommended', 'plugin:ember/recommended'],
  env: {
    browser: true,
    es6: true,
  },
  rules: {
    'no-unused-vars': ['error', { ignoreRestSiblings: true }],
  },
  globals: {
    TextEncoderLite: true,
    TextDecoderLite: true,
  },
  overrides: [
    // node files
    {
      files: [
        '.template-lintrc.js',
        'ember-cli-build.js',
        'testem.js',
        'blueprints/*/index.js',
        'config/**/*.js',
        'lib/*/index.js',
        'scripts/start-vault.js',
      ],
      parserOptions: {
        sourceType: 'script',
        ecmaVersion: 2015,
      },
      env: {
        browser: false,
        node: true,
      },
    },
  ],
};
