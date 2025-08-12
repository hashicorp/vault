/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * Plugin Status Indicator Component
 *
 * Displays plugin version information with badges to show builtin vs external status.
 * Provides consistent plugin status display across the application.
 *
 * @example
 * ```js
 * <PluginStatusIndicator
 *   @version="v1.12.0+builtin.vault"
 *   @builtin={{true}}
 *   @deprecationStatus="supported"
 * />
 * ```
 * @param {string} version - Plugin version string
 * @param {boolean} builtin - Whether the plugin is builtin or external
 * @param {string} [deprecationStatus] - Deprecation status of the plugin
 * @param {string} [size="small"] - Size of the badges (small, medium, large)
 */

export default class PluginStatusIndicatorComponent extends Component {
  /**
   * Get the version to display
   */
  get displayVersion() {
    return this.args.version;
  }

  /**
   * Get the plugin type badge configuration
   */
  get typeBadgeConfig() {
    if (this.args.builtin === undefined) return null;

    return this.args.builtin
      ? { text: 'Builtin', color: 'neutral' }
      : { text: 'External', color: 'highlight' };
  }

  /**
   * Get the deprecation badge configuration if applicable
   */
  get deprecationBadgeConfig() {
    if (!this.args.deprecationStatus || this.args.deprecationStatus === 'supported') {
      return null;
    }

    if (this.args.deprecationStatus === 'pending-removal') {
      return { text: 'Pending Removal', color: 'warning' };
    }

    if (this.args.deprecationStatus === 'removed') {
      return { text: 'Removed', color: 'critical' };
    }

    return { text: 'Deprecated', color: 'warning' };
  }

  /**
   * Get the badge size to use
   */
  get badgeSize() {
    return this.args.size || 'small';
  }
}
