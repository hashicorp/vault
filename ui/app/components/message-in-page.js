import Ember from 'ember';

import { messageTypes } from 'vault/helpers/message-types';
const { computed } = Ember;

export default Ember.Component.extend({
  type: null,

  classNameBindings: ['containerClass'],

  containerClass: computed('type', function() {
    return 'message ' + messageTypes([this.get('type')]).class;
  }),

  alertType: computed('type', function() {
    return messageTypes([this.get('type')]);
  }),

  messageClass: 'message-body',
});
