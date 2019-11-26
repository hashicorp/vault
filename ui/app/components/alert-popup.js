import OuterHTML from './outer-html';

/**
 * @module AlertPopup
 * The `AlertPopup` is an implementation of the [ember-cli-flash](https://github.com/poteto/ember-cli-flash) `flashMessage`.
 *
 * @example ```js
 * // All properties are passed in from the flashMessage service.
 *   <AlertPopup @type={{message-types flash.type}} @message={{flash.message}} @close={{close}}/>```
 *
 * @param type=null {String} - The alert type. This comes from the message-types helper.
 * @param [message=null] {String} - The alert message.
 * @param close=null {Func} - The close action which will close the alert.
 *
 */

export default OuterHTML.extend({
  type: null,
  message: null,
  close: null,
});
