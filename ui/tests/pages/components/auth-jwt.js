import { text, isPresent, clickable, fillable } from 'ember-cli-page-object';

export default {
  jwt: fillable('[data-test-jwt]'),
  jwtPresent: isPresent('[data-test-jwt]'),
  role: fillable('[data-test-role]'),
  rolePresent: isPresent('[data-test-role]'),
  login: clickable('[data-test-auth-submit]'),
  loginButtonText: text('[data-test-auth-submit]'),
  yieldContent: text('[data-test-yield-content]'),
};
