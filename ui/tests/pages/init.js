import { text, create, collection, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/init'),
  submit: clickable('[data-test-init-submit]'),
  shares: fillable('[data-test-key-shares]'),
  threshold: fillable('[data-test-key-threshold]'),
  keys: collection('[data-test-key-box]'),
  buttonText: text('[data-test-advance-button]'),
  init: async function(shares, threshold) {
    await this.visit();
    return this.shares(shares)
      .threshold(threshold)
      .submit();
  },
});
