import FlashMessages from 'ember-cli-flash/services/flash-messages';

export default FlashMessages.extend({
  stickyInfo(message) {
    return this.info(message, {
      sticky: true,
      priority: 300,
    });
  },
});
