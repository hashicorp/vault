/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { formatNumbers, formatTooltipNumber } from 'vault/utils/chart-helpers';
import type { SerializedChartData } from 'vault/client-counts';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format } from 'date-fns';

interface Args {
  data: SerializedChartData[];
  dataKey: string;
  chartTitle: string;
}

interface ChartData {
  x: string;
  y: number | null;
  tooltip: string;
  legendX: string;
  legendY: string;
}

/**
 * @module ClientsSyncBarChartComponent
 * Renders a bar chart of secret syncs over time.
 *
 * @example
 * ```js
 * <Clients::SyncBarChart
    @chartTitle="Secret Sync client counts"
    @data={{this.model}}
    @dataKey="secret_syncs"
    @showTable={{true}} />
 * ```
 */
export default class ClientsSyncBarChartComponent extends Component<Args> {
  barWidth = 7;
  chartHeight = 190;
  chartWidth = 600;

  @tracked activeDatum: ChartData | null = null;

  get chartData() {
    return this.args.data.map((d): ChartData => {
      const date = parseAPITimestamp(d.timestamp) as Date;
      const count = (d[this.args.dataKey] as number) ?? null;
      return {
        x: format(date, 'M/yy'),
        y: count,
        tooltip: count === null ? 'No data' : `${formatTooltipNumber(count)} secret syncs`,
        legendX: format(date, 'MMM yyyy'),
        legendY: (count ?? 'No data').toString(),
      };
    });
  }

  get countDomain() {
    const counts: number[] = this.chartData.map((d) => d.y).flatMap((num) => (num ? [num] : []));
    const upper = Math.round(Math.max(...counts) / 1000) * 1000;
    return [0, upper];
  }

  get monthDomain() {
    const months = this.chartData.map((d) => d.x);
    return new Set(months);
  }

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
  formatCount = (num: number): string => {
    return formatNumbers(num) || num.toString();
  };
}
