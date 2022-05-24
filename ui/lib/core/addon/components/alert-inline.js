import Component from '@ember/component';
import { computed } from '@ember/object';
import { messageTypes } from 'core/helpers/message-types';
import layout from '../templates/components/alert-inline';
import { later } from '@ember/runloop';

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
 * @param [sizeSmall=false]{Boolean} - Whether or not to display a small font with padding below of alert message.
 * @param [paddingTop=false]{Boolean} - Whether or not to add padding above component.
 * @param [isMarginless=false]{Boolean} - Whether or not to remove margin bottom below component.
 * @param [mimicRefresh=false]{Boolean} - If true will display a loading icon when attributes change (e.g. when a form submits and the alert message changes).
 */

export default Component.extend({
  layout,
  type: null,
  message: null,
  sizeSmall: false,
  paddingTop: false,
  classNames: ['message-inline'],
  classNameBindings: ['sizeSmall:size-small', 'paddingTop:padding-top', 'isMarginless:is-marginless'],
  isLoading: false,
  mimicRefresh: false,

  textClass: computed('type', function () {
    if (this.type == 'danger') {
      return messageTypes([this.type]).glyphClass;
    }

    return;
  }),

  alertType: computed('type', function () {
    return messageTypes([this.type]);
  }),

  didReceiveAttrs() {
    this._super(...arguments);
    if (this.mimicRefresh) {
      this.set('isLoading', true);
      later(() => {
        this.set('isLoading', false);
      }, 200);
    }
  },
});
