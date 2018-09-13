import { waitFor, settled } from '@ember/test-helpers';
import { collection, text, clickable } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

export default {
  latestMessage: getter(function() {
    return this.latestItem.text;
  }),
  latestItem: getter(function() {
    const count = this.messages.length;
    return this.messages.objectAt(count - 1);
  }),
  messages: collection('[data-test-flash-message-body]', {
    click: clickable(),
    text: text(),
  }),
  waitForFlash() {
    return waitFor('[data-test-flash-message-body]');
  },
  clickLast() {
    return this.latestItem.click();
  },
  async clickAll() {
    for (let message of this.messages) {
      message.click();
    }
    await settled();
  },
};
