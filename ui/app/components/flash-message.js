import Ember from 'ember';
import FlashMessage from 'ember-cli-flash/components/flash-message';

const { computed, getWithDefault } = Ember;

export default FlashMessage.extend({
  // override alertType to get Bulma specific prefix
  //https://github.com/poteto/ember-cli-flash/blob/master/addon/components/flash-message.js#L35
  alertType: computed('flash.type', {
    get() {
      const flashType = getWithDefault(this, 'flash.type', '');
      let prefix = 'notification has-border is-';

      return `${prefix}${flashType}`;
    },
  }),
});
