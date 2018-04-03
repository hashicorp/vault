import { create, visitable } from 'ember-cli-page-object';

export const Base = {
  visit: visitable('/vault/secrets/:backend/credentials/:id'),
  visitRoot: visitable('/vault/secrets/:backend/credentials'),
};

export default create(Base);
