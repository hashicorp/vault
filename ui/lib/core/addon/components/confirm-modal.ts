/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

interface ConfirmModalArgs {
  onConfirm: () => void;
  confirmText?: string;
  confirmLabel?: string;
  confirmButtonText?: string;
}

/**
 * @module ConfirmModal
 * @description
 * ConfirmModal components are used to allow users to confirm an action.
 * Supports an optional @confirmText arg — when provided,
 * the user must type the given string before the confirm button is enabled (type-to-confirm).
 *
 * @example
 * Simple confirmation:
 * <ConfirmModal @confirmTitle="Delete item?" @onClose={{this.close}} @onConfirm={{this.delete}} />
 *
 * Type-to-confirm:
 * <ConfirmModal @confirmTitle="Delete item?" @confirmText="my-item" @onClose={{this.close}} @onConfirm={{this.delete}} />
 */

export default class ConfirmModal extends Component<ConfirmModalArgs> {
  @tracked confirmInput = '';
  @tracked showConfirmWarning = false;

  get isConfirmDisabled() {
    const { confirmText } = this.args;
    if (!confirmText) return false;
    return this.confirmInput !== confirmText;
  }

  @action
  onInput(event: Event) {
    this.confirmInput = (event.target as HTMLInputElement).value;
    this.showConfirmWarning = false;
  }

  @action
  saveAndClose(close: () => void) {
    if (this.isConfirmDisabled) {
      this.showConfirmWarning = true;
      return;
    }
    close();
    if (this.args?.onConfirm) {
      this.args.onConfirm();
    }
  }
}
