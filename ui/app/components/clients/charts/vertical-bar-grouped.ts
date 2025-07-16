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
import ClientsVersionHistoryModel from 'vault/vault/models/clients/version-history';

interface Args {
  legend: Legend[];
  data: MonthlyChartData[];
  upgradeData?: ClientsVersionHistoryModel[];
  chartTitle?: string;
  chartHeight?: number;
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

type ChartDatum = {
  timestamp: string;
  clientType: string;
} & {
  [key in ClientTypes]?: number | undefined;
};

interface UpgradeByMonth {
  [key: string]: ClientsVersionHistoryModel;
}

/**
 * @module VerticalBarGrouped
 * Renders a grouped bar chart of counts for different client types over time. Which client types render
 * is mapped from the "key" values of the @legend arg.
 *
 * @example
 * <Clients::Charts::VerticalBarGrouped
 * @chartTitle="Total monthly usage"
 * @data={{this.flattenedByMonthData}}
 * @legend={{array (hash key="clients" label="Total clients")}}
 * @chartHeight={{250}}
 * />
 */
export default class VerticalBarGrouped extends Component<Args> {
  barWidth = BAR_WIDTH;
  @tracked activeDatum: AggregatedDatum | null = null;

  get chartHeight() {
    return this.args.chartHeight || 190;
  }

  get dataKeys(): ClientTypes[] {
    return this.args.legend.map((l: Legend) => l.key);
  }

  label(legendKey: string) {
    return this.args.legend.find((l: Legend) => l.key === legendKey)?.label;
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
      const maximum = values.length
        ? values.reduce((prev, currentValue) => (prev > currentValue ? prev : currentValue), 0)
        : null;
      const xValue = datum.timestamp;
      const legend = {
        x: xValue,
        y: maximum ?? 0, // y-axis point where tooltip renders
        legendX: parseAPITimestamp(xValue, 'MMMM yyyy') as string,
        legendY:
          maximum === null
            ? ['No data']
            : this.dataKeys.map((k) => `${formatNumber([datum[k]])} ${this.label(k)}`),
        tooltipUpgrade: this.upgradeMessage(datum),
      };
      return legend;
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
    const domain = this.args.data.map((d) => d.timestamp);
    return new Set(domain);
  }

  // UPGRADE STUFF
  get upgradeByMonthYear(): UpgradeByMonth {
    const empty: UpgradeByMonth = {};
    if (!Array.isArray(this.args.upgradeData)) return empty;
    return (
      this.args.upgradeData?.reduce((acc, upgrade) => {
        if (upgrade.timestampInstalled) {
          const key = parseAPITimestamp(upgrade.timestampInstalled, 'M/yy');
          acc[key as string] = upgrade;
        }
        return acc;
      }, empty) || empty
    );
  }

  upgradeMessage(datum: MonthlyChartData) {
    const upgradeInfo = this.upgradeByMonthYear[datum.month as string];
    if (upgradeInfo) {
      const { version, previousVersion } = upgradeInfo;
      return `Vault was upgraded
        ${previousVersion ? 'from ' + previousVersion : ''} to ${version}`;
    }
    return null;
  }

  // TEMPLATE HELPERS
  barOffset = (bandwidth: number, idx = 0) => {
    const withPadding = this.barWidth + 4;
    const moved = (bandwidth - withPadding * this.args.legend.length) / 2;
    return moved + idx * withPadding;
  };

  tooltipX = (original: number, bandwidth: number) => (original + bandwidth / 2).toString();

  tooltipY = (original: number) => (!original ? '0' : `${original}`);

  formatTicksX = (timestamp: string): string => parseAPITimestamp(timestamp, 'M/yy') as string;

  formatTicksY = (num: number): string => numericalAxisLabel(num) || num.toString();
}
