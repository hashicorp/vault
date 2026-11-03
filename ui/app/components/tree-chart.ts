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
  private toolbarObserver: MutationObserver | null = null;

  @action
  setupChart(element: HTMLDivElement): void {
    // Create the TreeChart instance
    this.chart = new TreeChart(element, {
      data: this.args.data,
      options: this.args.options,
    });
    this.toolbarObserver = new MutationObserver(() => this.updateToolbarAccessibility(element));
    this.toolbarObserver.observe(element, { childList: true, subtree: true });
    this.updateToolbarAccessibility(element);
  }

  // Carbon 1.27 wraps native buttons in interactive roles, including in its dynamically rendered overflow menu.
  private updateToolbarAccessibility(element: HTMLElement): void {
    element.querySelectorAll('.toolbar-control[role="button"]').forEach((control) => {
      const button = control.querySelector('button');
      if (button) {
        const disabled = control.getAttribute('aria-disabled');
        if (disabled !== null) {
          button.setAttribute('aria-disabled', disabled);
        }
        control.removeAttribute('role');
        control.removeAttribute('aria-disabled');
        control.removeAttribute('aria-label');
      }
    });
    element.querySelectorAll('.cds--overflow-menu-options [role="menu"] > li').forEach((item) => {
      const button = item.querySelector('button');
      if (button) {
        button.setAttribute('role', 'menuitem');
        const disabled = item.getAttribute('aria-disabled');
        if (disabled !== null) {
          button.setAttribute('aria-disabled', disabled);
        }
        item.setAttribute('role', 'none');
        item.removeAttribute('tabindex');
        item.removeAttribute('aria-disabled');
      }
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
    this.toolbarObserver?.disconnect();
    this.toolbarObserver = null;
    if (this.chart) {
      this.chart.destroy();
      this.chart = null;
    }
  }
}
