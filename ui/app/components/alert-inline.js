import Component from '@ember/component';
import { computed } from '@ember/object';

import { messageTypes } from 'vault/helpers/message-types';

export default Component.extend({
  type: null,
  message: null,

  classNames: ['message-inline'],

  textClass: computed('type', function() {
    if (this.get('type') == 'danger') {
      return messageTypes([this.get('type')]).glyphClass;
    }

    return;
  }),

  alertType: computed('type', function() {
    return messageTypes([this.get('type')]);
  }),
});
