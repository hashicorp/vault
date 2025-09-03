/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

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
 *   @plugin={{this.flyoutPlugin}}
 *   @pluginType="secret"
 *   @onClose={{this.closeFlyout}}
 * />
 * ```
 * @param {boolean} isOpen - Whether the flyout is currently open
 * @param {Object} plugin - Plugin object containing type, displayName, etc.
 * @param {string} pluginType - Type of plugin (secret, auth)
 * @param {Function} onClose - Callback when flyout is closed
 */

export default class PluginDocumentationFlyoutComponent extends Component {
  // Get the display name for the plugin
  get displayName() {
    // Support both plugin object and direct displayName/pluginName args
    return (
      this.args.displayName ||
      this.args.plugin?.displayName ||
      this.args.pluginName ||
      this.args.plugin?.type ||
      'Plugin'
    );
  }

  // Get the plugin type for display
  get pluginTypeDisplay() {
    return this.args.pluginType === 'secret' ? 'secrets engine' : 'auth method';
  }

  // Check if we have plugin-specific data
  get hasPlugin() {
    return !!(this.args.plugin || this.args.pluginName || this.args.displayName);
  }
}
