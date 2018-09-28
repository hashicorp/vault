import { create, attribute, visitable, collection, hasClass, text } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/access/'),
  navLinks: collection('[data-test-link]', {
    isActive: hasClass('is-active'),
    text: text(),
    scope: '[data-test-sidebar]',
  }),

  backendLinks: collection('[data-test-auth-backend-link]', {
    path: text('[data-test-path]'),
    id: attribute('data-test-id', '[data-test-path]'),
  }),

  findLinkById(id) {
    return this.backendLinks.filterBy('id', id)[0];
  },
});
