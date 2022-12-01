import FlashMessages from 'ember-cli-flash/services/flash-messages';

export default class FlashMessageService extends FlashMessages {
  stickyInfo(message: string) {
    return this.info(message, {
      sticky: true,
      priority: 300,
    });
  }
}
