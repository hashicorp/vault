import { clickable, create, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policy/:type/:name'),
  toggleEdit: clickable('[data-test-policy-edit-toggle]'),
});
