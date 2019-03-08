import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  login: async function(token) {
    await this.visit({ with: 'token' });
    await settled();
    if (token) {
      await this.tokenInput(token).submit();
      await settled();
      return;
    }
    await this.tokenInput('root').submit();
    await settled();
  },
});
