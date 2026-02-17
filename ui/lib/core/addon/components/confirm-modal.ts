/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

interface ConfirmModalArgs {
  onConfirm: () => void;
}

/**
 * @module ConfirmModal
 * @description
 * ConfirmModal components are used to allow users to select any number of predetermined options, aligned in a 3-column grid.
 *
 *
 * @example
 * <ConfirmModal @name="extKeyUsage" @label="Extended key usage" @fields={{array (hash key="EmailProtection" label="Email Protection") (hash key="TimeStamping" label="Time Stamping") (hash key="ServerAuth" label="Server Auth") }} @value={{array "TimeStamping"}} />
 */

export default class ConfirmModal extends Component<ConfirmModalArgs> {
  @action
  saveAndClose(close: () => void) {
    close();
    if (this.args?.onConfirm) {
      this.args.onConfirm();
    }
  }
}
