import { clickable, fillable, text, value } from 'ember-cli-page-object';
import fields from './form-field';
import errorText from './message-in-page';

export default {
  ...fields,
  ...errorText,
  header: text('[data-test-mount-form-header]'),
  submit: clickable('[data-test-mount-submit]'),
  path: fillable('[data-test-input="path"]'),
  pathValue: value('[data-test-input="path"]'),
  type: fillable('[data-test-input="type"]'),
  typeValue: value('[data-test-input="type"]'),
};
