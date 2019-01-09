import { text, create, collection, clickable, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policies/:type'),
  policies: collection('[data-test-policy-item]', {
    name: text('[data-test-policy-name]'),
    menu: clickable('[data-test-popup-menu-trigger]'),
  }),
  findPolicyByName(name) {
    return this.policies.filterBy('name', name)[0];
  },
  delete: clickable('[data-test-confirm-action-trigger]', {
    scope: '[data-test-item-delete]',
    testContainer: '#ember-testing',
  }),
  confirmDelete: clickable('[data-test-confirm-button]', {
    scope: '[data-test-item-delete]',
    testContainer: '#ember-testing',
  }),
});
