import { create, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    await this.visit({ with: 'token' });
    if (token) {
      return this.tokenInput(token).submit();
    }

    return this.tokenInput('root').submit();
  },
});
