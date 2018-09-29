import { text, create, collection, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policies/:type'),
  policies: collection('[data-test-policy-item]', {
    name: text('[data-test-policy-name]'),
  }),
  findPolicyByName(name) {
    return this.policies.filterBy('name', name)[0];
  },
});
