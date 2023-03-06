import Component from '@glimmer/component';

/**
 * @module AlertPopup
 * The `AlertPopup` is an implementation of the [ember-cli-flash](https://github.com/poteto/ember-cli-flash) `flashMessage`.
 *
 * @example ```js
 * // All properties are passed in from the flashMessage service.
 *   <AlertPopup @type={{message-types flash.type}} @message={{flash.message}} @close={{close}}/>```
 *
 * @param {string} type=null - The alert type. This comes from the message-types helper.
 * @param {string} [message=null] - The alert message.
 * @param {function} close=null - The close action which will close the alert.
 * @param {boolean} isPreformatted - if true modifies class.
 *
 */

export default class AlertPopup extends Component {
  get type() {
    return this.args.type || null;
  }
  get message() {
    return this.args.message || null;
  }
  get close() {
    return this.args.close || null;
  }
}
