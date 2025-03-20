import Component from '@glimmer/component';
import SSUReportingCounter from '../reporting/counter.js';
import SSUReportingHorizontalBarChart from '../reporting/horizontal-bar-chart.js';
import GlobalLease from '../reporting/global-lease.js';
import ClusterReplication from '../reporting/cluster-replication.js';
import { tracked } from '@glimmer/tracking';
import { on } from '@ember/modifier';
import { HdsTextDisplay, HdsSeparator, HdsCardContainer, HdsButton, HdsLinkInline, HdsPageHeader } from '@hashicorp/design-system-components/components';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUViewDashboard extends Component {
  static {
    g(this.prototype, "data", [tracked]);
  }
  #data = (i(this, "data"), undefined);
  static {
    g(this.prototype, "lastUpdatedTime", [tracked], function () {
      return '';
    });
  }
  #lastUpdatedTime = (i(this, "lastUpdatedTime"), undefined);
  constructor(owner, args) {
    super(owner, args);
    this.fetchAllData().catch(error => {
      console.error('Error fetching data', error);
    });
  }
  fetchAllData = async () => {
    this.data = await this.args.service.getUsageData();
    this.lastUpdatedTime = new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'medium'
    }).format(new Date());
  };
  getBarChartData = map => {
    return Object.entries(map).map(([label, value]) => {
      return {
        label,
        value
      };
    });
  };
  get counters() {
    return [{
      title: 'Child namespaces',
      data: this.data?.namespaces ?? 0,
      link: '/access/namespaces'
    }, {
      title: 'KV secrets',
      data: (this.data?.kvv1_secrets ?? 0) + (this.data?.kvv2_secrets ?? 0)
    }, {
      title: 'Secrets sync',
      data: this.data?.secrets_sync ?? 0,
      link: '/sync/secrets'
    }, {
      title: 'PKI roles',
      data: this.data?.pki.total_roles ?? 0
    }];
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div class=\"dashboard\">\n      <HdsPageHeader as |PH|>\n        <PH.Title>Vault Usage</PH.Title>\n        <PH.Description>\n          {{#if this.lastUpdatedTime}}\n            Last updated\n            {{this.lastUpdatedTime}}.\n          {{/if}}\n          Don't see what you're looking for?\n          <HdsLinkInline @icon=\"external-link\" @href=\"#\">Share feedback.</HdsLinkInline>\n        </PH.Description>\n        <PH.Actions>\n          <HdsButton @text=\"Refresh\" @icon=\"reload\" @iconPosition=\"leading\" @color=\"secondary\" {{on \"click\" this.fetchAllData}} />\n        </PH.Actions>\n      </HdsPageHeader>\n      {{#if this.data}}\n        <HdsCardContainer @hasBorder={{true}} class=\"dashboard__counters\">\n          {{#each this.counters as |counter|}}\n            <ReportingCounter @title={{counter.title}} @count={{counter.data}} @icon={{counter.icon}} @suffix={{counter.suffix}} @link={{counter.link}} />\n          {{/each}}\n        </HdsCardContainer>\n        <HdsSeparator />\n        <HdsTextDisplay @tag=\"h2\" @size=\"400\" class=\"dashboard__inventory-header\">Resource inventory</HdsTextDisplay>\n        <div class=\"dashboard__viz-blocks\">\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.secret_engines}} @title=\"Secret engines\" @description=\"Breakdown of secret engines for this namespace(s)\" @linkUrl=\"secrets\" class=\"dashboard__viz-block\" />\n            <GlobalLease @count={{this.data.lease_count_quotas.global_lease_count_quota.count}} @quota={{this.data.lease_count_quotas.global_lease_count_quota.capacity}} class=\"dashboard__viz-block\" />\n          </div>\n\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.auth_methods}} @title=\"Authentication methods\" @description=\"Breakdown of authentication methods\" @linkUrl=\"access\" class=\"dashboard__viz-block\" />\n\n            <ClusterReplication @isDisasterRecoveryPrimary={{this.data.replication_status.dr_primary}} @disasterRecoveryState={{this.data.replication_status.dr_state}} @isPerformancePrimary={{this.data.replication_status.pr_primary}} @performanceState={{this.data.replication_status.pr_state}} />\n          </div>\n        </div>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsPageHeader,
        HdsLinkInline,
        HdsButton,
        on,
        HdsCardContainer,
        ReportingCounter: SSUReportingCounter,
        HdsSeparator,
        HdsTextDisplay,
        ReportingHorizontalBarChart: SSUReportingHorizontalBarChart,
        GlobalLease,
        ClusterReplication
      })
    }), this);
  }
}

export { SSUViewDashboard as default };
//# sourceMappingURL=dashboard.js.map
