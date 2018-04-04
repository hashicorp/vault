import { collection } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

export default {
  latestMessage: getter(function() {
    const count = this.messages().count;
    return this.messages(count - 1).text;
  }),
  messages: collection({
    itemScope: '[data-test-flash-message-body]',
  }),
};
