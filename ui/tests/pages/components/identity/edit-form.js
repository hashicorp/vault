import { clickable, fillable, attribute } from 'ember-cli-page-object';
import fields from '../form-field';
import confirmAction from '../confirm-action';

export default {
  ...fields,
  ...confirmAction,
  cancelLinkHref: attribute('href', '[data-test-cancel-link]'),
  cancelLink: clickable('[data-test-cancel-link]'),
  name: fillable('[data-test-input="name"]'),
  disabled: clickable('[data-test-input="disabled"]'),
  type: fillable('[data-test-input="type"]'),
  submit: clickable('[data-test-identity-submit]'),
};
