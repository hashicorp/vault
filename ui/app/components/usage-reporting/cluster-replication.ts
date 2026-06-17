/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { htmlSafe } from '@ember/template';

import { REPLICATION_ENABLED_STATE } from 'vault/types/usage-reporting';

interface VaultReportingClusterReplicationSignature {
  Args: {
    disasterRecoveryState: REPLICATION_ENABLED_STATE | 'disabled';
    performanceState: REPLICATION_ENABLED_STATE | 'disabled';
    isVaultDedicated: boolean;
  };
}

export default class VaultReportingClusterReplication extends Component<VaultReportingClusterReplicationSignature> {
  getState = (state: REPLICATION_ENABLED_STATE | 'disabled' = 'disabled') => {
    return state;
  };

  get isEmpty() {
    return (
      this.getState(this.args.disasterRecoveryState) === 'disabled' &&
      this.getState(this.args.performanceState) === 'disabled'
    );
  }

  get description() {
    if (this.isEmpty) {
      return htmlSafe(
        'Enable <a class="hds-link-inline--color-secondary" href="https://developer.hashicorp.com/vault/docs/internals/replication" target="_blank" rel="noopener noreferrer" data-test-vault-reporting-cluster-replication-description-link>replication</a> to replicate data across clusters.'
      );
    }

    return 'Status of disaster recovery and performance replication.';
  }

  getIcon = (state: REPLICATION_ENABLED_STATE | 'disabled' = 'disabled') => {
    const iconMap: Record<string, string> = {
      disabled: 'x',
      primary: 'check',
      secondary: 'check',
      bootstrapping: 'loading',
    };

    return iconMap[state] || iconMap['disabled'];
  };

  getColor = (state: REPLICATION_ENABLED_STATE | 'disabled' = 'disabled') => {
    const colorMap: Record<string, string> = {
      disabled: 'neutral',
      primary: 'success',
      secondary: 'success',
      bootstrapping: 'neutral',
    };

    return colorMap[state] || colorMap['disabled'];
  };

  get linkRoute() {
    if (this.args.isVaultDedicated) {
      return;
    }

    return 'vault.cluster.replication';
  }
}
