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
 * @param {String} [bannerType=alert] - Defaults to 'alert', can be used to specify an alert banner's test selector
 * @param {boolean} [marginless=false] - xx
 * @param {String} [message=null] - The message to display within the banner.
 * @param {Object} [progressBar=null] - An object containing a value and maximum for a progress bar. Will be displayed next to the message title.
 * @param {boolean} [showLoading=false] - Showing a loading icon to the right of the title.
 * @param {String} [title=null] - A title to show above the message. If this is not provided, there are default values for each type of alert.
 * @param {boolean} [yieldWithoutColumn=false] - If true, do not show message or title, just yield with no formatting.
 * @param {string} [learnLink=null] -display a DocLink pointing to this path.
 * @param {string} [learnLinkMessage=null] -show a learn link message to the right of a link icon, displayed below the main message.
 *
 */

export default class AlertBanner extends Component {
  get bannerType() {
    return this.args.bannerType || 'alert';
  }

  get containerClass() {
    const base = this.args.marginless ? 'message message-marginless ' : 'message ';
    return base + messageTypes([this.args.type]).class;
  }

  get alertType() {
    return messageTypes([this.args.type]);
  }
}
