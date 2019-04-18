import { collection, clickable, fillable, text, value } from 'ember-cli-page-object';

export default {
  tabs: collection('[data-test-auth-method]', {
    name: text(),
    link: clickable('[data-test-auth-method-link]'),
  }),
  selectMethod: fillable('[data-test-method-select]'),
  username: fillable('[data-test-username]'),
  token: fillable('[data-test-token]'),
  tokenValue: value('[data-test-token]'),
  password: fillable('[data-test-password]'),
  errorText: text('[data-test-auth-error]'),
  login: clickable('[data-test-auth-submit]'),
};
