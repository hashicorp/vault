/**
 * @module ConfirmationModal
 * ConfirmationModal components are used to provide an alternative to ConfirmationButton that automatically prompts the user to fill in confirmation text before they can continue with a potentially destructive action. It is built off the Modal component
 *
 * @example
 * ```js
 * <ConfirmationModal
 *   @onConfirm={action "destructiveAction"}
 *   @title="Do Dangerous Thing?"
 *   @isActive={{isModalActive}}
 *   @onClose={{action (mut isModalActive) false}}
 * />
 * ```
 * @param {function} onConfirm - onConfirm is the action that happens when user clicks onConfirm after filling in the confirmation block
 * @param {boolean} isActive - Controls whether the modal is "active" eg. visible or not.
 * @param {string} title - Title of the modal
 * @param {function} onClose - specify what to do when user attempts to close modal
 * @param {string} [buttonText=Confirm] - Button text on the confirm button
 * @param {string} [confirmText=Yes] - The confirmation text that the user must type before continuing
 * @param {string} [buttonClass=is-danger] - extra class to add to confirm button (eg. "is-danger")
 * @param {sting} [type=warning] - Applies message-type styling to header. Override to default with empty string
 * @param {string} [toConfirmMsg] - Finishes the sentence "Type YES to confirm ..."
 */

import Component from '@ember/component';
import layout from '../templates/components/confirmation-modal';

export default Component.extend({
  layout,
  buttonClass: 'is-danger',
  buttonText: 'Confirm',
  confirmText: 'Yes',
  type: 'warning',
  actionDescription: '',
  toConfirmMsg: '',
});
