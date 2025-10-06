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

import type { MonthlyChartData } from 'vault/vault/client-counts/charts';
import type { ClientTypes } from 'vault/vault/client-counts/activity-api';

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

interface DatumBase {
  timestamp: string;
  clientType: string;
}
// separated because "A mapped type may not declare properties or methods."
type ChartDatum = DatumBase & {
  [key in ClientTypes]?: number | undefined;
};

/**
 * @module VerticalBarStacked
 * Renders a stacked bar chart of counts for different client types over time. Which client types render
 * is mapped from the "key" values of the @legend arg
 *
 * @example
 * <Clients::Charts::VerticalBarStacked
 * @chartTitle="Total monthly usage"
 * @data={{this.byMonthNewClients}}
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
    return this.args.chartLegend.find((l: Legend) => l.key === legendKey)?.label;
  }

  get chartData() {
    let dataset: [string, number | undefined, string, ChartDatum[]][] = [];
    // each datum needs to be its own object
    for (const key of this.dataKeys) {
      const chartData: ChartDatum[] = this.args.data.map((d: MonthlyChartData) => ({
        timestamp: d.timestamp,
        clientType: key,
        [key]: d[key],
      }));

      const group = flatGroup(
        chartData,
        // order here must match destructure order in return below
        (d) => d.timestamp,
        (d) => d[key],
        (d) => d.clientType
      );
      dataset = [...dataset, ...group];
    }

    return dataset.map(([timestamp, counts, clientType]) => ({
      timestamp, // x value
      counts, // y value
      clientType, // corresponds to chart's @color arg
    }));
  }

  // for yBounds scale, tooltip target area and tooltip text data
  get aggregatedData(): AggregatedDatum[] {
    return this.args.data.map((datum: MonthlyChartData) => {
      const values = this.dataKeys
        .map((k: string) => datum[k as ClientTypes])
        .filter((count) => Number.isInteger(count));
      const sum = values.length ? values.reduce((sum, currentValue) => sum + currentValue, 0) : null;
      const xValue = datum.timestamp;
      return {
        x: xValue,
        y: sum ?? 0, // y-axis point where tooltip renders
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY:
          sum === null
            ? ['No data']
            : this.dataKeys.map((k) => `${formatNumber([datum[k]])} ${this.label(k)}`),
      };
    });
  }

  get yBounds() {
    const counts: number[] = this.aggregatedData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=4, hardcode 4 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 4 ? 4 : max];
  }

  get xBounds() {
    const domain = this.chartData.map((d) => d.timestamp);
    return new Set(domain);
  }

  // TEMPLATE HELPERS
  barOffset = (bandwidth: number) => (bandwidth - this.barWidth) / 2;

  tooltipX = (original: number, bandwidth: number) => (original + bandwidth / 2).toString();

  tooltipY = (original: number) => (!original ? '0' : `${original}`);

  formatTicksX = (timestamp: string): string => parseAPITimestamp(timestamp, 'M/yy') as string;

  formatTicksY = (num: number): string => numericalAxisLabel(num) || num.toString();
}
