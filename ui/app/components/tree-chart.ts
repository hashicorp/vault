/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { TreeChart, TreeChartOptions } from '@carbon/charts';
import '@carbon/charts/styles.css';

// Carbon itself does not define the shape of this data anywhere beyond providing a few dataset examples
interface TreeChartData {
  name: string; // Node name
  value?: number; // Leaf node value
  children?: TreeChartData[];
}

interface Args {
  data: TreeChartData[];
  options: TreeChartOptions;
  title: string;
}

interface CarbonTreeChartSignature {
  Element: HTMLElement;
  Args: Args;
}

// Partial docs are available here
// https://charts.carbondesignsystem.com/api/classes/treechart
export default class CarbonTreeChart extends Component<CarbonTreeChartSignature> {
  @tracked chart: TreeChart | null = null;

  @action
  setupChart(element: HTMLDivElement): void {
    // Create the TreeChart instance
    this.chart = new TreeChart(element, {
      data: this.args.data,
      options: this.args.options,
    });
  }

  @action
  updateChart(): void {
    if (this.chart) {
      // Update the chart with new data
      this.chart.model.setData(this.args.data);
    }
  }

  @action
  destroyChart(): void {
    if (this.chart) {
      this.chart.destroy();
      this.chart = null;
    }
  }
}
