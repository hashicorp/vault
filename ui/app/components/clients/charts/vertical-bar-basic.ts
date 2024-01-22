/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { BAR_WIDTH, formatNumbers } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import type { Count, MonthlyChartData } from 'vault/vault/charts/client-counts';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format } from 'date-fns';
import { toLabel } from 'core/helpers/to-label';

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
 * Renders a bar chart of secret syncs over time.
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
      const date = parseAPITimestamp(d.timestamp) as Date;
      const count = d[this.args.dataKey as keyof Count] ?? null;
      return {
        x: format(date, 'M/yy'),
        y: count,
        tooltip: count === null ? 'No data' : `${formatNumber([count])} ${toLabel([this.args.dataKey])}`,
        legendX: format(date, 'MMMM yyyy'),
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
