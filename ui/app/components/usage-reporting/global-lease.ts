/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { htmlSafe } from '@ember/template';
import type { MeterChartOptions } from '@carbon/charts/dist/interfaces';

import { CHART_TYPES } from 'vault/modifiers/carbon-chart';

interface MeterDatum {
  group: string;
  value: number;
}

interface VaultReportingGlobalLeaseSignature {
  Args: {
    count?: number;
    quota?: number;
  };
}
const CHART_HEIGHT = '26px';

export default class VaultReportingGlobalLease extends Component<VaultReportingGlobalLeaseSignature> {
  chartType = CHART_TYPES.METER;

  get percentage() {
    const { count = 0, quota = 0 } = this.args;
    if (!quota) {
      return 0;
    }

    return Math.round(Math.min((count / quota) * 100, 100));
  }

  get chartData(): MeterDatum[] {
    return [{ group: 'Leases', value: this.percentage }];
  }

  get chartOptions(): MeterChartOptions {
    return {
      meter: {
        statusBar: {
          percentageIndicator: {
            enabled: false,
          },
        },
        status: {
          ranges: [
            {
              range: [0, 94],
              status: 'success',
            },
            {
              range: [95, 99],
              status: 'warning',
            },
            {
              range: [100, 100],
              status: 'danger',
            },
          ],
        },
      },
      toolbar: {
        enabled: false,
      },
      legend: {
        enabled: false,
      },
      height: CHART_HEIGHT,
      accessibility: {
        svgAriaLabel: 'Global lease count quota',
      },
    } as MeterChartOptions;
  }

  get formattedCount() {
    const formatter = new Intl.NumberFormat('en-US', {
      notation: 'compact',
      compactDisplay: 'short',
    });

    const { count = 0, quota = 0 } = this.args;
    const formattedCount = formatter.format(count);
    const formattedTotal = formatter.format(quota);

    return `${formattedCount} / ${formattedTotal}`;
  }

  get hasData() {
    return this.args.quota && typeof this.args.quota === 'number';
  }

  get description() {
    if (this.hasData) {
      return htmlSafe(
        'Total number of active <a class="hds-link-inline--color-secondary" href="https://developer.hashicorp.com/vault/docs/concepts/lease" target="_blank" rel="noopener noreferrer" data-test-vault-reporting-global-lease-description-link>leases</a> for this quota.'
      );
    }
    return undefined;
  }

  get linkUrl() {
    if (this.hasData) {
      return 'https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota';
    }
    return undefined;
  }

  get alert(): { color: 'warning' | 'neutral'; description: string } | undefined {
    if (this.percentage >= 100) {
      return {
        color: 'warning',
        description:
          'Global lease quota limit reached. If lease creation is blocked, reduce usage or increase the limit.',
      };
    }

    if (this.percentage >= 95) {
      return {
        color: 'neutral',
        description:
          'Approaching quota limit. Reduce usage or increase the lease limit to avoid blocking new leases.',
      };
    }

    return;
  }
}
