import { text, isPresent } from 'ember-cli-page-object';

export default {
  hasLabel: isPresent('[data-test-row-label]'),
  rowLabel: text('[data-test-row-label]'),
  rowValue: text('[data-test-row-value]'),
};
