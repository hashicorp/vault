/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module PageModeIndex
 *
 * @example
 * <Page::ModeIndex
 *  @cluster={{this.model}}
 *  @onEnableSuccess={{this.onEnableSuccess}}
 *  @replicationDisabled={{this.replicationForMode.replicationDisabled}
 *  @replicationMode={{this.replicationMode}}
 * />
 *
 * @param {model} cluster - cluster route model
 * @param {function} onEnableSuccess - callback after enabling is successful, handles transition if enabled from the top-level index route
 * @param {boolean} replicationDisabled - whether or not replication is enabled
 * @param {string} replicationMode - should be "dr" or "performance"
 */
export default class PageModeIndex extends Component {
  get canEnablePrimary() {
    const { cluster } = this.args;
    switch (this.args.replicationMode) {
      case 'dr':
        return cluster.canEnablePrimaryDr;
      case 'performance':
        return cluster.canEnablePrimaryPerformance;
      default:
        // if there's a problem checking capabilities, default to true
        // since the backend will gate as a fallback
        return true;
    }
  }
  get canEnableSecondary() {
    const { cluster } = this.args;
    switch (this.args.replicationMode) {
      case 'dr':
        return cluster.canEnableSecondaryDr;
      case 'performance':
        return cluster.canEnableSecondaryPerformance;
      default:
        // if there's a problem checking capabilities, default to true
        // since the backend will gate as a fallback
        return true;
    }
  }
}
