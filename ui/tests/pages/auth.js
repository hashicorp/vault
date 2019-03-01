import { create, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    await this.visit({ with: 'token' });
    if (token) {
      return await this.tokenInput(token).submit();
    }

    return await this.tokenInput('root').submit();
  },
});
