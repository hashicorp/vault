import { collection, isPresent, property, clickable } from 'ember-cli-page-object';

export default {
  wizardItems: collection('[data-test-select-input]', {
    hasDisabledTooltip: isPresent('[data-test-tooltip]'),
  }),
  hasDisabledStartButton: property('disabled', '[data-test-start-button]'),
  selectSecrets: clickable('[data-test-checkbox=Secrets]'),
};
