import { clickable, fillable, attribute } from 'ember-cli-page-object';
import fields from '../form-field';

export default {
  ...fields,

  cancelLinkHref: attribute('href', '[data-test-cancel-link]'),
  name: fillable('[data-test-input="name"]'),
  disabled: clickable('[data-test-input="disabled"]'),
  submit: clickable('[data-test-identity-submit]'),
};
