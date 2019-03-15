import OuterHTML from './outer-html';

/**
 * @module AlertPopup
 * The `AlertPopup` is an implementation of the {@link https://github.com/poteto/ember-cli-flash|ember-cli-flash} `flashMessage`.
 *
 * @example ```js
 * // All properties are passed in from the flashMessage service.
 *   <AlertPopup @type={{message-types flash.type}} @message={{flash.message}} @close={{close}}/>```
 *
 * @property [AlertPopup.type=null] {String} - The alert type. This comes from the message-types helper.
 * @property message=null {String} - The alert message.
 * @property close=null {Func} - The close action which will close the alert.
 *
 * @see {@link https://github.com/hashicorp/vault/search?l=Handlebars&q=AlertPopup|Uses of AlertPopup}
 * @see {@link https://github.com/hashicorp/vault/blob/master/ui/app/components/alert-popup.js|AlertPopup Source Code}
 */

export default OuterHTML.extend({
  type: null,
  message: null,
  close: null,
});
