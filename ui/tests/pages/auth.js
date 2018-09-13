import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import withFlash from 'vault/tests/helpers/with-flash';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    await this.visit({ with: 'token' });
    if (token) {
      return this.tokenInput(token).submit();
    }

    return withFlash(this.tokenInput('root').submit());
  },
});
