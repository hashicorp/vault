/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { SVG_DIMENSIONS, numericalAxisLabel } from 'vault/utils/chart-helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format, isValid } from 'date-fns';
import { debug } from '@ember/debug';

import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type { MonthlyChartData, Timestamp } from 'vault/vault/charts/client-counts';
import type { TotalClients } from 'core/utils/client-count-utils';

interface Args {
  dataset: MonthlyChartData[];
  upgradeData: ClientsVersionHistoryModel[];
  xKey?: string;
  yKey?: string;
  chartHeight?: number;
}

interface ChartData {
  x: Date;
  y: number | null;
  new: number;
  tooltipUpgrade: string | null;
  month: string; // used for test selectors and to match key on upgradeData
}

interface UpgradeByMonth {
  [key: string]: ClientsVersionHistoryModel;
}

/**
 * @module LineChart
 * LineChart components are used to display time-based data in a line plot with accompanying tooltip
 *
 * @example
 * ```js
 * <LineChart @dataset={{dataset}} @upgradeData={{this.versionHistory}}/>
 * ```
 * @param {array} dataset - array of objects containing data to be plotted
 * @param {string} [xKey=clients] - string denoting key for x-axis data of dataset. Should reference a timestamp string.
 * @param {string} [yKey=timestamp] - string denoting key for y-axis data of dataset. Should reference a number or null.
 * @param {array} [upgradeData=null] - array of objects containing version history from the /version-history endpoint
 * @param {number} [chartHeight=190] - height of chart in pixels
 */
export default class LineChart extends Component<Args> {
  // Chart settings
  get yKey() {
    return this.args.yKey || 'clients';
  }
  get xKey() {
    return this.args.xKey || 'timestamp';
  }
  get chartHeight() {
    return this.args.chartHeight || SVG_DIMENSIONS.height;
  }
  // Plot points
  get data(): ChartData[] {
    try {
      return this.args.dataset?.map((datum) => {
        const timestamp = parseAPITimestamp(datum[this.xKey as keyof Timestamp]) as Date;
        if (isValid(timestamp) === false)
          throw new Error(`Unable to parse value "${datum[this.xKey as keyof Timestamp]}" as date`);
        const upgradeMessage = this.getUpgradeMessage(datum);
        return {
          x: timestamp,
          y: (datum[this.yKey as keyof TotalClients] as number) ?? null,
          new: this.getNewClients(datum),
          tooltipUpgrade: upgradeMessage,
          month: datum.month,
        };
      });
    } catch (e) {
      debug(e as string);
      return [];
    }
  }
  get upgradedMonths() {
    // only render upgrade month circle if datum has client count data (the y value)
    return this.data.filter((datum) => datum.tooltipUpgrade && datum.y);
  }
  // Domains
  get yDomain() {
    const counts: number[] = this.data
      .map((d) => d.y)
      .flatMap((num) => (typeof num === 'number' ? [num] : []));
    const max = Math.max(...counts);
    // if max is <=4, hardcode 4 which is the y-axis tickCount so y-axes are not decimals
    return [0, max <= 4 ? 4 : max];
  }

  get xDomain() {
    // these values are date objects but are already in chronological order so we use scale-point (instead of scale-time)
    // which calculates the x-scale based on the number of data points
    return this.data.map((d) => d.x);
  }

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

  getUpgradeMessage(datum: MonthlyChartData) {
    const upgradeInfo = this.upgradeByMonthYear[datum.month as string];
    if (upgradeInfo) {
      const { version, previousVersion } = upgradeInfo;
      return `Vault was upgraded
        ${previousVersion ? 'from ' + previousVersion : ''} to ${version}`;
    }
    return null;
  }
  getNewClients(datum: MonthlyChartData) {
    if (!datum?.new_clients) return 0;
    return (datum?.new_clients[this.yKey as keyof TotalClients] as number) || 0;
  }

  // TEMPLATE HELPERS
  hasValue = (count: number | null) => {
    return typeof count === 'number' ? true : false;
  };
  formatCount = (num: number): string => {
    return numericalAxisLabel(num) || num.toString();
  };
  formatMonth = (date: Date) => {
    return format(date, 'M/yy');
  };
  tooltipX = (original: number) => {
    return original.toString();
  };
  tooltipY = (original: number) => {
    return `${this.chartHeight - original + 15}`;
  };
}
