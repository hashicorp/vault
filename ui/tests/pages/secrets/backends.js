import { create, visitable, collection, text, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets'),
  links: collection({
    itemScope: '[data-test-secret-backend-link]',
    item: {
      path: text('[data-test-secret-path]'),
      toggleDetails: clickable('[data-test-secret-backend-detail]'),
      defaultTTL: text('[data-test-secret-backend-details="default-ttl"]'),
      maxTTL: text('[data-test-secret-backend-details="max-ttl"]'),
    },
    findByPath(path) {
      return this.toArray().findBy('path', path + '/');
    },
  }),
});
