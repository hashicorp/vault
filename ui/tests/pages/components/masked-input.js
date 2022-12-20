import { clickable, isPresent } from 'ember-cli-page-object';

export default {
  textareaIsPresent: isPresent('[data-test-textarea]'),
  copyButtonIsPresent: isPresent('[data-test-copy-button]'),
  toggleMasked: clickable('[data-test-button="toggle-masked"]'),
};
