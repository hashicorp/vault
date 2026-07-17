/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { toLabel } from 'core/helpers/to-label';
import { ScaleTypes, TruncationTypes } from '@carbon/charts';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import type { BarChartOptions, DonutChartOptions } from '@carbon/charts/dist/interfaces';
import { CHART_TYPES } from 'vault/modifiers/carbon-chart';

import type { ByMonthClients, ByMonthNewClients, TotalClients } from 'vault/vault/client-counts/activity-api';
import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';

interface ChartDataPoint {
  group: string;
  key: string;
  value: number | null;
  legendX?: string;
}

interface Args {
  byMonthClients: ByMonthClients[] | ByMonthNewClients[];
  runningTotals: TotalClients;
}

interface DonutTooltipDatum {
  label?: string;
  group?: string;
  value: number | null;
}

interface DonutTooltipItem {
  label?: string;
  group?: string;
  value?: number | string | null;
}

const CHART_HEIGHT = '300px';
const MIN_STACKED_Y_AXIS_MAX = 4;

export default class RunningTotal extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  @tracked showStacked = false;

  // Export chart types for use in template
  chartTypes = CHART_TYPES;

  get chartContainerText() {
    const range = this.version.isEnterprise ? 'billing period' : 'date range';
    return this.flags.isHvdManaged
      ? 'Number of total unique clients in the data period by client type, and total number of unique clients per month. The monthly total is the relevant billing metric.'
      : `Number of clients in the ${range} by client type, and a breakdown of new clients per month during the ${range}. `;
  }

  get dataKey() {
    return this.flags.isHvdManaged ? 'clients' : 'new_clients';
  }

  get runningTotalData() {
    // The parent component determines whether `monthly.clients` in @byMonthClients represents "new" or "total" clients per month.
    // (We render "new" for self-managed clusters and "total" for HVD-managed.)
    // As a result, we do not use `this.dataKey` to select a property from `monthly` but to add a superficial key
    // to the data that ensures the chart tooltip and legend text render appropriately.
    return this.args.byMonthClients.map((monthly) => ({
      ...monthly,
      [this.dataKey]: monthly.clients,
    }));
  }

  get donutChartData() {
    return [
      { group: 'Entity clients', value: this.args.runningTotals.entity_clients },
      { group: 'Non-entity clients', value: this.args.runningTotals.non_entity_clients },
      { group: 'ACME clients', value: this.args.runningTotals.acme_clients },
      ...(this.flags.secretsSyncIsActivated
        ? [{ group: 'Secret sync clients', value: this.args.runningTotals.secret_syncs }]
        : []),
    ];
  }

  formatTooltipLabel(label: string) {
    if (!label) {
      return label;
    }

    if (label.startsWith('ACME')) {
      return label;
    }

    return `${label.charAt(0).toLowerCase()}${label.slice(1)}`;
  }

  formatTooltipValue(value: number, label: string) {
    return `${value.toLocaleString()} ${this.formatTooltipLabel(label)}`;
  }

  getDonutTooltipLabel(data: unknown, datum: DonutTooltipDatum) {
    const tooltipItem = Array.isArray(data) ? (data[0] as DonutTooltipItem | undefined) : undefined;

    return tooltipItem?.label || tooltipItem?.group || datum?.group || datum?.label || '';
  }

  get donutChartOptions(): DonutChartOptions {
    const title = 'Client count and type distribution';
    const donutLabel = this.flags.isHvdManaged ? 'Total unique clients' : 'Total clients';
    const total = this.args.runningTotals.clients;
    const legendOrder = this.donutChartData.map((d) => d.group);
    return {
      title,
      color: {
        scale: this.categoricalColorScale,
      },
      pie: {
        sortFunction: () => 0,
      },
      donut: {
        center: {
          label: donutLabel,
          number: total,
          numberFormatter: (value: number) => value.toLocaleString(),
        },
      },
      legend: {
        enabled: true,
        alignment: 'center',
        order: legendOrder,
        truncation: {
          type: TruncationTypes.NONE,
        },
      },
      toolbar: {
        enabled: false,
      },
      tooltip: {
        customHTML: (data: unknown, _defaultHTML: string, datum: DonutTooltipDatum) => {
          if (!datum || datum.value === null) return '';

          const label = this.getDonutTooltipLabel(data, datum);
          if (!label) return '';

          return this.formatTooltipValue(datum.value, label);
        },
      },
      accessibility: {
        svgAriaLabel: title,
      },
      height: CHART_HEIGHT,
    };
  }

  get chartLegend() {
    if (this.showStacked) {
      return this.stackedLegend;
    }
    return [{ key: this.dataKey, label: toLabel([this.dataKey]) }];
  }

  get stackedLegend() {
    return [
      { key: 'entity_clients', label: 'Entity clients' },
      { key: 'non_entity_clients', label: 'Non-entity clients' },
      { key: 'acme_clients', label: 'ACME clients' },
      // MUST BE LAST because conditionally renders and legend color mapping for stacked bars will be off otherwise
      ...(this.flags.secretsSyncIsActivated ? [{ key: 'secret_syncs', label: 'Secret sync clients' }] : []),
    ];
  }

  get categoricalColorScale() {
    return {
      'Entity clients': 'var(--clients-chart-color-first)',
      'Non-entity clients': 'var(--clients-chart-color-second)',
      'ACME clients': 'var(--clients-chart-color-third)',
      'Secret sync clients': 'var(--clients-chart-color-fourth)',
    };
  }

  get simpleColorScale() {
    return {
      'New clients': 'var(--clients-chart-color-single)',
      Clients: 'var(--clients-chart-color-single)',
    };
  }

  /**
   * Transforms monthly data into simple bar chart format
   */
  get simpleChartData(): ChartDataPoint[] {
    if (!this.runningTotalData || this.runningTotalData.length === 0) {
      return [];
    }

    return this.runningTotalData.map((monthData) => {
      const value = (monthData as unknown as Record<string, number | string>)[this.dataKey];
      const month = parseAPITimestamp(monthData.timestamp, 'M/yy');
      const label = toLabel([this.dataKey]);

      return {
        group: label,
        key: month,
        value: typeof value === 'number' ? value : null,
        legendX: parseAPITimestamp(monthData.timestamp, 'MMMM yyyy'),
      };
    });
  }

  /**
   * Transforms monthly data into stacked bar chart format
   */
  get stackedChartData(): ChartDataPoint[] {
    if (!this.runningTotalData || this.runningTotalData.length === 0) {
      return [];
    }

    const result: ChartDataPoint[] = [];

    this.runningTotalData.forEach((monthData) => {
      const month = parseAPITimestamp(monthData.timestamp, 'M/yy');
      const formattedMonth = parseAPITimestamp(monthData.timestamp, 'MMMM yyyy');

      this.stackedLegend.forEach((legend) => {
        const rawValue = (monthData as unknown as Record<string, number | null>)[legend.key];
        result.push({
          group: legend.label,
          key: month,
          value: typeof rawValue === 'number' ? rawValue : null,
          legendX: formattedMonth,
        });
      });
    });

    return result;
  }

  /**
   * Calculates the y-axis domain for simple chart
   */
  get simpleYDomain(): [number, number] {
    const values = this.simpleChartData.map((d) => d.value).filter((v): v is number => v !== null);
    const max = Math.max(...values, 0);
    return [0, Math.ceil(max * 1.1)];
  }

  /**
   * Calculates the y-axis domain for stacked chart
   * Calculates the maximum stacked total for each month
   */
  get stackedYDomain(): [number, number] {
    // Calculate the sum of all client types for each month
    const stackedTotals = new Map<string, number>();

    this.runningTotalData.forEach((monthData) => {
      const timestamp = monthData.timestamp;
      const total = this.stackedLegend.reduce((sum, legend) => {
        const value = (monthData as unknown as Record<string, number | string>)[legend.key] || 0;
        return sum + (typeof value === 'number' ? value : 0);
      }, 0);
      stackedTotals.set(timestamp, total);
    });

    const max = Math.max(...Array.from(stackedTotals.values()), 0);
    return [0, Math.max(max, MIN_STACKED_Y_AXIS_MAX)];
  }

  /**
   * Generates Carbon Charts configuration options for simple bar chart
   */
  get simpleChartOptions(): BarChartOptions {
    return {
      title: 'Client usage by month',
      color: {
        pairing: {
          option: 3,
        },
        scale: this.simpleColorScale,
      },
      height: CHART_HEIGHT,
      axes: {
        left: {
          mapsTo: 'value',
          scaleType: ScaleTypes.LINEAR,
          domain: this.simpleYDomain,
        },
        bottom: {
          mapsTo: 'key',
          scaleType: ScaleTypes.LABELS,
        },
      },
      legend: {
        alignment: 'center',
        enabled: true,
        truncation: {
          type: TruncationTypes.NONE,
        },
      },
      toolbar: {
        enabled: false,
      },
      bars: {
        maxWidth: 20,
      },
      tooltip: {
        customHTML: (data: ChartDataPoint[]) => {
          if (!data || data.length === 0) return '';

          const firstPoint = data[0];
          if (!firstPoint) return '';

          const month = firstPoint.legendX || firstPoint.key;
          const value = firstPoint.value;
          const label = toLabel([this.dataKey]);

          if (value === null) {
            return `
              <div class="cds--tooltip cds--tooltip--shown carbon-chart-tooltip">
                <p class="tooltip-month">${month}</p>
                <p class="tooltip-value">No data</p>
              </div>
            `;
          }

          return `
            <div class="cds--tooltip cds--tooltip--shown carbon-chart-tooltip">
              <p class="tooltip-month">${month}</p>
              <p class="tooltip-value">${this.formatTooltipValue(value, label)}</p>
            </div>
          `;
        },
      },
    };
  }

  /**
   * Generates Carbon Charts configuration options for stacked bar chart
   */
  get stackedChartOptions(): BarChartOptions {
    return {
      title: 'Client usage by month',
      height: CHART_HEIGHT,
      color: {
        pairing: {
          option: 1,
        },
        scale: this.categoricalColorScale,
      },
      axes: {
        bottom: {
          mapsTo: 'key',
          scaleType: ScaleTypes.LABELS,
        },
        left: {
          mapsTo: 'value',
          scaleType: ScaleTypes.LINEAR,
          domain: this.stackedYDomain,
        },
      },
      bars: {
        maxWidth: 20,
      },
      legend: {
        enabled: true,
        alignment: 'center',
        truncation: {
          type: TruncationTypes.NONE,
        },
      },
      toolbar: {
        enabled: false,
      },
      tooltip: {
        customHTML: (data: ChartDataPoint[]) => {
          if (!data || data.length === 0) return '';

          const firstPoint = data[0];
          if (!firstPoint) return '';

          const month = firstPoint.legendX || firstPoint.key;
          const hasData = data.some((d) => d.value !== null);

          if (!hasData) {
            return `
              <div class="cds--tooltip cds--tooltip--shown carbon-chart-tooltip">
                <p class="tooltip-month">${month}</p>
                <p class="tooltip-value">No data</p>
              </div>
            `;
          }

          const rows = data
            .map((d) => {
              return `<p class="tooltip-value">${this.formatTooltipValue(d.value ?? 0, d.group)}</p>`;
            })
            .join('');

          return `
            <div class="cds--tooltip cds--tooltip--shown carbon-chart-tooltip">
              <p class="tooltip-month">${month}</p>
              ${rows}
            </div>
          `;
        },
      },
    };
  }
}
