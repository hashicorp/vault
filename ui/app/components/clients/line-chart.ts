/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { SVG_DIMENSIONS, formatNumbers } from 'vault/utils/chart-helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { format } from 'date-fns';
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
 * @param {string} xKey - string denoting key for x-axis data of dataset. Should point to timestamp string.
 * @param {string} yKey - string denoting key for y-axis data of dataset. Should point to number.
 * @param {array} upgradeData - array of objects containing version history from the /version-history endpoint
 * @param {number} [chartHeight=190] - height of chart in pixels
 */
export default class LineChart extends Component<Args> {
  get yKey() {
    return this.args.yKey || 'clients';
  }
  get xKey() {
    return this.args.xKey || 'timestamp';
  }
  get chartHeight() {
    return this.args.chartHeight || SVG_DIMENSIONS.height;
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

  tooltipData(datum: SerializedChartData) {
    const upgradeKey = parseAPITimestamp(datum[this.xKey], 'M/yy') as string;
    const upgradeInfo = this.upgradeByMonthYear[upgradeKey];
    if (upgradeInfo) {
      const { version, previousVersion } = upgradeInfo;
      return `Vault was upgraded
        ${previousVersion ? 'from ' + previousVersion : ''} to ${version}`;
    }
    return null;
  }

  get upgradedMonths() {
    return this.data.filter((datum) => datum.tooltipUpgrade);
  }
  newClients(datum: SerializedChartData) {
    if (!datum?.new_clients) return 0;
    return (datum?.new_clients[this.yKey] as number) || 0;
  }
  get data(): ChartData[] {
    return this.args.dataset?.map((datum) => {
      const date = parseAPITimestamp(datum[this.xKey]) as Date;
      const upgradeMessage = this.tooltipData(datum);
      return {
        x: date,
        y: (datum[this.yKey] as number) ?? null,
        new: this.newClients(datum),
        tooltipUpgrade: upgradeMessage,
      };
    });
  }

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
    const offset = `${this.chartHeight - original + 20}`;
    return offset;
  };
}
