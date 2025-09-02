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
 *
 *
 * @example
 * <ConfirmAction @buttonColor="critical" @buttonText="Delete" @confirmMessage="This action cannot be undone." @onConfirmAction={{fn (mut this.showConfirmModal) false}} />
 *
 *
 * @param {Function} onConfirmAction - The action to take upon confirming.
 * @param {String} [confirmTitle="Are you sure?"] - The title to display in the confirmation modal.
 * @param {String} [confirmMessage="You will not be able to recover it later."] - The message to display when confirming.
 * @param {Boolean} isInDropdown - If true styles for dropdowns, button color is 'critical', and renders inside `<li>` elements via `<Hds::Dropdown::ListItem::Interactive`
 * @param {String} buttonText - Text for the button that toggles modal to open.
 * @param {String} [buttonColor="primary"] - Color of button that toggles modal, only applies when @isInDropdown=false. Options are primary, secondary (use for toolbars), tertiary, or critical
 * @param {String} [modalColor="critical"] - Styles modal color, if 'critical' confirm button is also 'critical'. Possible values: critical, warning or neutral ('neutral' used for @disabledMessage modal)
 * @param {Boolean} [isRunning] - Disables the modal confirm button - usually a concurrency task that informs the modal if a process is still running
 * @param {String} [disabledMessage] - A message explaining why the confirm action is not allowed, usually combined with a conditional that returns a string if true
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
    assert(`@buttonText is required for ConfirmAction components`, this.args.buttonText);
  }

  get confirmMessage() {
    return this.args.confirmMessage || 'You will not be able to recover it later.';
  }

  get modalColor() {
    if (this.args.disabledMessage) return 'neutral';
    return this.args.modalColor || 'critical';
  }

  @action
  async onConfirm() {
    await this.args.onConfirmAction();
    // close modal after destructive operation
    this.showConfirmModal = false;
  }
}
