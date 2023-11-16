/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { tracked } from '@glimmer/tracking';

/**
 * @module ConfirmAction
 * ConfirmAction is a button that opens a modal containing a confirmation message with confirm or cancel action.
 * Splattributes are spread to the button element to apply styling directly without adding extra args.
 * The default button is the same as Hds::Button which is primary blue
 *
 * @example
 * ```js
 *  <ConfirmAction
 *    @buttonText="Delete item"
 *    @onConfirmAction={{ () => { console.log('Action!') } }}
 *    @confirmMessage="Are you sure you want to delete this config?"
 *  />
 * ```
 *
 * @param {Func} onConfirmAction - The action to take upon confirming.
 * @param {String} [confirmTitle=Are you sure?] - The title to display in the confirmation modal.
 * @param {String} [confirmMessage=You will not be able to recover it later.] - The message to display when confirming.
 * @param {String} buttonText - Text for the button that triggers modal to open.
 * @param {String} buttonColor - Color of button that triggers modal. Default is primary, other options are secondary, tertiary, and critical
 * @param {Boolean} [isRunning] - Disables the confirm button if action is still running
 * @param {Boolean} [disabled] - To disable the confirm button.
 * @param {String} [color=critical] - Changes the color of modal and confirm button
 *
 */

export default class ConfirmActionComponent extends Component {
  @tracked showConfirmModal = false;

  constructor() {
    super(...arguments);
    assert(
      '<ConfirmAction> component expects @onConfirmAction arg to be a function',
      typeof this.args.onConfirmAction === 'function'
    );
  }

  get confirmMessage() {
    return this.args.confirmMessage || 'You will not be able to recover it later.';
  }

  @action
  async onConfirm() {
    await this.args.onConfirmAction();
    // close modal after destructive operation
    this.showConfirmModal = false;
  }
}
