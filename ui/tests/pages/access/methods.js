import { create, attribute, visitable, collection, hasClass, text } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/access/'),
  navLinks: collection({
    scope: '[data-test-sidebar]',
    itemScope: '[data-test-link]',
    item: {
      isActive: hasClass('is-active'),
      text: text(),
    },
  }),

  backendLinks: collection({
    itemScope: '[data-test-auth-backend-link]',
    item: {
      path: text('[data-test-path]'),
      id: attribute('data-test-id', '[data-test-path]'),
    },
    findById(id) {
      return this.toArray().findBy('id', id);
    },
  }),
});
