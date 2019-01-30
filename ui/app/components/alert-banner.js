import Component from '@ember/component';
import { computed } from '@ember/object';

import { messageTypes } from 'vault/helpers/message-types';

export default Component.extend({
  type: null,

  yieldWithoutColumn: false,

  classNameBindings: ['containerClass'],

  containerClass: computed('type', function() {
    return 'message ' + messageTypes([this.get('type')]).class;
  }),

  alertType: computed('type', function() {
    return messageTypes([this.get('type')]);
  }),
});
