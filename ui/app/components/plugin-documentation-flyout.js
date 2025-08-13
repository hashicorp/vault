/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * Plugin Documentation Flyout Component
 *
 * Modal/flyout that opens when disabled plugin is clicked.
 * Displays information about plugins and links to general Vault plugin documentation.
 *
 * @example
 * ```js
 * <PluginDocumentationFlyout
 *   @isOpen={{this.showFlyout}}
 *   @pluginName="my-plugin"
 *   @pluginType="secret"
 *   @onClose={{this.closeFlyout}}
 * />
 * ```
 * @param {boolean} isOpen - Whether the flyout is currently open
 * @param {string} pluginName - Name of the plugin
 * @param {string} pluginType - Type of plugin (secret, auth)
 * @param {string} [displayName] - Display name of the plugin (defaults to pluginName)
 * @param {Function} onClose - Callback when flyout is closed
 */

export default class PluginDocumentationFlyoutComponent extends Component {
  /**
   * Get the display name for the plugin
   */
  get displayName() {
    return this.args.displayName || this.args.pluginName;
  }

  /**
   * Get the plugin type for display
   */
  get pluginTypeDisplay() {
    return this.args.pluginType === 'secret' ? 'secrets engine' : 'auth method';
  }

  @action
  close() {
    if (this.args.onClose) {
      this.args.onClose();
    }
  }
}
