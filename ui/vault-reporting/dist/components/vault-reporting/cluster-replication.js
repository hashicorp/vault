import Component from '@glimmer/component';
import { HdsBadge, HdsTextBody, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import { REPLICATION_ENABLED_STATE } from '../../types/index.js';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class ClusterReplication extends Component {
  get disasterRecoveryBadge() {
    return Object.values(REPLICATION_ENABLED_STATE).includes(this.args.disasterRecoveryState) ? {
      icon: 'check',
      text: 'Enabled',
      color: 'success'
    } : {
      icon: 'x',
      text: 'Not set up',
      color: 'neutral'
    };
  }
  get performanceBadge() {
    return Object.values(REPLICATION_ENABLED_STATE).includes(this.args.performanceState) ? {
      icon: 'check',
      text: 'Enabled',
      color: 'success'
    } : {
      icon: 'x',
      text: 'Not set up',
      color: 'neutral'
    };
  }
  get disasterRecoveryRole() {
    return this.args.isDisasterRecoveryPrimary ? 'Primary' : 'Secondary';
  }
  get performanceRole() {
    return this.args.isPerformancePrimary ? 'Primary' : 'Secondary';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-cluster-replication @hasBorder={{true}} class=\"ssu-cluster-replication\" ...attributes>\n      <TitleRow @title=\"Cluster replication status\" @description=\"Check the status and health of Vault clusters\" @linkUrl=\"replication\" />\n\n      <div class=\"ssu-cluster-replication__list-row\" data-test-cluster-replication-dr-row>\n        <HdsTextBody @size=\"300\">\n          Disaster Recovery\n          <HdsBadge class=\"ssu-cluster-replication__list-row__badge\" data-test-cluster-replication-dr-badge @icon={{this.disasterRecoveryBadge.icon}} @text={{this.disasterRecoveryBadge.text}} @color={{this.disasterRecoveryBadge.color}} @type=\"outlined\" @size=\"small\" />\n        </HdsTextBody>\n\n        <HdsTextBody class=\"ssu-cluster-replication__list-row__role\" data-test-cluster-replication-dr-role @size=\"300\" @color=\"var(--token-color-palette-green-400)\">\n          {{this.disasterRecoveryRole}}\n        </HdsTextBody>\n      </div>\n\n      <div class=\"ssu-cluster-replication__list-row\" data-test-cluster-replication-perf-row>\n        <HdsTextBody @size=\"300\">\n          Performance\n          <HdsBadge class=\"ssu-cluster-replication__list-row__badge\" data-test-cluster-replication-perf-badge @icon={{this.performanceBadge.icon}} @text={{this.performanceBadge.text}} @color={{this.performanceBadge.color}} @type=\"outlined\" @size=\"small\" />\n        </HdsTextBody>\n\n        <HdsTextBody class=\"ssu-cluster-replication__list-row__role\" data-test-cluster-replication-perf-role @size=\"300\" @color=\"var(--token-color-palette-green-400)\">\n          {{this.performanceRole}}\n        </HdsTextBody>\n      </div>\n\n    </HdsCardContainer>\n  ", {
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
