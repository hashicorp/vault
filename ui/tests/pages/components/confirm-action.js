import { clickable } from 'ember-cli-page-object';

export default {
  delete: clickable('[data-test-confirm-action-trigger]', {
    testContainer: '#ember-testing',
  }),
  confirmDelete: clickable('[data-test-confirm-button]', {
    testContainer: '#ember-testing',
  }),
};
