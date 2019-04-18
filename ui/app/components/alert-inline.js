import Component from '@ember/component';
import { computed } from '@ember/object';
import { messageTypes } from 'vault/helpers/message-types';

/**
 * @module AlertInline
 * `AlertInline` components are used to inform users of important messages.
 *
 * @example
 * ```js
 * <AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
 * ```
 *
 * @param type=null{String} - The alert type. This comes from the message-types helper.
 * @param [message=null]{String} - The message to display within the alert.
 *
 */

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
