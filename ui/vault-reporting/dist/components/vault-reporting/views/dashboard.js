import Component from '@glimmer/component';
import { concat } from '@ember/helper';
import SSUReportingCounter from '../counter.js';
import SSUReportingHorizontalBarChart from '../horizontal-bar-chart.js';
import GlobalLease from '../global-lease.js';
import ClusterReplication from '../cluster-replication.js';
import DashboardExport from '../dashboard/export.js';
import { tracked } from '@glimmer/tracking';
import { on } from '@ember/modifier';
import { HdsTextDisplay, HdsSeparator, HdsCardContainer, HdsAlert, HdsButton, HdsLinkInline, HdsBadge, HdsPageHeader } from '@hashicorp/design-system-components/components';
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
  #data = (i(this, "data"), void 0);
  static {
    g(this.prototype, "lastUpdatedTime", [tracked], function () {
      return '';
    });
  }
  #lastUpdatedTime = (i(this, "lastUpdatedTime"), void 0);
  static {
    g(this.prototype, "error", [tracked]);
  }
  #error = (i(this, "error"), void 0);
  constructor(owner, args) {
    super(owner, args);
    this.fetchAllData();
  }
  fetchAllData = async () => {
    try {
      this.error = undefined;
      this.data = await this.args.onFetchUsageData();
      this.lastUpdatedTime = new Intl.DateTimeFormat('en-US', {
        dateStyle: 'medium',
        timeStyle: 'medium'
      }).format(new Date());
    } catch (e) {
      this.error = e;
    }
  };
  getBarChartData = (map = {}) => {
    return Object.entries(map).map(([label, value]) => {
      return {
        label,
        value
      };
    });
  };
  get isVaultDedicated() {
    return this.args.isVaultDedicated ?? true;
  }
  get kvSecretsTooltipMessage() {
    const kvv1Secrets = this.data?.kvv1_secrets ?? 0;
    const kvv2Secrets = this.data?.kvv2_secrets ?? 0;
    const kvv1Formatted = Intl.NumberFormat().format(kvv1Secrets);
    const kvv2Formatted = Intl.NumberFormat().format(kvv2Secrets);
    if (kvv1Secrets && kvv2Secrets) {
      return `Combined count of ${kvv1Formatted} KV v1 secrets and ${kvv2Formatted} KV v2 secrets in this namespace`;
    }
    if (kvv1Secrets) {
      return `Total number of ${kvv1Formatted} KV v1 secrets in this namespace`;
    }
    if (kvv2Secrets) {
      return `Total number of ${kvv2Formatted} KV v2 secrets in this namespace`;
    }
    return '';
  }
  get counters() {
    return [{
      title: 'Child namespaces',
      tooltipMessage: this.isVaultDedicated ? 'Total number of direct child namespaces under the root/ namespace' : 'Total number of direct child namespaces under the admin/ namespace',
      data: this.data?.namespaces ?? 0,
      link: 'access/namespaces'
    }, {
      title: 'KV secrets',
      tooltipMessage: this.kvSecretsTooltipMessage,
      data: (this.data?.kvv1_secrets ?? 0) + (this.data?.kvv2_secrets ?? 0),
      emptyText: 'No secrets stored'
    }, {
      title: 'Secrets sync',
      tooltipMessage: this.isVaultDedicated ? 'Total number of destinations (e.g. third-party integrations) synced with secrets from this namespace' : '',
      data: this.data?.secrets_sync ?? 0,
      link: 'sync/secrets/overview',
      emptyText: 'Not configured',
      suffix: 'destinations'
    }, {
      title: 'PKI roles',
      tooltipMessage: this.isVaultDedicated ? 'Total number of PKI roles configured in this namespace' : '',
      data: this.data?.pki?.total_roles ?? 0,
      emptyText: 'No roles created'
    }];
  }
  get namespace() {
    return this.isVaultDedicated ? 'Admin' : 'Root';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div class=\"dashboard\">\n      <HdsPageHeader as |PH|>\n        <PH.Title>Vault Usage\n          <HdsBadge class=\"dashboard__badge\" @size=\"small\" @type=\"inverted\" @text={{concat \"Namespace: \" this.namespace}} />\n        </PH.Title>\n        <PH.Description>\n          {{#if this.lastUpdatedTime}}\n            Last updated\n            {{this.lastUpdatedTime}}.\n          {{/if}}\n          Don't see what you're looking for?\n          <HdsLinkInline @icon=\"external-link\" @href=\"#\">Share feedback.</HdsLinkInline>\n        </PH.Description>\n        <PH.Actions>\n          <HdsButton @text=\"Refresh\" @icon=\"reload\" @iconPosition=\"leading\" @color=\"secondary\" data-test-dashboard-refresh-button {{on \"click\" this.fetchAllData}} />\n          <DashboardExport @data={{this.data}} />\n        </PH.Actions>\n      </HdsPageHeader>\n      {{#if this.error}}\n        <HdsAlert data-test-dashboard-error @type=\"inline\" @color=\"critical\" class=\"dashboard__error\" as |A|>\n          <A.Title>Error</A.Title>\n          <A.Description data-test-dashboard-error-description>An error\n            occurred, please try again.</A.Description>\n        </HdsAlert>\n      {{/if}}\n      {{#if this.data}}\n        <HdsCardContainer @hasBorder={{true}} class=\"dashboard__counters\" data-test-dashboard-counters>\n          {{#each this.counters as |counter|}}\n            <ReportingCounter @title={{counter.title}} @tooltipMessage={{counter.tooltipMessage}} @count={{counter.data}} @icon={{counter.icon}} @suffix={{counter.suffix}} @link={{counter.link}} @emptyText={{counter.emptyText}} />\n          {{/each}}\n        </HdsCardContainer>\n        <HdsSeparator />\n        <HdsTextDisplay @tag=\"h2\" @size=\"400\" class=\"dashboard__inventory-header\">Resource inventory</HdsTextDisplay>\n        <div class=\"dashboard__viz-blocks\">\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.secret_engines}} @title=\"Secret engines\" @description=\"Breakdown of secret engines for this namespace(s)\" @linkUrl=\"secrets\" class=\"dashboard__viz-block\" data-test-dashboard-secret-engines />\n            <GlobalLease @count={{this.data.lease_count_quotas.global_lease_count_quota.count}} @quota={{this.data.lease_count_quotas.global_lease_count_quota.capacity}} class=\"dashboard__viz-block\" data-test-dashboard-lease-count />\n          </div>\n\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.auth_methods}} @title=\"Authentication methods\" @description=\"Breakdown of authentication methods\" @linkUrl=\"access\" class=\"dashboard__viz-block\" data-test-dashboard-auth-methods />\n\n            <ClusterReplication @isDisasterRecoveryPrimary={{this.data.replication_status.dr_primary}} @disasterRecoveryState={{this.data.replication_status.dr_state}} @isPerformancePrimary={{this.data.replication_status.pr_primary}} @performanceState={{this.data.replication_status.pr_state}} data-test-dashboard-cluster-replication />\n          </div>\n        </div>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsPageHeader,
        HdsBadge,
        concat,
        HdsLinkInline,
        HdsButton,
        on,
        DashboardExport,
        HdsAlert,
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
