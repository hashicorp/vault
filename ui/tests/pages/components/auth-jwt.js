import { clickable, fillable } from 'ember-cli-page-object';

export default {
  jwt: fillable('[data-test-jwt]'),
  role: fillable('[data-test-role]'),
  login: clickable('[data-test-auth-submit]'),
};
