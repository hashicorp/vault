/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
/**
 * @module ConfirmationModal
 * ConfirmationModal components wrap the <Hds::Modal> component to present a critical (red) type-to-confirm modal.
 * They are used for extremely destructive actions that require extra consideration before confirming.
 *
 * @example
 * <ConfirmationModal
 *   @onConfirm={action "destructiveAction"}
 *   @title="Do Dangerous Thing?"
 *   @isActive={{isModalActive}}
 *   @onClose={{action (mut isModalActive) false}}
 *   @confirmText="yes"
 *   @onConfirmMsg="deleting this thing to delete."
 * />
 *
 * @param {function} onConfirm - onConfirm is the action that happens when user clicks onConfirm after filling in the confirmation block
 * @param {function} onClose - specify what to do when user attempts to close modal
 * @param {boolean} isActive - Controls whether the modal is "active" eg. visible or not.
 * @param {string} title - Title of the modal
 * @param {string} [confirmText=Yes] - The confirmation text that the user must type before continuing
 * @param {string} [toConfirmMsg] - Finishes the sentence "Type <confirmText> to confirm <toConfirmMsg>", default is an empty string (ex. 'secret deletion')
 * @param {string} [buttonText=Confirm] - Button text on the confirm button
 */

export default class ConfirmationModal extends Component {
  get confirmText() {
    return this.args.confirmText || 'Yes';
  }
}
