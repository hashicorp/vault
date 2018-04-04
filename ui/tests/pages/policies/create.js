import { create, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policies/create'),
});
