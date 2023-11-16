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
 * The default button is the same as Hds::Button which is primary blue.
 * 
 *
 * @example
//  in dropdown
  <ConfirmAction
    class="hds-confirm-action-critical"
    @buttonColor="critical"
    @isInDropdown={{true}}
    @buttonText="Delete"
    @confirmMessage="This action cannot be undone."
    @onConfirmAction={{log "my action!"}}
  />

  // in toolbar 
  <ConfirmAction
    class="toolbar-button"
    @buttonColor="secondary"
    @buttonText="Delete item"
    @confirmTitle="Delete item?"
    @onConfirmAction={{log "my action!"}}
    @confirmMessage="Are you sure you want to delete this config?"
   />
 * ```
 *
 * @param {Function} onConfirmAction - The action to take upon confirming.
 * @param {String} [confirmTitle=Are you sure?] - The title to display in the confirmation modal.
 * @param {String} [confirmMessage=You will not be able to recover it later.] - The message to display when confirming.
 * @param {Boolean} isInDropdown - If true styles for dropdowns by rendering inside `<li>` elements via `<Hds::Dropdown::ListItem::Interactive`
 * @param {String} buttonText - Text for the button that toggles modal to open.
 * @param {String} [buttonColor] - Color of button that toggles modal. Default is primary, other options are secondary, tertiary, and critical (recommended for dropdowns)
 * @param {String} [modalColor=critical] - Styles modal color, if 'critical' confirm button is also 'critical'. Possible values: critical, warning or neutral
 * @param {Boolean} [isRunning] - Disables the confirm button if action is still running
 * @param {Boolean} [disabled] - Disables the modal's confirm button.
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

  @action
  async onConfirm() {
    await this.args.onConfirmAction();
    // close modal after destructive operation
    this.showConfirmModal = false;
  }
}
