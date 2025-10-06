/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { BAR_WIDTH, numericalAxisLabel } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { parseAPITimestamp } from 'core/utils/date-formatters';

import type { MonthlyChartData } from 'vault/vault/client-counts/charts';
import type { TotalClients } from 'vault/vault/client-counts/activity-api';

interface Args {
  data: MonthlyChartData[];
  dataKey: string;
  chartTitle: string;
  chartHeight?: number;
}

interface ChartData {
  x: string;
  y: number | null;
  tooltip: string;
  legendX: string;
  legendY: string;
}

/**
 * @module VerticalBarBasic
 * Renders a vertical bar chart of counts fora single data point (@dataKey) over time.
 *
 * @example
 <Clients::Charts::VerticalBarBasic
    @chartTitle="Secret Sync client counts"
    @data={{this.model}}
    @dataKey="secret_syncs"
    @showTable={{true}}
    @chartHeight={{200}}
  />
 */
export default class VerticalBarBasic extends Component<Args> {
  barWidth = BAR_WIDTH;

  @tracked activeDatum: ChartData | null = null;

  get chartHeight() {
    return this.args.chartHeight || 190;
  }

  get chartData() {
    return this.args.data.map((d): ChartData => {
      const xValue = d.timestamp as string;
      const yValue = (d[this.args.dataKey as keyof TotalClients] as number) ?? null;
      return {
        x: parseAPITimestamp(xValue, 'M/yy') as string,
        y: yValue,
        tooltip:
          yValue === null ? 'No data' : `${formatNumber([yValue])} ${this.args.dataKey.replace(/_/g, ' ')}`,
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY: (yValue ?? 'No data').toString(),
      };
    });
  }

  get yDomain() {
    const counts: number[] = this.chartData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=4, hardcode 4 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 4 ? 4 : max];
  }

  get xDomain() {
    const months = this.chartData.map((d) => d.x);
    return new Set(months);
  }

  // TEMPLATE HELPERS
  barOffset = (bandwidth: number) => {
    return (bandwidth - this.barWidth) / 2;
  };

  tooltipX = (original: number, bandwidth: number) => {
    return (original + bandwidth / 2).toString();
  };

  tooltipY = (original: number) => {
    if (!original) return `0`;
    return `${original}`;
  };

  formatTicksY = (num: number): string => {
    return numericalAxisLabel(num) || num.toString();
  };
}
