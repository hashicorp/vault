import { text, isPresent } from 'ember-cli-page-object';
//  ARG TODO: add test that shows table value concatenated when over 10 in length.
export default {
  hasLabel: isPresent('[data-test-row-label]'),
  rowLabel: text('[data-test-row-label]'),
  rowValue: text('[data-test-row-value]'),
};
