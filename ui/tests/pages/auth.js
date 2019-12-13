import { create, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/auth'),
  logout: visitable('/vault/logout'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await this.visit({ with: 'token' });
    if (token) {
      return this.tokenInput(token).submit();
    }

    return this.tokenInput('root').submit();
  },
});
