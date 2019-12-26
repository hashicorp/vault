import { isPresent, collection, text, clickable } from 'ember-cli-page-object';

export default {
  hasSearchSelect: isPresent('[data-test-component=search-select]'),
  hasTrigger: isPresent('.ember-power-select-trigger'),
  hasLabel: isPresent('[data-test-field-label]'),
  labelText: text('[data-test-field-label]'),
  options: collection('.ember-power-select-option'),
  selectedOptions: collection('[data-test-selected-option]'),
  deleteButtons: collection('[data-test-selected-list-button="delete"]'),
  selectedOptionText: text('[aria-current=true]'),
  selectOption: clickable('[aria-current=true]'),
  hasStringList: isPresent('[data-test-string-list-input]'),
  smallOptionIds: collection('[data-test-smaller-id]'),
};
