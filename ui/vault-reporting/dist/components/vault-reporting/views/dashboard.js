import Component from '@glimmer/component';
import { array } from '@ember/helper';
import SSUReportingCounter from '../counter.js';
import SSUReportingHorizontalBarChart from '../horizontal-bar-chart.js';
import GlobalLease from '../global-lease.js';
import ClusterReplication from '../cluster-replication.js';
import DashboardExport from '../dashboard/export.js';
import { tracked } from '@glimmer/tracking';
import { HdsCardContainer, HdsAlert, HdsLinkInline, HdsTextBody, HdsBadge, HdsPageHeader } from '@hashicorp/design-system-components/components';
import { service } from '@ember/service';
import { on } from '@ember/modifier';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUViewDashboard extends Component {
  static {
    g(this.prototype, "reportingAnalytics", [service]);
  }
  #reportingAnalytics = (i(this, "reportingAnalytics"), void 0);
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
        timeStyle: 'medium'
      }).format(new Date());
    } catch (e) {
      this.error = e;
    }
  };
  handleTrackAnalyticsEvent = (eventName, properties, options) => {
    this.reportingAnalytics.trackEvent(eventName, properties, options);
  };
  handleTrackSurveyLink = () => {
    this.handleTrackAnalyticsEvent('survey_link');
  };
  handleRefresh = () => {
    this.fetchAllData();
    this.handleTrackAnalyticsEvent('refresh_button');
  };
  getBarChartData = (map = {}, exclude) => {
    return Object.entries(map).map(([label, value]) => {
      return {
        label,
        value
      };
    }).filter(item => {
      return !exclude?.includes(item.label);
    });
  };
  get isVaultDedicated() {
    return this.args.isVaultDedicated ?? false;
  }
  get kvSecretsTooltipMessage() {
    const {
      kvv1Secrets = 0,
      kvv2Secrets = 0
    } = this.data ?? {};
    const kvv1Formatted = Intl.NumberFormat().format(kvv1Secrets);
    const kvv2Formatted = Intl.NumberFormat().format(kvv2Secrets);
    if (kvv1Secrets && kvv2Secrets) {
      return `Combined count of ${kvv1Formatted} KV version 1 secrets and ${kvv2Formatted} KV version 2 secrets.`;
    }
    if (kvv1Secrets) {
      return `Total number of ${kvv1Formatted} KV version 1 secrets.`;
    }
    if (kvv2Secrets) {
      return `Total number of ${kvv2Formatted} KV version 2 secrets.`;
    }
    return '';
  }
  get counters() {
    const {
      kvv1Secrets = 0,
      kvv2Secrets = 0
    } = this.data ?? {};
    return [{
      title: 'Child namespaces',
      tooltipMessage: 'Total number of namespaces for this cluster.',
      data: this.data?.namespaces ?? 0,
      link: 'access/namespaces'
    }, {
      title: 'KV secrets',
      tooltipMessage: this.kvSecretsTooltipMessage,
      data: kvv1Secrets + kvv2Secrets,
      emptyText: 'No secrets stored',
      emptyLink: 'secrets'
    }, {
      title: 'Secrets sync',
      tooltipMessage: 'Total number of destinations (e.g. third-party integrations) synced with secrets from this namespace.',
      data: this.data?.secretSync?.totalDestinations ?? 0,
      link: 'sync/secrets/overview',
      emptyText: 'Not activated',
      suffix: 'destinations'
    }, {
      title: 'PKI roles',
      tooltipMessage: 'Total number of PKI roles configured.',
      data: this.data?.pki?.totalRoles ?? 0,
      emptyText: 'No roles created'
    }];
  }
  get namespace() {
    return this.isVaultDedicated ? 'admin' : 'root';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div class=\"dashboard\" data-test-dashboard-container>\n      <HdsPageHeader as |PH|>\n        <PH.Title>\n          Vault Usage\n          <HdsBadge class=\"dashboard__badge\" @size=\"medium\" @icon=\"org\" @color=\"neutral\" @text={{this.namespace}} />\n        </PH.Title>\n        <PH.Description class=\"dashboard__description\">\n          {{#if this.lastUpdatedTime}}\n            <HdsTextBody @tag=\"p\" @size=\"200\" @color=\"faint\">\n              Updated today at\n              {{this.lastUpdatedTime}}.\n\n            </HdsTextBody>\n          {{/if}}\n          <HdsTextBody @tag=\"p\" @size=\"200\" @color=\"primary\">\n            View and export your Vault usage. Don't see what you're looking for?\n            <HdsLinkInline data-test-vault-reporting-dashboard-survey-link @icon=\"external-link\" @href=\"https://hashicorp.sjc1.qualtrics.com/jfe/form/SV_bqhLeB3deLd2caa\" target=\"_blank\" {{on \"click\" this.handleTrackSurveyLink}}>Share feedback</HdsLinkInline>\n          </HdsTextBody>\n        </PH.Description>\n        <PH.Actions>\n          <DashboardExport @data={{this.data}} />\n        </PH.Actions>\n      </HdsPageHeader>\n      {{#if this.error}}\n        <HdsAlert data-test-vault-reporting-dashboard-error @type=\"inline\" @color=\"critical\" class=\"dashboard__error\" as |A|>\n          <A.Title>Error</A.Title>\n          <A.Description data-test-vault-reporting-dashboard-error-description>An error occurred, please try again.</A.Description>\n        </HdsAlert>\n      {{/if}}\n      {{#if this.data}}\n        <HdsCardContainer @hasBorder={{true}} {{!-- @glint-expect-error --}} @background=\"neutral-secondary\" class=\"dashboard__counters\" data-test-vault-reporting-dashboard-counters>\n          {{#each this.counters as |counter|}}\n            <ReportingCounter @title={{counter.title}} @tooltipMessage={{counter.tooltipMessage}} @count={{counter.data}} @icon={{counter.icon}} @suffix={{counter.suffix}} @link={{counter.link}} @emptyText={{counter.emptyText}} @emptyLink={{counter.emptyLink}} />\n          {{/each}}\n        </HdsCardContainer>\n        <div data-test-vault-reporting-dashboard-viz-blocks class=\"dashboard__viz-blocks\">\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.authMethods}} @title=\"Authentication methods\" @description=\"Enabled authentication methods for this cluster.\" @linkUrl=\"access\" class=\"dashboard__viz-block\" data-test-vault-reporting-dashboard-auth-methods />\n\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.secretEngines (array \"system\" \"identity\")}} @title=\"Secret engines\" @description=\"Enabled secret engines for this cluster.\" @linkUrl=\"secrets\" class=\"dashboard__viz-block\" data-test-vault-reporting-dashboard-secret-engines />\n\n            <ClusterReplication @disasterRecoveryState={{this.data.replicationStatus.drState}} @performanceState={{this.data.replicationStatus.prState}} @isVaultDedicated={{this.isVaultDedicated}} data-test-vault-reporting-dashboard-cluster-replication />\n          </div>\n\n          <div>\n            <ReportingHorizontalBarChart @data={{this.getBarChartData this.data.leasesByAuthMethod}} @title=\"Leases by authentication methods\" @description=\"Active leases issued per authentication method.\" @linkUrl=\"https://developer.hashicorp.com/vault/docs/concepts/auth#auth-leases\" @linkText=\"Documentation\" @linkIcon=\"docs-link\" @linkTarget=\"_blank\" class=\"dashboard__viz-block\" data-test-vault-reporting-dashboard-leases-by-auth-method>\n              <:empty as |A|>\n                <A.Body @text=\"Lease are created when clients authenticate. Add an authentication method to monitor leases across this namespace.\" />\n                <A.Footer as |F|>\n                  <F.LinkStandalone @icon=\"docs-link\" @text=\"Authentication leases\" @href=\"https://developer.hashicorp.com/vault/docs/concepts/auth#auth-leases\" />\n                </A.Footer>\n              </:empty>\n            </ReportingHorizontalBarChart>\n            <GlobalLease @count={{this.data.leaseCountQuotas.globalLeaseCountQuota.count}} @quota={{this.data.leaseCountQuotas.globalLeaseCountQuota.capacity}} class=\"dashboard__viz-block\" data-test-vault-reporting-dashboard-lease-count />\n          </div>\n\n        </div>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsPageHeader,
        HdsBadge,
        HdsTextBody,
        HdsLinkInline,
        on,
        DashboardExport,
        HdsAlert,
        HdsCardContainer,
        ReportingCounter: SSUReportingCounter,
        ReportingHorizontalBarChart: SSUReportingHorizontalBarChart,
        array,
        ClusterReplication,
        GlobalLease
      })
    }), this);
  }
}

export { SSUViewDashboard as default };
//# sourceMappingURL=dashboard.js.map
