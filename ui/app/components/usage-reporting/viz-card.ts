/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { ScaleTypes } from '@carbon/charts';
import type { BarChartOptions } from '@carbon/charts/dist/interfaces';

import type RouterService from '@ember/routing/router-service';
import { CHART_TYPES } from 'vault/modifiers/carbon-chart';
import { toSentenceCase } from 'vault/utils/to-sentence-case';
import type { SimpleDatum } from 'vault/types/usage-reporting';

interface CarbonBarDatum {
  group: string;
  key: string;
  value: number;
}

interface VaultReportingVizCardSignature {
  Args: {
    data: SimpleDatum[];
    title: string;
    description?: string;
    linkText?: string;
    linkIcon?: string;
    linkUrl?: string;
    linkRoute?: string;
    linkTarget?: '_blank' | '_self';
  };
}

const CHART_BAR_WIDTH = 8;

export default class VaultReportingVizCard extends Component<VaultReportingVizCardSignature> {
  @service declare readonly router: RouterService;

  chartType = CHART_TYPES.SIMPLE_BAR;
  private readonly numberFormatter = new Intl.NumberFormat('en-US');
  private readonly singleRowChartHeight = 48;

  get hasData() {
    return this.data.length > 0;
  }

  get data() {
    const values = Array.isArray(this.args.data) ? this.args.data : [];
    return values.filter(({ value }) => value !== 0).sort((a, b) => b.value - a.value);
  }

  get total() {
    return this.data.reduce((runningTotal, { value }) => runningTotal + value, 0);
  }

  get chartData(): CarbonBarDatum[] {
    return this.data.map(({ label, value }) => ({
      group: this.args.title,
      key: toSentenceCase(label),
      value,
    }));
  }

  get tooltipLabel() {
    return this.args.title.slice(0, -1).toLowerCase();
  }

  tooltipCountLabel(value: number) {
    const pluralLabel = this.args.title.toLowerCase();

    if (value === 1) {
      return pluralLabel.endsWith('s') ? pluralLabel.slice(0, -1) : pluralLabel;
    }

    return pluralLabel;
  }

  get maxCount() {
    return Math.max(0, ...this.data.map(({ value }) => value));
  }

  get xAxisIntegerTicks() {
    const maxTicks = 6;

    if (this.maxCount <= 0) {
      return [0];
    }

    const step = Math.max(1, Math.ceil(this.maxCount / (maxTicks - 1)));
    const values: number[] = [];

    for (let i = 0; i <= this.maxCount; i += step) {
      values.push(i);
    }

    if (values[values.length - 1] !== this.maxCount) {
      values.push(this.maxCount);
    }

    return values;
  }

  get categoricalColorScale() {
    return {
      [this.args.title]: 'var(--clients-chart-color-first)',
    };
  }

  get chartOptions(): BarChartOptions {
    const height = this.data.length === 1 ? this.singleRowChartHeight : this.data.length * 20 + 20;

    return {
      height: `${height}px`,
      color: {
        pairing: {
          option: 1,
        },
        scale: this.categoricalColorScale,
      },
      axes: {
        left: {
          mapsTo: 'key',
          scaleType: ScaleTypes.LABELS,
        },
        bottom: {
          mapsTo: 'value',
          scaleType: ScaleTypes.LINEAR,
          domain: [0, this.maxCount],
          ticks: {
            values: this.xAxisIntegerTicks,
            formatter: (tick: number | Date) => {
              if (typeof tick !== 'number') {
                return '';
              }

              return this.numberFormatter.format(tick);
            },
          },
        },
      },
      grid: {
        x: {
          enabled: false,
        },
        y: {
          enabled: false,
        },
      },
      resizable: true,
      legend: {
        enabled: false,
      },
      toolbar: {
        enabled: false,
      },
      bars: {
        width: CHART_BAR_WIDTH,
      },
      tooltip: {
        customHTML: (data: CarbonBarDatum[]) => {
          if (!data?.length) {
            return '';
          }

          const firstPoint = data[0];
          if (!firstPoint) {
            return '';
          }

          const countLabel = this.tooltipCountLabel(firstPoint.value);

          return `
            <div class="usage-reporting-chart-tooltip">
              <p class="usage-reporting-chart-tooltip__label">${firstPoint.key}</p>
              <p class="usage-reporting-chart-tooltip__value">${this.numberFormatter.format(
                firstPoint.value
              )} ${countLabel}</p>
            </div>
          `;
        },
      },
      accessibility: {
        svgAriaLabel: this.args.title,
      },
    };
  }

  get emptyStateTitle() {
    return 'None enabled';
  }

  get emptyStateDescription() {
    return `${this.args.title} in this namespace will appear here.`;
  }

  get emptyStateLinkText() {
    return `Enable ${this.args.title.toLowerCase()}`;
  }

  get description() {
    if (this.hasData) {
      return this.args.description;
    }

    return;
  }

  get linkUrl() {
    if (this.hasData) {
      return this.args.linkUrl;
    }

    return;
  }

  get linkRoute() {
    if (!this.hasData) {
      return;
    }

    return this.args.linkRoute;
  }

  get emptyStateLinkUrl() {
    return this.args.linkUrl;
  }

  get emptyStateLinkRoute() {
    return this.args.linkRoute;
  }
}
