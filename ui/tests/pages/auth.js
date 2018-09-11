import { create, visitable, fillable, clickable, collection } from 'ember-cli-page-object';
import { waitFor } from 'ember-test-helpers';

export default create({
  visit: visitable('/vault/auth'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  flashes: collection('[data-test-flash-message-body]', {
    click: clickable(),
  }),
  clickLast: async function() {
    const count = this.flashes.length;
    let last = this.flashes.objectAt(count - 1);
    return count && last.click();
  },
  login: async function(token) {
    this.visit({ with: 'token' })
      .tokenInput(token || 'root')
      .submit();
    await waitFor('[data-test-flash-message-body]');
    this.clickLast();
  },
});
