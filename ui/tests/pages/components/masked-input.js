import { attribute, clickable, fillable, isPresent } from 'ember-cli-page-object';
import { focus, blur } from '@ember/test-helpers';

export default {
  wrapperClass: attribute('class', '[data-test-masked-input]'),
  enterText: fillable('[data-test-textarea]'),
  textareaIsPresent: isPresent('[data-test-textarea]'),
  copyButtonIsPresent: isPresent('[data-test-copy-button]'),
  toggleMasked: clickable('[data-test-button]'),
  async focusField() {
    return focus('[data-test-textarea]');
  },
  async blurField() {
    return blur('[data-test-textarea]');
  },
};
