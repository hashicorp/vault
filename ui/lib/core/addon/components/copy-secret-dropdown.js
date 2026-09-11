/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

/**
 * @module CopySecretDropdown
 * Renders the toolbar "Copy" menu for a secret. @onWrap is recommended to be a
 * concurrency task, see <Page::Secret::Details> in the KV addon for an example.
 *
 * @param {string} clipboardText - JSON string copied to the clipboard
 * @param {function} onWrap - fires when "Wrap secret" is clicked
 * @param {boolean} isWrapping - renders the wrap item in a loading state
 * @param {string} wrappedData - when present, replaces the wrap item with the wrapped token
 * @param {function} onClose - fires when the dropdown closes
 */
export default class CopySecretDropdown extends Component {
  @service flashMessages;

  @action
  async copyJson(close) {
    try {
      await navigator.clipboard.writeText(this.args.clipboardText);
      this.flashMessages.success('JSON copied!');
    } catch {
      this.flashMessages.danger('Clipboard copy failed. The Clipboard API requires a secure context.');
    }
    close();
  }
}
