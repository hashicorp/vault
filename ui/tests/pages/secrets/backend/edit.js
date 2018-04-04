import { create, visitable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets/:backend/edit/:id'),
  visitRoot: visitable('/vault/secrets/:backend/edit'),
});
