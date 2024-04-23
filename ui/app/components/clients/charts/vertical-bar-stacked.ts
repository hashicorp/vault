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

  get chartData() {
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
    const counts: number[] = this.chartData
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=4, hardcode 4 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 6 ? 6 : max];
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
    return formatNumbers(num) || num.toString();
  };

  get monthlyClientsByType() {
    // each datum needs to be its own object

    const entity_clients = this.args.data.map(({ month, timestamp, entity_clients }) => ({
      month,
      timestamp,
      entity_clients,
      type: 'entity',
    }));
    const non_entity_clients = this.args.data.map(({ month, timestamp, non_entity_clients }) => ({
      month,
      timestamp,
      non_entity_clients,
      type: 'non-entity',
    }));

    const groupedEntity = flatGroup(
      entity_clients,
      (d) => d.month,
      (d) => d.type,
      (d) => d.entity_clients
    );
    const groupedNonEntity = flatGroup(
      non_entity_clients,
      (d) => d.month,
      (d) => d.type,
      (d) => d.non_entity_clients
    );

    return [...groupedEntity, ...groupedNonEntity].map(([month, type, clients]) => ({
      month,
      type,
      clients,
    }));
  }

  get months() {
    return Array.from(new Set(this.monthlyClientsByType.map((d) => d.month)));
  }
}
