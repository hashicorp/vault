/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type { UsageDashboardData } from 'vault/types/usage-reporting';
import { toSentenceCase } from 'vault/utils/to-sentence-case';

interface VaultReportingDashboardExportSignature {
  Args: {
    data?: UsageDashboardData;
  };
}

export default class VaultReportingDashboardExport extends Component<VaultReportingDashboardExportSignature> {
  private getNestedRows(records: Record<string, number>, prefix = '') {
    return Object.entries(records).map(([key, value]) => [
      `${prefix} ${toSentenceCase(key, { acronymsOnly: true })}`,
      value,
    ]);
  }

  get dataAsDownloadableJSONString() {
    const jsonString = JSON.stringify(this.args.data, null, '    ');
    return `data:application/json;charset=utf-8,${encodeURIComponent(jsonString)}`;
  }

  get dataAsDownloadableCSVString() {
    const headers = ['Metric', 'Count/breakdown'];

    const rows = [
      headers,
      ['Child namespaces', this.args.data?.namespaces || 0],
      ['Total KV secrets', (this.args.data?.kvv1Secrets || 0) + (this.args.data?.kvv2Secrets || 0)],
      ['KV V1 secrets', this.args.data?.kvv1Secrets || 0],
      ['KV V2 secrets', this.args.data?.kvv2Secrets || 0],
      ['PKI roles', this.args.data?.pki?.totalRoles || 0],
      ...this.getNestedRows(this.args.data?.secretEngines || {}, 'Secret engine'),
      ...this.getNestedRows(this.args.data?.authMethods || {}, 'Auth method'),
      ['Global lease count', this.args.data?.leaseCountQuotas.globalLeaseCountQuota.count || 0],
      ['Global lease quota', this.args.data?.leaseCountQuotas.globalLeaseCountQuota.capacity || 0],
      ['Cluster disaster recovery', this.args.data?.replicationStatus.drState || '-'],
      ['Cluster disaster recovery primary', this.args.data?.replicationStatus.drPrimary ?? '-'],
      ['Cluster performance', this.args.data?.replicationStatus.prState || '-'],
      ['Cluster performance primary', this.args.data?.replicationStatus.prPrimary ?? '-'],
      ['Secrets sync', this.args.data?.secretSync?.totalDestinations || 0],
      ...this.getNestedRows(this.args.data?.secretSync.destinations || {}, 'Secrets sync destination'),
    ];

    const csvString = rows
      .map((row) =>
        row
          .map((cell) => {
            const escaped = String(cell).replace(/"/g, '""');
            return `"${escaped}"`;
          })
          .join(',')
      )
      .join('\r\n');

    return `data:text/csv;charset=utf-8,${encodeURIComponent(csvString)}`;
  }
}
