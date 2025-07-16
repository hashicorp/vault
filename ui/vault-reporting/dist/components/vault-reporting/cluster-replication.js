import Component from '@glimmer/component';
import { HdsBadge, HdsTextBody, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import { htmlSafe } from '@ember/template';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class ClusterReplication extends Component {
  getState = (state = 'disabled') => {
    return state;
  };
  get isEmpty() {
    return this.getState(this.args.disasterRecoveryState) === 'disabled' && this.getState(this.args.performanceState) === 'disabled';
  }
  get description() {
    if (this.isEmpty) {
      return htmlSafe('Enable <a class="hds-link-inline--color-secondary" href="https://developer.hashicorp.com/vault/docs/internals/replication" target="_blank" data-test-vault-reporting-cluster-replication-description-link>replication</a> to replicate data across clusters.');
    } else {
      return 'Status of disaster recovery and performance replication.';
    }
  }
  getIcon = (state = 'disabled') => {
    const iconMap = {
      disabled: 'x',
      primary: 'check',
      secondary: 'check',
      bootstrapping: 'loading'
    };
    return iconMap[state] || iconMap['disabled'];
  };
  getColor = (state = 'disabled') => {
    const colorMap = {
      disabled: 'neutral',
      primary: 'success',
      secondary: 'success',
      bootstrapping: 'neutral'
    };
    return colorMap[state] || colorMap['disabled'];
  };
  get linkUrl() {
    const {
      isVaultDedicated = false
    } = this.args;
    if (isVaultDedicated) {
      return;
    }
    return 'replication';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-vault-reporting-cluster-replication @hasBorder={{true}} class=\"ssu-cluster-replication\" ...attributes>\n      <TitleRow @title=\"Cluster replication\" @description={{this.description}} @linkUrl={{this.linkUrl}} />\n\n      <HdsTextBody @size=\"300\" data-test-vault-reporting-cluster-replication-dr-row>\n        Disaster Recovery\n        <HdsBadge class=\"ssu-cluster-replication__list-row__badge\" data-test-vault-reporting-cluster-replication-dr-badge @icon={{this.getIcon @disasterRecoveryState}} @text={{this.getState @disasterRecoveryState}} @color={{this.getColor @disasterRecoveryState}} @type=\"outlined\" @size=\"small\" />\n      </HdsTextBody>\n\n      <HdsTextBody @size=\"300\" data-test-vault-reporting-cluster-replication-perf-row>\n        Performance\n        <HdsBadge class=\"ssu-cluster-replication__list-row__badge\" data-test-vault-reporting-cluster-replication-perf-badge @icon={{this.getIcon @performanceState}} @text={{this.getState @performanceState}} @color={{this.getColor @performanceState}} @type=\"outlined\" @size=\"small\" />\n      </HdsTextBody>\n\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        TitleRow,
        HdsTextBody,
        HdsBadge
      })
    }), this);
  }
}

export { ClusterReplication as default };
//# sourceMappingURL=cluster-replication.js.map
