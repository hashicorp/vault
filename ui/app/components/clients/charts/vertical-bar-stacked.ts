/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { BAR_WIDTH, numericalAxisLabel } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { flatGroup } from 'd3-array';
import type { MonthlyChartData } from 'vault/vault/charts/client-counts';
import type { ClientTypes } from 'core/utils/client-count-utils';

interface Args {
  chartHeight?: number;
  chartLegend: Legend[];
  chartTitle: string;
  data: MonthlyChartData[];
}

interface Legend {
  key: ClientTypes;
  label: string;
}
interface AggregatedDatum {
  x: string;
  y: number;
  legendX: string;
  legendY: string[];
}

interface Base {
  timestamp: string;
  clientType: string;
}

type KeyDataItem = Base & {
  [key in ClientTypes]: number | undefined;
};

/**
 * @module VerticalBarStacked
 * Renders a vertical bar chart of counts fora single data point (@dataKey) over time.
 *
 * @example
 * <Clients::Charts::VerticalBarStacked
 * @chartTitle="Total monthly usage"
 * @data={{this.byMonthActivityData}}
 * @chartLegend={{this.legend}}
 * @chartHeight={{250}}
 * />
 */
export default class VerticalBarStacked extends Component<Args> {
  barWidth = BAR_WIDTH;
  @tracked activeDatum: AggregatedDatum | null = null;

  get chartHeight() {
    return this.args.chartHeight || 190;
  }

  get dataKeys(): ClientTypes[] {
    return this.args.chartLegend.map((l: Legend) => l.key);
  }

  label(legendKey: string) {
    return this.args.chartLegend.find((l: Legend) => l.key === legendKey).label;
  }

  get chartData() {
    let dataset: [string, string, number | undefined, KeyDataItem[]][] = [];
    // each datum needs to be its own object
    for (const key of this.dataKeys) {
      const keyData: KeyDataItem[] = this.args.data.map((d: MonthlyChartData) => ({
        timestamp: d.timestamp,
        clientType: key,
        [key]: d[key],
      }));

      const group = flatGroup(
        keyData,
        // order here must match destructure order in return below
        (d) => d.timestamp,
        (d) => d.clientType,
        (d) => d[key]
      );
      dataset = [...dataset, ...group];
    }

    return dataset.map(([timestamp, clientType, counts]) => ({
      timestamp,
      clientType, // key name matches the chart's @color arg
      counts,
    }));
  }

  // for yRange scale, tooltip target area and tooltip text data
  get aggregatedData(): AggregatedDatum[] {
    return this.args.data.map((datum: MonthlyChartData) => {
      const values = this.dataKeys
        .map((k: string) => datum[k as ClientTypes])
        .filter((count) => Number.isInteger(count));
      const sum = values.length ? values.reduce((sum, currentValue) => sum + currentValue, 0) : null;
      const xValue = datum.timestamp;
      return {
        x: xValue,
        y: sum ?? 0,
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY: sum ? this.dataKeys.map((k) => `${formatNumber([datum[k]])} ${this.label(k)}`) : ['No data'],
      };
    });
  }

  get yRange() {
    const counts: number[] = this.aggregatedData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=4, hardcode 4 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 4 ? 4 : max];
  }

  get xDomain() {
    const domain = this.chartData.map((d) => d.timestamp);
    return new Set(domain);
  }

  // TEMPLATE HELPERS
  barOffset = (bandwidth: number) => (bandwidth - this.barWidth) / 2;

  tooltipX = (original: number, bandwidth: number) => (original + bandwidth / 2).toString();

  tooltipY = (original: number) => (!original ? '0' : `${original}`);

  formatTicksX = (timestamp: string): string => parseAPITimestamp(timestamp, 'M/yy');

  formatTicksY = (num: number): string => numericalAxisLabel(num) || num.toString();
}
