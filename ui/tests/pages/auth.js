import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';

export default create({
  visit: visitable('/vault/auth'),
  logout: visitable('/vault/logout'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    await this.visit({ with: 'token' });
    await settled();
    if (token) {
      return await this.tokenInput(token).submit();
    }

    return await this.tokenInput('root').submit();
  },
});
