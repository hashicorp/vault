import { getWithDefault, computed } from '@ember/object';
import FlashMessage from 'ember-cli-flash/components/flash-message';

export default FlashMessage.extend({
  // override alertType to get Bulma specific prefix
  //https://github.com/poteto/ember-cli-flash/blob/master/addon/components/flash-message.js#L35
  alertType: computed('flash.type', {
    get() {
      const flashType = getWithDefault(this, 'flash.type', '');
      let prefix = 'is-';

      return `${prefix}${flashType}`;
    },
  }),
});
