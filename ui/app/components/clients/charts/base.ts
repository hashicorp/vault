/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { BAR_WIDTH, numericalAxisLabel } from 'vault/utils/chart-helpers';

interface Args {
  chartHeight?: number;
}
export default class ChartsBase extends Component<Args> {
  barWidth = BAR_WIDTH;

  get chartHeight() {
    return this.args.chartHeight || 190;
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
    return numericalAxisLabel(num) || num.toString();
  };
}
