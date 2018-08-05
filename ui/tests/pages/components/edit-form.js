import { clickable, isPresent, text } from 'ember-cli-page-object';
import fields from './form-field';
export default {
  ...fields,
  submit: clickable('[data-test-edit-form-submit]'),
  deleteButton: clickable('[data-test-confirm-action-trigger]'),
  deleteConfirm: clickable('[data-test-confirm-button]'),
  deleteText: text('[data-test-edit-delete-text]'),
  showsDelete: isPresent('[data-test-edit-delete-text]'),
  errorText: text('[data-test-edit-form-error]'),
};
