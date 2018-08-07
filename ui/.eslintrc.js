module.exports = {
  root: true,
  parserOptions: {
    ecmaVersion: 2017,
    sourceType: 'module',
    ecmaFeatures: {
      experimentalObjectRestSpread: true,
    },
  },
  extends: 'eslint:recommended',
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
};
