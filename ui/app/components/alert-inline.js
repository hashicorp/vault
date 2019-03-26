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
 * @property [AlertInline.type=null]{String} - The alert type. This comes from the message-types helper.
 * @property message=null{String} - The message to display within the alert.
 *
 * @see {@link https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertInline|Uses of AlertInline}
 * @see {@link https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-inline.js|AlertInline Source Code}
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
