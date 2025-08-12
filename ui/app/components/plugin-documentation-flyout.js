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
 * Contains link to Vault plugin documentation and explains how to enable the plugin.
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

  /**
   * Get the documentation URL for the plugin
   */
  get documentationUrl() {
    const baseUrl = 'https://developer.hashicorp.com/vault/docs';
    const pluginName = this.args.pluginName;

    if (this.args.pluginType === 'secret') {
      // Most secret engines follow this pattern
      return `${baseUrl}/secrets/${pluginName}`;
    } else {
      // Auth methods follow this pattern
      return `${baseUrl}/auth/${pluginName}`;
    }
  }

  /**
   * Get CLI command for enabling the plugin
   */
  get enableCommand() {
    const pluginName = this.args.pluginName;
    const pluginType = this.args.pluginType;

    if (pluginType === 'secret') {
      return `vault secrets enable ${pluginName}`;
    } else {
      return `vault auth enable ${pluginName}`;
    }
  }

  /**
   * Get API endpoint for enabling the plugin
   */
  get apiEndpoint() {
    const pluginName = this.args.pluginName;
    const pluginType = this.args.pluginType;

    if (pluginType === 'secret') {
      return `POST /v1/sys/mounts/${pluginName}`;
    } else {
      return `POST /v1/sys/auth/${pluginName}`;
    }
  }

  @action
  close() {
    if (this.args.onClose) {
      this.args.onClose();
    }
  }
}
