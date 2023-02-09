import Component from '@glimmer/component';
import { messageTypes } from 'core/helpers/message-types';

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
 * @param {String} [message=null] - The message to display within the banner.
 * @param {Object} [progressBar=null] - An object containing a value and maximum for a progress bar. Will be displayed next to the message title.
 * @param {Boolean} [showLoading=false] - Shows a loading icon to the right of the title.
 * @param {String} [title=null] - A title to show above the message. If this is not provided, there are default values for each type of alert.
 */

export default class AlertBanner extends Component {
  get alertType() {
    return messageTypes([this.args.type]);
  }
}
