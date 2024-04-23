/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// @ts-nocheck
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { BAR_WIDTH, formatNumbers } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { flatGroup } from 'd3-array';
import type { MonthlyChartData } from 'vault/vault/charts/client-counts';
import type { TotalClients } from 'core/utils/client-count-utils';

interface Args {
  data: MonthlyChartData[];
  dataKeys: string[];
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
 * @module VerticalBarStacked
 * Renders a vertical bar chart of counts fora single data point (@dataKey) over time.
 *
 * @example
 <Clients::Charts::VerticalBarStacked
    @chartTitle="Secret Sync client counts"
    @data={{this.model}}
    @dataKey="secret_syncs"
    @showTable={{true}}
    @chartHeight={{200}}
  />
 */
export default class VerticalBarStacked extends Component<Args> {
  barWidth = BAR_WIDTH;

  @tracked activeDatum: ChartData | null = null;

  get chartHeight() {
    return this.args.chartHeight || 190;
  }

  get monthlyClientsByType() {
    // each datum needs to be its own object
    let dataset = [];
    for (const key of this.args.dataKeys) {
      const keyData = this.args.data.map((d) => ({
        month: parseAPITimestamp(d.timestamp, 'M/yy'),
        clientType: key,
        [key]: d[key],
      }));
      dataset = [
        ...dataset,
        ...flatGroup(
          keyData,
          // order here must match destructure order in return below
          (d) => d.month,
          (d) => d.clientType,
          (d) => d[key] // client counts for the datakey
        ),
      ];
    }

    return dataset.map(([month, clientType, counts]) => ({
      month,
      clientType,
      counts,
    }));
  }

  getSum = (array, key) => {
    return array.reduce((sum, { counts }) => sum + counts[key], 0);
  };
  // for yDomain scale, tooltip target area and tooltip text data
  get aggregatedData() {
    return this.args.data.map((d): ChartData => {
      const xValue = d.timestamp as string;
      const yValue =
        d?.entity_clients && d?.non_entity_clients ? d.entity_clients + d.non_entity_clients : null;
      const entity = (d['entity_clients'] as number) ?? null;
      const nonEntity = (d['non_entity_clients'] as number) ?? null;

      return {
        x: parseAPITimestamp(xValue, 'M/yy') as string,
        y: yValue ?? 0,
        tooltip:
          yValue === null
            ? 'No data'
            : `${formatNumber([entity])} entity, ${formatNumber([nonEntity])} non-entity `,
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY: (yValue ?? 'No data').toString(),
      };
    });
  }

  get yDomain() {
    const counts: number[] = this.aggregatedData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=6, hardcode 6 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 6 ? 6 : max];
  }

  get xDomain() {
    const months = this.monthlyClientsByType.map((d) => d.month);
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
    return formatNumbers(num) || num.toString();
  };
}
