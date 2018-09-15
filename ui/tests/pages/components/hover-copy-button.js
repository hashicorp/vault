import { attribute, clickable, isVisible, focusable, text } from 'ember-cli-page-object';
import { triggerEvent, focus } from '@ember/test-helpers';

export default {
  async focusContainer() {
    await focus('.has-copy-button');
  },
  tooltipText: text('[data-test-hover-copy-tooltip-text]', {
    testContainer: '#ember-testing',
  }),
  wrapperClass: attribute('class', '[data-test-hover-copy]'),
  buttonIsVisible: isVisible('[data-test-hover-copy-button]'),
  click: clickable('[data-test-hover-copy-button]'),
  focus: focusable('[data-test-hover-copy-button]'),

  async mouseEnter() {
    await triggerEvent('[data-test-tooltip-trigger]', 'mouseenter');
  },
};
