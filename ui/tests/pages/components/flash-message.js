import { collection } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

export default {
  latestMessage: getter(function() {
    return this.latestItem.text;
  }),
  latestItem: getter(function() {
    const count = this.messages().count;
    return this.messages(count - 1);
  }),
  messages: collection({
    itemScope: '[data-test-flash-message-body]',
  }),

  clickLast() {
    return this.latestItem.click();
  },
};
