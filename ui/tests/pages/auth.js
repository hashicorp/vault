import { create, visitable, fillable, clickable, collection } from 'ember-cli-page-object';
import { waitFor } from 'ember-test-helpers';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  flashes: collection({
    itemScope: '[data-test-flash-message-body]',
    click: clickable(),
  }),
  clickLast() {
    const count = this.flashes().count;
    return this.flashes(count - 1).click();
  },
  login: async function(token) {
    this.visit({ with: 'token' })
      .tokenInput(token || 'root')
      .submit();
    await waitFor('[data-test-flash-message-body]');
    this.clickLast();
  },
});
