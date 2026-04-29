/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';
import { NormalizedMetricsData } from 'vault/vault/billing/overview';
import { NormalizedBillingMetrics } from 'vault/utils/metrics-helpers';

interface SummaryMetricInfo {
  label: string;
  tooltipText?: string;
  count?: number;
  showBadge?: boolean;
  total?: number | boolean | undefined;
}
interface Args {
  title: string;
  metrics: Record<string, number>;
  normalizedMetricData: NormalizedMetricsData;
}

export default class SummaryCard extends Component<Args> {
  summaryMetricKeys = [
    NormalizedBillingMetrics.STATIC_SECRETS_TOTAL,
    NormalizedBillingMetrics.CREDENTIAL_UNITS_TOTAL,
    NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TOTAL,
    NormalizedBillingMetrics.MANAGED_KEYS_TOTAL,
    NormalizedBillingMetrics.KMIP_USED_IN_MONTH,
    NormalizedBillingMetrics.EXTERNAL_PLUGINS_TOTAL,
  ];

  summaryMetricMap: Record<string, SummaryMetricInfo> = {
    [NormalizedBillingMetrics.STATIC_SECRETS_TOTAL]: {
      label: 'Secrets',
    },
    [NormalizedBillingMetrics.CREDENTIAL_UNITS_TOTAL]: {
      label: 'Credential units',
    },
    [NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TOTAL]: {
      label: 'Data protection calls',
    },
    [NormalizedBillingMetrics.MANAGED_KEYS_TOTAL]: {
      label: 'Managed keys',
    },
    [NormalizedBillingMetrics.KMIP_USED_IN_MONTH]: {
      label: 'KMIP',
      tooltipText: 'Whether KMIP was enabled on the cluster at any time during the month.',
      showBadge: true,
    },
    [NormalizedBillingMetrics.EXTERNAL_PLUGINS_TOTAL]: {
      label: 'Plugins',
      tooltipText: 'Highest number of plugins enabled on the cluster at any time during the month.',
    },
  };

  summaryMetric = (key: string): SummaryMetricInfo => {
    if (this.summaryMetricMap?.[key]) {
      this.summaryMetricMap[key].total = this.args.normalizedMetricData[key];
    }

    return this.summaryMetricMap[key] || { label: toLabel([key]) };
  };
}
