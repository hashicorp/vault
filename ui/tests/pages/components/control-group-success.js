import { isPresent, fillable, clickable } from 'ember-cli-page-object';

export default {
  showsJsonViewer: isPresent('[data-test-json-viewer]'),
  showsNavigateMessage: isPresent('[data-test-navigate-message]'),
  showsUnwrapForm: isPresent('[data-test-unwrap-form]'),
  navigate: clickable('[data-test-navigate-button]'),
  unwrap: clickable('[data-test-unwrap-button]'),
  token: fillable('[data-test-token-input]'),
};
