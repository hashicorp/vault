/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { SVG_DIMENSIONS, formatNumbers } from 'vault/utils/chart-helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format, isValid, parse } from 'date-fns';
import { SerializedChartData, UpgradeData } from 'vault/client-counts';
import { debug } from '@ember/debug';

interface Args {
  dataset: SerializedChartData[];
  upgradeData: UpgradeData[];
  xKey?: string;
  yKey?: string;
  chartHeight?: number;
  dateFormat?: string;
}

interface ChartData {
  x: Date;
  y: number | null;
  new: number;
  tooltipUpgrade: string | null;
}

interface UpgradeByMonth {
  [key: string]: UpgradeData;
}

/**
 * @module LineChart
 * LineChart components are used to display time-based data in a line plot with accompanying tooltip
 *
 * @example
 * ```js
 * <LineChart @dataset={{dataset}} @upgradeData={{this.versionHistory}}/>
 * ```
 * @param {string} xKey - string denoting key for x-axis data of dataset. Should reference a date string with format 'M/yy'.
 * @param {string} yKey - string denoting key for y-axis data of dataset. Should reference a number or null.
 * @param {array} upgradeData - array of objects containing version history from the /version-history endpoint
 * @param {number} [chartHeight=190] - height of chart in pixels
 * @param {string} [dateFormat='M/yy'] - format of date string in xKey
 */
export default class LineChart extends Component<Args> {
  // Chart settings
  get yKey() {
    return this.args.yKey || 'clients';
  }
  get xKey() {
    return this.args.xKey || 'month';
  }
  get chartHeight() {
    return this.args.chartHeight || SVG_DIMENSIONS.height;
  }
  get dateFormat() {
    return this.args.dateFormat || 'M/yy';
  }
  // Plot points
  get data(): ChartData[] {
    try {
      return this.args.dataset?.map((datum) => {
        const timestamp = parse(datum[this.xKey] as string, this.dateFormat, new Date()) as Date;
        if (isValid(timestamp) === false)
          throw new Error(
            `Unable to parse value "${datum[this.xKey]}" as date with format ${this.dateFormat}`
          );
        const upgradeMessage = this.getUpgradeMessage(datum);
        return {
          month: datum[this.xKey],
          x: timestamp,
          y: (datum[this.yKey] as number) ?? null,
          new: this.getNewClients(datum),
          tooltipUpgrade: upgradeMessage,
        };
      });
    } catch (e) {
      debug(e as string);
      return [];
    }
  }
  get upgradedMonths() {
    return this.data.filter((datum) => datum.tooltipUpgrade);
  }
  // Domains
  get yDomain() {
    const setMax = Math.max(...this.data.map((datum) => datum.y ?? 0));
    const nearest = setMax > 1000 ? 1000 : setMax > 100 ? 100 : 10;
    // round to nearest 10, 100, or 1000
    return [0, Math.ceil(setMax / nearest) * nearest];
  }
  get timeDomain() {
    // assume data is sorted by time
    const firstTime = this.data[0]?.x;
    const lastTime = this.data[this.data.length - 1]?.x;
    return [firstTime, lastTime];
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

  getUpgradeMessage(datum: SerializedChartData) {
    const upgradeInfo = this.upgradeByMonthYear[datum[this.xKey] as string];
    if (upgradeInfo) {
      const { version, previousVersion } = upgradeInfo;
      return `Vault was upgraded
        ${previousVersion ? 'from ' + previousVersion : ''} to ${version}`;
    }
    return null;
  }
  getNewClients(datum: SerializedChartData) {
    if (!datum?.new_clients) return 0;
    return (datum?.new_clients[this.yKey] as number) || 0;
  }

  hasValue = (count: number | null) => {
    return typeof count === 'number' ? true : false;
  };
  // These functions are used by the tooltip
  formatCount = (count: number) => {
    return formatNumbers([count]);
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
