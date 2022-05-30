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
 * @param type=null{String} - The alert type passed to the message-types helper.
 * @param [message=null]{String} - The message to display within the alert.
 * @param [paddingTop=false]{Boolean} - Whether or not to add padding above component.
 * @param [isMarginless=false]{Boolean} - Whether or not to remove margin bottom below component.
 * @param [sizeSmall=false]{Boolean} - Whether or not to display a small font with padding below of alert message.
 * @param [mimicRefresh=false]{Boolean} - If true will display a loading icon when attributes change (e.g. when a form submits and the alert message changes).
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
