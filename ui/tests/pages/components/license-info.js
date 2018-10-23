import { clickable, fillable, text, isPresent, collection } from 'ember-cli-page-object';

export default {
  text: fillable('[data-test-text-input]'),
  isTemp: isPresent('[data-test-temp-license]'),
  hasTextInput: isPresent('[data-test-text-input]'),
  saveButton: clickable('[data-test-save-button]'),
  hasSaveButton: isPresent('[data-test-save-button]'),
  enterButton: clickable('[data-test-enter-button]'),
  hasEnterButton: isPresent('[data-test-enter-button]'),
  cancelButton: clickable('[data-test-cancel-button]'),
  hasWarning: isPresent('[data-test-warning-text]'),
  warning: text('[data-test-warning-text]'),
  featureRows: collection('[data-test-feature-row]', {
    featureName: text('[data-test-row-label]'),
    featureStatus: text('[data-test-feature-status]'),
  }),
};
