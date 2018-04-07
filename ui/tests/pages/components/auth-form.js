import { clickable, text } from 'ember-cli-page-object';

export default {
  errorText: text('[data-test-auth-error]'),
  login: clickable('[data-test-auth-submit]'),
};
