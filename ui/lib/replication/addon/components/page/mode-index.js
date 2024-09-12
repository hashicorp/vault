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
  canEnable = (type) => {
    const { cluster, replicationMode } = this.args;
    let perm;
    if (replicationMode === 'dr') {
      // returns canEnablePrimaryDr or canEnableSecondaryDr
      perm = `canEnable${type}Dr`;
    }
    if (replicationMode === 'performance') {
      // returns canEnablePrimaryPerformance or canEnableSecondaryPerformance
      perm = `canEnable${type}Performance`;
    }
    // if there's a problem checking capabilities, default to true
    // since the backend can gate as a fallback
    return cluster[perm] ?? true;
  };
}
