import Component from '@ember/component';
import { computed } from '@ember/object';
import { messageTypes } from 'core/helpers/message-types';
import layout from '../templates/components/alert-banner';

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
 * @param secondIconType=null {String} - Optional. If you want a second icon to appear to the right of the title. This comes from the message-types helper.
 * @param [message=null {String}] - The message to display within the banner.
 * @param [title=null {String}] - A title to show above the message. If this is not provided, there are default values for each type of alert.
 *
 */

export default Component.extend({
  layout,
  type: null,
  message: null,
  title: null,
  secondIconType: null,
  yieldWithoutColumn: false,
  classNameBindings: ['containerClass'],

  containerClass: computed('type', function() {
    return 'message ' + messageTypes([this.get('type')]).class;
  }),

  alertType: computed('type', function() {
    return messageTypes([this.get('type')]);
  }),

  secondAlertType: computed('secondIconType', function() {
    return messageTypes([this.get('secondIconType')]);
  }),
});
