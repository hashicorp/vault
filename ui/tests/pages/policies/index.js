import { text, create, collection, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policies/:type'),
  policies: collection({
    itemScope: '[data-test-policy-item]',
    item: {
      name: text('[data-test-policy-name]'),
    },
    findByName(name) {
      return this.toArray().findBy('name', name);
    },
  }),
});
