/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { stack } from 'd3-shape';
// eslint-disable-next-line no-unused-vars
import { select, event, selectAll } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { axisLeft } from 'd3-axis';
import { max, maxIndex } from 'd3-array';
import { BAR_COLOR_HOVER, GREY, LIGHT_AND_DARK_BLUE, formatTooltipNumber } from 'vault/utils/chart-helpers';
import { tracked } from '@glimmer/tracking';
import { formatNumber } from 'core/helpers/format-number';

/**
 * @module HorizontalBarChart
 * HorizontalBarChart components are used to display data in the form of a horizontal, stacked bar chart with accompanying tooltip.
 *
 * @example
 * ```js
 * <HorizontalBarChart @dataset={{@dataset}} @chartLegend={{@chartLegend}}/>
 * ```
 * @param {array} dataset - dataset for the chart, must be an array of flattened objects
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {string} labelKey - string of key name for label value in chart data
 * @param {string} xKey - string of key name for x value in chart data
 * @param {object} totalCounts - object to calculate percentage for tooltip
 * @param {string} [noDataMessage] - custom empty state message that displays when no dataset is passed to the chart
 */

// SIZING CONSTANTS
const CHART_MARGIN = { top: 10, left: 95 }; // makes space for y-axis legend
const TRANSLATE = { down: 14, left: 99 };
const CHAR_LIMIT = 15; // character count limit for y-axis labels to trigger truncating
const LINE_HEIGHT = 24; // each bar w/ padding is 24 pixels thick

export default class HorizontalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipText = '';
  @tracked isLabel = null;

  get labelKey() {
    return this.args.labelKey || 'label';
  }

  get xKey() {
    return this.args.xKey || 'clients';
  }

  get chartLegend() {
    return this.args.chartLegend;
  }

  get topNamespace() {
    return this.args.dataset[maxIndex(this.args.dataset, (d) => d[this.xKey])];
  }

  get total() {
    return this.args.totalCounts[this.xKey] || null;
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, [chartData]) {
    // chart legend tells stackFunction how to stack/organize data
    // creates an array of data for each key name
    // each array contains coordinates for each data bar
    const stackFunction = stack().keys(this.chartLegend.map((l) => l.key));
    const dataset = chartData;
    const stackedData = stackFunction(dataset);
    const labelKey = this.labelKey;
    const xKey = this.xKey;
    const xScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d[xKey]))])
      .range([0, 75]); // 25% reserved for margins

    const yScale = scaleBand()
      .domain(dataset.map((d) => d[labelKey]))
      .range([0, dataset.length * LINE_HEIGHT])
      .paddingInner(0.765); // percent of the total width to reserve for padding between bars

    const chartSvg = select(element);
    chartSvg.attr('width', '100%').attr('viewBox', `0 0 564 ${(dataset.length + 1) * LINE_HEIGHT}`);

    const dataBarGroup = chartSvg
      .selectAll('g')
      .remove()
      .exit()
      .data(stackedData)
      .enter()
      .append('g')
      .attr('data-test-group', (d) => `${d.key}`)
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => LIGHT_AND_DARK_BLUE[i]);

    const yAxis = axisLeft(yScale).tickSize(0);

    const yLabelsGroup = chartSvg
      .append('g')
      .attr('data-test-group', 'y-labels')
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`);
    yAxis(yLabelsGroup);

    chartSvg.select('.domain').remove();

    const truncate = (selection) =>
      selection.text((string) =>
        string.length < CHAR_LIMIT ? string : string.slice(0, CHAR_LIMIT - 3) + '...'
      );

    chartSvg.selectAll('.tick text').call(truncate);

    dataBarGroup
      .selectAll('rect')
      .remove()
      .exit()
      // iterate through the stacked data and chart respectively
      .data((stackedData) => stackedData)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', (chartData) => `${xScale(Math.abs(chartData[1] - chartData[0]))}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', (chartData) => `${xScale(chartData[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3);

    const actionBarGroup = chartSvg.append('g').attr('data-test-group', 'action-bars');

    const actionBars = actionBarGroup
      .selectAll('.action-bar')
      .remove()
      .exit()
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'action-bar')
      .attr('width', '100%')
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', (chartData) => yScale(chartData[labelKey]))
      .style('fill', `${GREY}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    const labelActionBarGroup = chartSvg.append('g').attr('data-test-group', 'label-action-bars');

    const labelActionBar = labelActionBarGroup
      .selectAll('.label-action-bar')
      .remove()
      .exit()
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'label-action-bar')
      .attr('width', CHART_MARGIN.left)
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', (chartData) => yScale(chartData[labelKey]))
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    const dataBars = chartSvg.selectAll('rect.data-bar');
    const actionBarSelection = chartSvg.selectAll('rect.action-bar');

    const compareAttributes = (elementA, elementB, attr) =>
      select(elementA).attr(`${attr}`) === select(elementB).attr(`${attr}`);

    // MOUSE EVENTS FOR DATA BARS
    actionBars
      .on('mouseover', (data) => {
        const hoveredElement = actionBars.filter((bar) => bar[labelKey] === data[labelKey]).node();
        this.tooltipTarget = hoveredElement;
        this.isLabel = false;
        this.tooltipText = this.total
          ? `${Math.round((data[xKey] * 100) / this.total)}% 
        of total client counts:
        ${formatTooltipNumber(data.entity_clients)} entity clients, 
        ${formatTooltipNumber(data.non_entity_clients)} non-entity clients.`
          : '';

        select(hoveredElement).style('opacity', 1);

        dataBars
          .filter(function () {
            return compareAttributes(this, hoveredElement, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
      })
      .on('mouseout', function () {
        select(this).style('opacity', 0);
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${LIGHT_AND_DARK_BLUE[i]}`);
      });

    // MOUSE EVENTS FOR Y-AXIS LABELS
    labelActionBar
      .on('mouseover', (data) => {
        if (data[labelKey].length >= CHAR_LIMIT) {
          const hoveredElement = labelActionBar.filter((bar) => bar[labelKey] === data[labelKey]).node();
          this.tooltipTarget = hoveredElement;
          this.isLabel = true;
          this.tooltipText = data[labelKey];
        } else {
          this.tooltipTarget = null;
        }
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        actionBarSelection
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '1');
      })
      .on('mouseout', function () {
        this.tooltipTarget = null;
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${LIGHT_AND_DARK_BLUE[i]}`);
        actionBarSelection
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '0');
      });

    // client count total values to the right
    const totalValueGroup = chartSvg
      .append('g')
      .attr('data-test-group', 'total-values')
      .attr('transform', `translate(${TRANSLATE.left}, ${TRANSLATE.down})`);

    totalValueGroup
      .selectAll('text')
      .data(dataset)
      .enter()
      .append('text')
      .text((d) => formatNumber([d[xKey]]))
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .style('font-size', '.8rem')
      .attr('text-anchor', 'start')
      .attr('alignment-baseline', 'middle')
      .attr('x', (chartData) => `${xScale(chartData[xKey])}%`)
      .attr('y', (chartData) => yScale(chartData[labelKey]));
  }
}
