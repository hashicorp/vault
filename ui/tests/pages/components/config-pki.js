import { clickable, fillable, text, isPresent } from 'ember-cli-page-object';
import fields from './form-field';

export default {
  ...fields,
  scope: '.config-pki',
  text: text('[data-test-text]'),
  title: text('[data-test-title]'),
  hasTitle: isPresent('[data-test-title]'),
  hasError: isPresent('[data-test-error]'),
  submit: clickable('[data-test-submit]'),
  fillInField: fillable('[data-test-field]'),
};
