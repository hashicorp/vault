import { create, visitable, collection, clickable, text } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets'),
  rows: collection({
    itemScope: '[data-test-secret-backend-row]',
    item: {
      path: text('[data-test-secret-path]'),
      menu: clickable('[data-test-popup-menu-trigger]'),
    },
    findByPath(path) {
      return this.toArray().findBy('path', path + '/');
    },
  }),
  configLink: clickable('[data-test-engine-config]'),
});
