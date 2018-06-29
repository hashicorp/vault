import { attribute, isVisible, triggerable, focusable, text } from 'ember-cli-page-object';

export default {
  focusContainer: focusable('.has-copy-button'),
  mouseEnter: triggerable('mouseenter', '[data-test-tooltip-trigger]'),
  tooltipText: text('[data-test-hover-copy-tooltip-text]'),
  wrapperClass: attribute('class', '[data-test-hover-copy]'),
  buttonIsVisible: isVisible('[data-test-hover-copy-button]'),
};
