import Component from '@glimmer/component';
import { action } from '@ember/object';
import { later } from '@ember/runloop';
import { tracked } from '@glimmer/tracking';
import { messageTypes } from 'core/helpers/message-types';

/**
 * @module AlertInline
 * `AlertInline` components are used to inform users of important messages.
 *
 * @example
 * ```js
 * <AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
 * ```
 *
 * @param {string} type=null - The alert type passed to the message-types helper.
 * @param {string} [message=null] - The message to display within the alert.
 * @param {boolean} [paddingTop=false] - Whether or not to add padding above component.
 * @param {boolean} [isMarginless=false] - Whether or not to remove margin bottom below component.
 * @param {boolean} [sizeSmall=false] - Whether or not to display a small font with padding below of alert message.
 * @param {boolean} [mimicRefresh=false] - If true will display a loading icon when attributes change (e.g. when a form submits and the alert message changes).
 * @param {number} [formErrorCount=null] - Number of errors (i.e. error.errors.length) to display an inline alert (typically be the submit button).
 */

export default class AlertInlineComponent extends Component {
  @tracked isRefreshing = false;

  get mimicRefresh() {
    return this.args.mimicRefresh || false;
  }

  get paddingTop() {
    return this.args.paddingTop ? ' padding-top' : '';
  }

  get isMarginless() {
    return this.args.isMarginless ? ' is-marginless' : '';
  }

  get sizeSmall() {
    return this.args.sizeSmall ? ' size-small' : '';
  }

  get textClass() {
    if (this.args.type === 'danger') {
      return this.alertType.glyphClass;
    }
    return null;
  }

  get alertType() {
    return messageTypes([this.args.type]);
  }

  get displayMessage() {
    const { formErrorCount, message } = this.args;
    return formErrorCount ? this.generateErrorCountMessage(formErrorCount) : message;
  }

  generateErrorCountMessage(errorCount) {
    if (errorCount < 1) return null;
    // returns count specific message: 'There is an error/are N errors with this form.'
    let isPlural = errorCount > 1 ? `are ${errorCount} errors` : false;
    return `There ${isPlural ? isPlural : 'is an error'} with this form.`;
  }

  @action
  refresh() {
    if (this.mimicRefresh) {
      this.isRefreshing = true;
      later(() => {
        this.isRefreshing = false;
      }, 200);
    }
  }
}
