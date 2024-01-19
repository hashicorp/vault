/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { SVG_DIMENSIONS, formatNumbers } from 'vault/utils/chart-helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format, parse } from 'date-fns';
import { SerializedChartData, UpgradeData } from 'vault/client-counts';

interface Args {
  dataset: SerializedChartData[];
  upgradeData: UpgradeData[];
  xKey?: string;
  yKey?: string;
  chartHeight?: number;
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
 * LineChart components are used to display data in a line plot with accompanying tooltip
 *
 * @example
 * ```js
 * <LineChart @dataset={{dataset}} @upgradeData={{this.versionHistory}}/>
 * ```
 * @param {string} xKey - string denoting key for x-axis data of dataset. Should reference a date string with format 'M/yy'.
 * @param {string} yKey - string denoting key for y-axis data of dataset. Should reference a number or null.
 * @param {array} upgradeData - array of objects containing version history from the /version-history endpoint
 * @param {number} [chartHeight=190] - height of chart in pixels
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
  // Plot points
  get data(): ChartData[] {
    return this.args.dataset?.map((datum) => {
      // We expect the xKey to be formatted like 'M/yy'
      const timestamp = parse(datum[this.xKey] as string, 'M/yy', new Date()) as Date;
      const upgradeMessage = this.getUpgradeMessage(datum);
      return {
        month: datum[this.xKey],
        x: timestamp,
        y: (datum[this.yKey] as number) ?? null,
        new: this.getNewClients(datum),
        tooltipUpgrade: upgradeMessage,
      };
    });
  }
  get upgradedMonths() {
    return this.data.filter((datum) => datum.tooltipUpgrade);
  }
  get yDomain() {
    const setMax = Math.max(...this.data.map((datum) => datum.y ?? 0));
    const nearest = setMax < 1500 ? 200 : 2000;
    // Round to upper 200 or 2000
    return [0, Math.ceil(setMax / nearest) * nearest];
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
