import { attribute, clickable, fillable, focusable, blurrable } from 'ember-cli-page-object';

export default {
  wrapperClass: attribute('class','[data-test-masked-input]'),
  enterText: fillable('[data-test-textarea]'),
  toggleMasked: clickable('[data-test-button]'),
  focus: focusable('[data-test-textarea]'),
  blur: blurrable('[data-test-textarea]')
};
