/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { tracked } from '@glimmer/tracking';

/**
 * @module ConfirmAction
 * `ConfirmAction` is a button followed by a pop up confirmation message and button used to prevent users from performing actions they do not intend to.
 *
 * @example
 * ```js
 *  <ConfirmAction
 *    @onConfirmAction={{ () => { console.log('Action!') } }}
 *    @confirmMessage="Are you sure you want to delete this config?">
 *    Delete
 *  </ConfirmAction>
 * ```
 *
 * @param {Func} [onConfirmAction=null] - The action to take upon confirming.
 * @param {String} [confirmTitle=Delete this?] - The title to display when confirming.
 * @param {String} [confirmMessage=You will not be able to recover it later.] - The message to display when confirming.
 * @param {String} [confirmButtonText=Delete] - The confirm button text.
 * @param {String} [cancelButtonText=Cancel] - The cancel button text.
 * @param {String} [buttonClasses] - A string to indicate the button class.
 * @param {String} [horizontalPosition=auto-right] - For the position of the dropdown.
 * @param {String} [verticalPosition=below] - For the position of the dropdown.
 * @param {Boolean} [isRunning=false] - If action is still running disable the confirm.
 * @param {Boolean} [disable=false] - To disable the confirm action.
 *
 */

export default class ConfirmActionComponent extends Component {
  @tracked showConfirm = false;

  get horizontalPosition() {
    return this.args.horizontalPosition || 'auto-right';
  }

  get verticalPosition() {
    return this.args.verticalPosition || 'below';
  }

  get isRunning() {
    return this.args.isRunning || false;
  }

  get disabled() {
    return this.args.disabled || false;
  }

  get confirmTitle() {
    return this.args.confirmTitle || 'Delete this?';
  }

  get confirmMessage() {
    return this.args.confirmMessage || 'You will not be able to recover it later.';
  }

  get confirmButtonText() {
    return this.args.confirmButtonText || 'Delete';
  }

  get cancelButtonText() {
    return this.args.cancelButtonText || 'Cancel';
  }

  @action
  toggleConfirm() {
    // toggle
    this.showConfirm = !this.showConfirm;
  }

  @action
  onConfirm(actions) {
    const confirmAction = this.args.onConfirmAction;

    if (typeof confirmAction !== 'function') {
      assert('confirm-action components expects `onConfirmAction` attr to be a function');
    } else {
      confirmAction();
      // close the dropdown content
      actions.close();
    }
  }
}
