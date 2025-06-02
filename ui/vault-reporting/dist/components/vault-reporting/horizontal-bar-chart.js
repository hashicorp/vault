import Component from '@glimmer/component';
import { HdsApplicationState, HdsTextBody, HdsSeparator, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import LinealFluid from '@lineal-viz/lineal/components/lineal/fluid/index.js';
import LinealHBars from '@lineal-viz/lineal/components/lineal/h-bars/index.js';
import scaleLinear from '@lineal-viz/lineal/helpers/scale-linear.js';
import scaleBand from '@lineal-viz/lineal/helpers/scale-band.js';
import stackH from '@lineal-viz/lineal/helpers/stack-h.js';
import LinealAxis from '@lineal-viz/lineal/components/lineal/axis/index.js';
import axisOffset from '../../modifiers/axis-offset.js';
import { tracked } from '@glimmer/tracking';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUReportingHorizontalBarChart extends Component {
  static {
    g(this.prototype, "xRangeOffsetWidth", [tracked], function () {
      return 0;
    });
  }
  #xRangeOffsetWidth = (i(this, "xRangeOffsetWidth"), void 0);
  get hasData() {
    return this.args.data && Array.isArray(this.args.data) && this.args.data.length > 0;
  }
  get data() {
    if (!this.hasData) {
      return [];
    }
    // Filtering DESC for now per designs, could make configurable if needed
    return this.args.data.filter(({
      value
    }) => value !== 0).sort((a, b) => {
      return b.value - a.value;
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
  get yDomain() {
    return this.data.map(({
      label
    }) => label);
  }
  get xDomain() {
    return [0, Math.max(0, ...this.data.map(({
      value
    }) => value))];
  }
  get rangeHeight() {
    return this.data.length * 26;
  }
  get yRange() {
    return [0, this.rangeHeight];
  }
  get emptyStateTitle() {
    return 'None enabled';
  }
  get emptyStateDescription() {
    const entitiesTitle = this.args.title;
    return `${entitiesTitle} in this namespace will appear here.`;
  }
  get emptyStateLinkText() {
    const entities = this.args.title.toLowerCase();
    return `Enable ${entities}`;
  }
  get description() {
    if (this.hasData) {
      return this.args.description;
    }
  }
  get linkUrl() {
    if (this.hasData) {
      return this.args.linkUrl;
    }
  }
  getXRange = width => {
    return [0, Math.max(0, width - this.xRangeOffsetWidth - 32)];
  };
  handleAxisOffset = offsetWidth => {
    this.xRangeOffsetWidth = offsetWidth;
  };
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer ...attributes class=\"ssu-horizontal-bar-chart__container\" @hasBorder={{true}}>\n      <TitleRow @title={{@title}} @description={{this.description}} @linkUrl={{this.linkUrl}} @linkText={{@linkText}} @linkTarget={{@linkTarget}} @linkIcon={{@linkIcon}} />\n      {{#if this.hasData}}\n        {{!-- TODO: Figure out glint errors on lineal components --}}\n        {{!-- @glint-expect-error --}}\n        <LinealFluid class=\"ssu-horizontal-bar-chart__chart\" as |width|>\n          <svg height={{this.rangeHeight}} width=\"100%\" {{axisOffset this.handleAxisOffset 8}} data-test-vault-reporting-horizontal-bar-chart-svg>\n            {{!-- We are using the stacked version of the HBars as there seems to be an issue in the non-stacked version for how the x position is calculated.  --}}\n            {{#let (scaleLinear range=(this.getXRange width) domain=this.xDomain) (scaleBand range=this.yRange domain=this.yDomain) (stackH data=this.data x=\"value\" y=\"label\" z=\"\") as |xScale yScale stacked|}}\n              {{#if xScale.isValid}}\n                <LinealAxis @scale={{yScale}} {{!-- @glint-expect-error --}} @orientation=\"left\" @includeDomain={{false}} />\n                {{!-- TODO: Extra wrapper exists only for test attribute, figure out a better way --}}\n                <g data-test-vault-reporting-horizontal-bar-chart-bars>\n                  <LinealHBars @data={{stacked.data}} {{!-- @glint-expect-error --}} @x=\"x\" {{!-- @glint-expect-error --}} @y=\"y\" {{!-- @glint-expect-error --}} @height={{6}} @xScale={{xScale}} @yScale={{yScale}} />\n                </g>\n                <g>\n                  {{!-- @glint-expect-error --}}\n                  {{#each stacked.data as |dataset|}}\n                    {{#each dataset as |datum|}}\n                      <text class=\"ssu-horizontal-bar-chart__label\" {{!-- @glint-expect-error --}} y={{yScale.compute datum.y}} x={{xScale.compute datum.x}} dy=\"17.5px\" dx=\"8px\" data-test-vault-reporting-horizontal-bar-chart-inline-count aria-label=\"{{datum.y}} {{datum.x}}\">\n                        {{datum.x}}\n                      </text>\n                    {{/each}}\n                  {{/each}}\n                </g>\n              {{/if}}\n            {{/let}}\n          </svg>\n        </LinealFluid>\n        <HdsSeparator class=\"ssu-horizontal-bar-chart__separator\" @spacing=\"0\" />\n        <HdsTextBody class=\"ssu-horizontal-bar-chart__total\" @size=\"200\" @tag=\"p\" data-test-vault-reporting-horizontal-bar-chart-total>\n          Total:\n          {{this.total}}\n        </HdsTextBody>\n      {{else}}\n\n        <HdsApplicationState data-test-vault-reporting-horizontal-bar-chart-empty-state class=\"ssu-horizontal-bar-chart__empty-state\" as |A|>\n          {{#if (has-block \"empty\")}}\n            {{yield A to=\"empty\"}}\n          {{else}}\n            <A.Header data-test-vault-reporting-horizontal-bar-chart-empty-state-title @title={{this.emptyStateTitle}} />\n            <A.Body data-test-vault-reporting-horizontal-bar-chart-empty-state-description @text={{this.emptyStateDescription}} />\n            {{#if @linkUrl}}\n              <A.Footer as |F|>\n                <F.LinkStandalone data-test-vault-reporting-horizontal-bar-chart-empty-state-link @icon=\"plus\" @text={{this.emptyStateLinkText}} @href={{@linkUrl}} />\n              </A.Footer>\n            {{/if}}\n          {{/if}}\n        </HdsApplicationState>\n      {{/if}}\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        TitleRow,
        LinealFluid,
        axisOffset,
        scaleLinear,
        scaleBand,
        stackH,
        LinealAxis,
        LinealHBars,
        HdsSeparator,
        HdsTextBody,
        HdsApplicationState
      })
    }), this);
  }
}

export { SSUReportingHorizontalBarChart as default };
//# sourceMappingURL=horizontal-bar-chart.js.map
