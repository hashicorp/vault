/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';
import { calculateSum } from 'vault/utils/chart-helpers';
import { NormalizedBillingMetrics } from 'vault/utils/metrics-helpers';

interface Args {
  title: string;
  metrics: Record<string, number>;
}

export default class MetricCard extends Component<Args> {
  get total() {
    const sums = Object.values(this.args.metrics).filter((metric) => metric !== undefined);
    return calculateSum(sums);
  }

  get description() {
    switch (this.args.title) {
      case 'Secrets':
        return 'Highest number of static secrets, static roles, and dynamic roles managed on the cluster during the month. Secrets replicated to this cluster are not counted.';
      case 'Credential units':
        return 'Certificates, tokens, and other credentials issued during the month, adjusted by their duration.';
      case 'Data protection calls':
        return 'Total number of data elements processed.';
      case 'Managed keys':
        return 'Highest number of cryptographic keys managed on the cluster during the month. Keys replicated to this cluster are not counted.';
      default:
        return '';
    }
  }

  metricDetailsMap: Record<string, { label: string; tooltipText?: string }> = {
    [NormalizedBillingMetrics.STATIC_SECRETS_KV]: {
      label: 'KV Secrets',
    },
    [NormalizedBillingMetrics.DYNAMIC_ROLES]: {
      label: 'Dynamic roles',
      tooltipText: 'Highest number of dynamic roles for the month',
    },
    [NormalizedBillingMetrics.STATIC_ROLES]: {
      label: 'Static roles',
      tooltipText: 'Highest number of static roles for the month',
    },
    [NormalizedBillingMetrics.PKI_UNITS_TOTAL]: {
      label: 'PKI units',
      tooltipText: 'Total number of X.509 certificates issued, normalized by their duration.',
    },
    [NormalizedBillingMetrics.SSH_UNITS_OTP_UNITS]: {
      label: 'SSH OTP units',
      tooltipText:
        'Total number of SSH one-time passwords issued, normalized by their duration. Each OTP is 0.0014 units.',
    },
    [NormalizedBillingMetrics.SSH_UNITS_CERTIFICATE_UNITS]: {
      label: 'SSH certificate units',
      tooltipText: 'Total number of SSH certificates issued, normalized by their duration.',
    },
    [NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT]: {
      label: 'Transit',
    },
    [NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSFORM]: {
      label: 'Transform',
    },
    [NormalizedBillingMetrics.MANAGED_KEYS_TOTP]: {
      label: 'TOTP',
    },
    [NormalizedBillingMetrics.MANAGED_KEYS_KMSE]: {
      label: 'KMSE',
    },
  };

  metricDetails = (key: string): { label: string; tooltipText?: string; count?: number } => {
    return this.metricDetailsMap[key] || { label: toLabel([key]) };
  };
}
