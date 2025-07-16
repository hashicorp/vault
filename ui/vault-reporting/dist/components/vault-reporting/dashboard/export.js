import Component from '@glimmer/component';
import { HdsDropdown } from '@hashicorp/design-system-components/components';
import { on } from '@ember/modifier';
import { fn } from '@ember/helper';
import { service } from '@ember/service';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class DashboardExport extends Component {
  static {
    g(this.prototype, "reportingAnalytics", [service]);
  }
  #reportingAnalytics = (i(this, "reportingAnalytics"), void 0);
  #getNestedRows(records, prefix = '') {
    return Object.entries(records).map(([key, value]) => {
      return [`${prefix} ${key}`, value];
    });
  }
  handleTrackExportToggle = () => {
    this.reportingAnalytics.trackEvent('export_toggle');
  };
  handleTrackExportOption = option => {
    this.reportingAnalytics.trackEvent(`export_option`, {
      option
    });
  };
  get dataAsDownloadableJSONString() {
    const {
      data
    } = this.args;
    const file = new Blob([JSON.stringify(data, null, '    ')], {
      type: 'application/json'
    });
    return URL.createObjectURL(file);
  }
  get dataAsDownloadableCSVString() {
    const headers = ['Metric', 'Count/Breakdown'];
    // Manually define rows as looping through the data does not leave the most legible structure
    const rows = [headers, ['Child Namespaces', this.args?.data?.namespaces || 0], ['Total KV Secrets', (this.args.data?.kvv1Secrets || 0) + (this.args.data?.kvv2Secrets || 0)], ['KV V1 Secrets', this.args.data?.kvv1Secrets || 0], ['KV V2 Secrets', this.args.data?.kvv2Secrets || 0], ['Secret Syncs', this.args.data?.secretSync?.totalDestinations || 0], ['PKI Roles', this.args.data?.pki?.totalRoles || 0], ...this.#getNestedRows(this.args.data?.secretEngines || {}, 'Secret Engine'), ...this.#getNestedRows(this.args.data?.authMethods || {}, 'Auth Method'), ['Global Lease Count', this.args.data?.leaseCountQuotas.globalLeaseCountQuota.count || 0], ['Global Lease Quota', this.args.data?.leaseCountQuotas.globalLeaseCountQuota.capacity || 0], ['Cluster Disaster Recovery', this.args?.data?.replicationStatus.drState || '-'], ['Cluster Disaster Recovery Primary', this.args?.data?.replicationStatus.drPrimary ?? '-'], ['Cluster Performance', this.args?.data?.replicationStatus.prState || '-'], ['Cluster Performance Primary', this.args?.data?.replicationStatus.prPrimary ?? '-']];
    // Escape double quotes, quote cell content and separate with comma
    const csvString = rows.map(row => row.map(cell => {
      const escaped = String(cell).replace(/"/g, '""');
      return `"${escaped}"`;
    }).join(',')).join('\r\n');
    const blob = new Blob([csvString], {
      type: 'text/csv'
    });
    return URL.createObjectURL(blob);
  }
  static {
    setComponentTemplate(precompileTemplate("\n    {{#if @data}}\n      <HdsDropdown @matchToggleWidth={{true}} as |D|>\n        <D.ToggleButton data-test-vault-reporting-export-toggle @text=\"Export\" {{on \"click\" this.handleTrackExportToggle}} />\n        <D.Interactive data-test-vault-reporting-export-json @href={{this.dataAsDownloadableJSONString}} download=\"vault-usage-dashboard.json\" {{on \"click\" (fn this.handleTrackExportOption \"json\")}}>JSON</D.Interactive>\n        <D.Interactive data-test-vault-reporting-export-csv @href={{this.dataAsDownloadableCSVString}} download=\"vault-usage-dashboard.csv\" {{on \"click\" (fn this.handleTrackExportOption \"csv\")}}>CSV</D.Interactive>\n      </HdsDropdown>\n    {{/if}}\n  ", {
      strictMode: true,
      scope: () => ({
        HdsDropdown,
        on,
        fn
      })
    }), this);
  }
}

export { DashboardExport as default };
//# sourceMappingURL=export.js.map
