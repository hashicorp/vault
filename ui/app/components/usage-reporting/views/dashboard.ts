/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type {
  UsageDashboardData,
  NamespaceData,
  SimpleDatum,
  getUsageDataFunction,
  getNamespaceDataFunction,
} from 'vault/types/usage-reporting';

interface CounterBlock {
  title: string;
  tooltipMessage: string;
  data: number;
  icon?: string;
  suffix?: string;
  link?: string;
  emptyText?: string;
  emptyLink?: string;
}

interface VaultReportingViewsDashboardSignature {
  Args: {
    onFetchUsageData?: getUsageDataFunction;
    onFetchNamespaceData?: getNamespaceDataFunction;
    isVaultDedicated?: boolean;
  };
}

export default class VaultReportingViewsDashboard extends Component<VaultReportingViewsDashboardSignature> {
  @tracked data?: UsageDashboardData;
  @tracked namespaceData?: NamespaceData;
  @tracked lastUpdatedTime = '';
  @tracked error?: unknown;
  // Force remount of chart-heavy dashboard layers after namespace data refresh.
  // This avoids stale Carbon chart rendering during namespace switches.
  @tracked chartLayerRenderKey = 0;

  constructor(owner: unknown, args: VaultReportingViewsDashboardSignature['Args']) {
    super(owner, args);
    this.fetchAllData();
  }

  fetchAllData = async (namespace?: string) => {
    try {
      this.error = undefined;
      const usageData = await this.args.onFetchUsageData?.(namespace ?? 'root');
      const namespaceData = await this.args.onFetchNamespaceData?.();

      this.data = usageData;
      this.namespaceData = namespaceData;
      this.chartLayerRenderKey += 1;

      this.lastUpdatedTime = new Intl.DateTimeFormat('en-US', {
        timeStyle: 'medium',
      }).format(new Date());
    } catch (error) {
      this.error = error;
    }
  };

  handleNamespaceChange = (namespace?: string) => {
    this.fetchAllData(namespace);
  };

  getBarChartData = (map: Record<string, number> = {}, exclude?: string[]): SimpleDatum[] => {
    return Object.entries(map)
      .map(([label, value]) => ({ label, value }))
      .filter((item) => !exclude?.includes(item.label));
  };

  get isVaultDedicated() {
    return this.args.isVaultDedicated ?? false;
  }

  get kvSecretsTooltipMessage() {
    const { kvv1Secrets = 0, kvv2Secrets = 0 } = this.data ?? {};
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

  get counters(): CounterBlock[] {
    const { kvv1Secrets = 0, kvv2Secrets = 0 } = this.data ?? {};

    return [
      {
        title: 'Child namespaces',
        tooltipMessage: 'Total number of namespaces for this cluster.',
        data: this.data?.namespaces ?? 0,
        link: 'vault.cluster.access.namespaces',
      },
      {
        title: 'KV secrets',
        tooltipMessage: this.kvSecretsTooltipMessage,
        data: kvv1Secrets + kvv2Secrets,
        emptyText: 'No secrets stored',
        emptyLink: 'vault.cluster.secrets',
      },
      {
        title: 'PKI roles',
        tooltipMessage: 'Total number of PKI roles configured.',
        data: this.data?.pki?.totalRoles ?? 0,
        emptyText: 'No roles created',
      },
    ];
  }

  get namespace() {
    return this.isVaultDedicated ? 'admin' : 'root';
  }
}
