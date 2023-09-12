/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
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
 *   @onConfirmMsg="deleting this thing to delete."
 * />
 * ```
 * @param {function} onConfirm - onConfirm is the action that happens when user clicks onConfirm after filling in the confirmation block
 * @param {function} onClose - specify what to do when user attempts to close modal
 * @param {boolean} isActive - Controls whether the modal is "active" eg. visible or not.
 * @param {string} title - Title of the modal
 * @param {string} [confirmText=Yes] - The confirmation text that the user must type before continuing
 * @param {string} [toConfirmMsg=''] - Finishes the sentence "Type <confirmText> to confirm <toConfirmMsg>", default is an empty string (ex. 'secret deletion')
 * @param {string} [buttonText=Confirm] - Button text on the confirm button
 * @param {string} [buttonClass=is-danger] - extra class to add to confirm button (eg. "is-danger")
 * @param {string} [type=warning] - The header styling based on type, passed into the message-types helper (in the Modal component).
 */

export default class ConfirmationModal extends Component {
  get buttonClass() {
    return this.args.buttonClass || 'is-danger';
  }

  get buttonText() {
    return this.args.buttonText || 'Confirm';
  }

  get confirmText() {
    return this.args.confirmText || 'Yes';
  }

  get type() {
    return this.args.type || 'warning';
  }

  get toConfirmMsg() {
    return this.args.toConfirmMsg || '';
  }
}
