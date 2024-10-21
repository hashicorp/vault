/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module ReplicationOverviewModeComponent
 * ReplicationOverviewMode components are used on the Replication index page to display
 * details about a given mode (DR or Performance) status.
 *
 * @example
 * <ReplicationOverviewModeComponent @mode="dr" @model={{this.cluster.dr}} @clusterName={{this.cluster.name}} />
 *
 * @param {string} mode - should be "dr" or "performance"
 * @param {ReplicationAttributesModel} model - either the dr or performance attribute of the cluster model
 * @param {string} clusterName - used for the link to the mode details
 */
export default class ReplicationOverviewModeComponent extends Component {
  get details() {
    if (this.args.mode === 'dr') {
      return {
        blockTitle: 'Disaster Recovery (DR)',
        upgradeTitle: 'Disaster Recovery is a feature of Vault Enterprise Premium.',
        upgradeLink: 'https://hashicorp.com/products/vault/trial?source=vaultui_DR%20Replication',
        feature: 'DR Replication',
        icon: 'replication-direct',
      };
    }
    return {
      blockTitle: 'Performance',
      upgradeTitle: 'Performance Replication is a feature of Vault Enterprise Premium.',
      upgradeLink: 'https://hashicorp.com/products/vault/trial?source=vaultui_Performance%20Replication',
      feature: 'Performance Replication',
      icon: 'replication-perf',
    };
  }
}
