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
 * @param {String} type=null  - The banner type. This comes from the message-types helper.
 * @param {String} [secondIconType=null] - If you want a second icon to appear to the right of the title. This comes from the message-types helper.
 * @param {Object} [progressBar=null] - An object containing a value and maximum for a progress bar. Will be displayed next to the message title.
 * @param {String} [message=null] - The message to display within the banner.
 * @param {String} [title=null] - A title to show above the message. If this is not provided, there are default values for each type of alert.
 * @param {String} [bannerType=alert] - Defaults to 'alert', can be used to specify an alert banner's test selector
 *
 */

export default Component.extend({
  layout,
  type: null,
  message: null,
  title: null,
  secondIconType: null,
  progressBar: null,
  yieldWithoutColumn: false,
  marginless: false,
  classNameBindings: ['containerClass'],
  bannerType: 'alert',

  containerClass: computed('type', 'marginless', function () {
    const base = this.marginless ? 'message message-marginless ' : 'message ';
    return base + messageTypes([this.type]).class;
  }),

  alertType: computed('type', function () {
    return messageTypes([this.type]);
  }),

  secondAlertType: computed('secondIconType', function () {
    return messageTypes([this.secondIconType]);
  }),
});
