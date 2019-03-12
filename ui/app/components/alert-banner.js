import Component from '@ember/component';
import { computed } from '@ember/object';

import { messageTypes } from 'vault/helpers/message-types';
/**
 * `AlertBanners` are used to inform users of important messages.
 * @example
 **  <AlertBanner
 **   @type="danger"
 **   @message="{{model.keyId}} is not a valid lease ID"
 **  />
 **/
export default Component.extend({
  /**
   * The banner type. Should either be `info`, `warning`, `success`, or `danger`.
   * @type {'info' | 'warning' | 'success' | 'danger'}
   * @default null
   * @example 'warning'
   *
   */
  type: null,

  /**
   * The message to display within the banner.
   * @type String
   * @default null
   * @example 'hello'
   */
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
