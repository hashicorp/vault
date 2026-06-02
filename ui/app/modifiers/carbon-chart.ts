/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { modifier } from 'ember-modifier';
import { SimpleBarChart, StackedBarChart, DonutChart } from '@carbon/charts';
import type { BarChartOptions, DonutChartOptions } from '@carbon/charts/dist/interfaces';

/**
 * Chart type constants for Carbon Charts
 */
export const CHART_TYPES = {
  SIMPLE_BAR: 'simple',
  STACKED_BAR: 'stacked',
  DONUT: 'donut',
} as const;

export type ChartType = (typeof CHART_TYPES)[keyof typeof CHART_TYPES];

interface ChartDataPoint {
  group: string;
  value: number | null;
  [key: string]: string | number | null;
}

interface CarbonChartModifierSignature {
  Element: HTMLDivElement;
  Args: {
    Positional: [ChartDataPoint[], BarChartOptions | DonutChartOptions, ChartType];
  };
}

/**
 * Custom modifier for managing Carbon Chart lifecycle.
 * Replaces the need for did-insert, did-update, and will-destroy render modifiers.
 *
 * @example
 * ```hbs
 * <div {{carbon-chart @chartData @chartOptions @chartType}}></div>
 * ```
 */
const CHART_CLASS_MAP = {
  [CHART_TYPES.SIMPLE_BAR]: SimpleBarChart,
  [CHART_TYPES.STACKED_BAR]: StackedBarChart,
  [CHART_TYPES.DONUT]: DonutChart,
} as const;

export default modifier<CarbonChartModifierSignature>((element, [chartData, chartOptions, chartType]) => {
  let chart: SimpleBarChart | StackedBarChart | DonutChart | null = null;

  if (chartData && Array.isArray(chartData) && chartData.length > 0) {
    const ChartClass = CHART_CLASS_MAP[chartType];
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    chart = new ChartClass(element as HTMLDivElement, { data: chartData, options: chartOptions as any });
  }

  // Return cleanup function
  return () => {
    if (chart) {
      chart.destroy();
      chart = null;
    }
  };
});
