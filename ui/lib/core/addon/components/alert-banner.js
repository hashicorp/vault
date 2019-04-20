import Component from '@ember/component';
import { computed } from '@ember/object';
import { messageTypes } from 'vault/helpers/message-types';

/**
 * @module AlertBanner
 * `AlertBanner` components are used to inform users of important messages.
 *
 * @example
 * ```js
 * <AlertBanner @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
 * ```
 *
 * @param type=null {String} - The banner type. This comes from the message-types helper.
 * @param [message=null {String}] - The message to display within the banner.
 *
 */

export default Component.extend({
  type: null,
  message: null,
  yieldWithoutColumn: false,
  classNameBindings: ['containerClass'],

  containerClass: computed('type', function() {
    return 'message ' + messageTypes([this.get('type')]).class;
  }),

  alertType: computed('type', function() {
    return messageTypes([this.get('type')]);
  }),
});
