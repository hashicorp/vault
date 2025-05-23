import Component from '@glimmer/component';
import { HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
import LinealArc from '@lineal-viz/lineal/components/lineal/arc/index.js';
import LinealArcs from '@lineal-viz/lineal/components/lineal/arcs/index.js';
import LinealFluid from '@lineal-viz/lineal/components/lineal/fluid/index.js';
import { concat } from '@ember/helper';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUReportingDonutChart extends Component {
  get data() {
    return (this.args.data || []).map((datum, index) => {
      return {
        ...datum,
        scaleIndex: index + 1
      };
    });
  }
  get total() {
    return this.data.reduce((runningTotal, {
      value
    }) => {
      return runningTotal + value;
    }, 0);
  }
  get a11yLabel() {
    const title = `Total of ${this.total} ${this.args.title}.`;
    const itemsDescription = this.data.map(({
      value,
      label
    }) => {
      return `${value} ${label}`;
    }).join(', ');
    return `${title} Comprised of ${itemsDescription}.`;
  }
  getOffset(width, height) {
    return `translate(${width / 2}, ${height / 2})`;
  }
  getInnerRadius(width, height) {
    const computedRadius = Math.min(width, height) / 2 - 50;
    // Smallest inner radius is 60 to allow for text
    return Math.max(computedRadius, 60);
  }
  getOuterRadius(width, height) {
    // Smallest inner radius is 60 to allow for text
    const computedRadius = Math.min(width, height) / 2;
    return Math.max(computedRadius, 110);
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer ...attributes class=\"ssu-donut-chart__container\" @hasBorder={{true}}>\n      <div class=\"ssu-donut-chart__row\">\n        {{!-- TODO: Figure out glint errors on lineal components --}}\n        {{!-- @glint-expect-error --}}\n        <LinealFluid class=\"ssu-donut-chart__fluid\" as |width height|>\n          <svg width=\"100%\" height=\"100%\" class=\"ssu-donut-chart__chart\" tabindex=\"0\" role=\"img\" aria-label={{this.a11yLabel}}>\n            <g transform={{this.getOffset width height}}>\n              {{!-- @glint-expect-error --}}\n              <LinealArcs @data={{this.data}} {{!-- @glint-expect-error --}} @theta=\"value\" @colorScale=\"nominal\" as |pie|>\n                {{#each pie as |slice|}}\n                  {{!-- @glint-expect-error --}}\n                  <LinealArc data-test-vault-reporting-slice={{slice.data.label}} @startAngle={{slice.startAngle}} @endAngle={{slice.endAngle}} @outerRadius={{this.getOuterRadius width height}} @innerRadius={{this.getInnerRadius width height}} stroke-width=\"2\" class={{slice.cssClass}} />\n                {{/each}}\n              </LinealArcs>\n              <foreignObject transform=\"translate(-60 -30)\" width=\"120\" height=\"120\">\n                <div class=\"ssu-donut-chart__total-summary\">\n                  <HdsTextDisplay @size=\"500\">\n                    {{this.total}}\n                  </HdsTextDisplay>\n                  <HdsTextDisplay @size=\"200\">\n                    {{@title}}\n                  </HdsTextDisplay>\n                </div>\n              </foreignObject>\n            </g>\n          </svg>\n        </LinealFluid>\n        <div class=\"ssu-donut-chart__legend\">\n          {{#each this.data as |datum|}}\n            <HdsTextDisplay data-test-vault-reporting-legend-item={{datum.label}} class=\"ssu-donut-chart__legend-item\n                {{concat \"ssu-donut-chart__legend-item-\" datum.scaleIndex}}\">{{datum.value}} {{datum.label}} </HdsTextDisplay>\n          {{/each}}\n        </div>\n      </div>\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        LinealFluid,
        LinealArcs,
        LinealArc,
        HdsTextDisplay,
        concat
      })
    }), this);
  }
}

export { SSUReportingDonutChart as default };
//# sourceMappingURL=donut-chart.js.map
