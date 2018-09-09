import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import flashMessages from '../../components/flash-message';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  flash: flashMessages,
  login: async function(token) {
    await this.visit({ with: 'token' })
      .tokenInput(token || 'root')
      .submit();
  },
});
