/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
// eslint-disable-next-line no-unused-vars
import { select, selectAll, node } from 'd3-selection';
import { axisLeft, axisBottom } from 'd3-axis';
import { scaleLinear, scalePoint } from 'd3-scale';
import { stack } from 'd3-shape';
import {
  BAR_WIDTH,
  GREY,
  BLUE_PALETTE,
  SVG_DIMENSIONS,
  TRANSLATE,
  calculateSum,
  formatNumbers,
} from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';

/**
 * @module VerticalBarChart
 * VerticalBarChart components are used to display stacked data in a vertical bar chart with accompanying tooltip
 *
 * @example
 * ```js
 * <VerticalBarChart @dataset={dataset} @chartLegend={chartLegend} />
 * ```
 * @param {array} dataset - dataset for the chart, must be an array of flattened objects
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {string} xKey - string denoting key for x-axis data (data[xKey]) of dataset
 * @param {string} yKey - string denoting key for y-axis data (data[yKey]) of dataset
 * @param {string} [noDataMessage] - custom empty state message that displays when no dataset is passed to the chart
 */

export default class VerticalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipTotal = '';
  @tracked tooltipStats = [];

  get xKey() {
    return this.args.xKey || 'month';
  }

  get yKey() {
    return this.args.yKey || 'clients';
  }

  @action
  renderChart(element, [chartData]) {
    const dataset = chartData;
    const filteredData = dataset.filter((e) => Object.keys(e).includes('clients')); // months with data will contain a 'clients' key (otherwise only a timestamp)
    const stackFunction = stack().keys(this.args.chartLegend.map((l) => l.key));
    const stackedData = stackFunction(filteredData);
    const chartSvg = select(element);
    const domainMax = max(filteredData.map((d) => d[this.yKey]));

    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    // DEFINE DATA BAR SCALES
    const yScale = scaleLinear().domain([0, domainMax]).range([0, 100]).nice();

    const xScale = scalePoint()
      .domain(dataset.map((d) => d[this.xKey]))
      .range([0, SVG_DIMENSIONS.width]) // set width to fix number of pixels
      .padding(0.2);

    // clear out DOM before appending anything
    chartSvg.selectAll('g').remove().exit().data(stackedData).enter();
    const dataBars = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .style('fill', (d, i) => BLUE_PALETTE[i]);

    dataBars
      .selectAll('rect')
      .data((stackedData) => stackedData)
      .enter()
      .append('rect')
      .attr('width', `${BAR_WIDTH}px`)
      .attr('class', 'data-bar')
      .attr('data-test-vertical-chart', 'data-bar')
      .attr('height', (stackedData) => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data[this.xKey])) // uses destructuring because was data.data.month
      .attr('y', (data) => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

    const tooltipTether = chartSvg
      .append('g')
      .attr('transform', `translate(${BAR_WIDTH / 2})`)
      .attr('data-test-vertical-chart', 'tooltip-tethers')
      .selectAll('circle')
      .data(filteredData)
      .enter()
      .append('circle')
      .style('opacity', '0')
      .attr('cy', (d) => `${100 - yScale(d[this.yKey])}%`)
      .attr('cx', (d) => xScale(d[this.xKey]))
      .attr('r', 1);

    // MAKE AXES //
    const yAxisScale = scaleLinear()
      .domain([0, max(filteredData.map((d) => d[this.yKey]))])
      .range([`${SVG_DIMENSIONS.height}`, 0])
      .nice();

    const yAxis = axisLeft(yAxisScale)
      .ticks(4)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width)
      .tickFormat(formatNumbers);

    const xAxis = axisBottom(xScale).tickSize(0);

    yAxis(chartSvg.append('g').attr('data-test-vertical-chart', 'y-axis-labels'));
    xAxis(
      chartSvg
        .append('g')
        .attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`)
        .attr('data-test-vertical-chart', 'x-axis-labels')
    );

    chartSvg.selectAll('.domain').remove(); // remove domain lines

    // WIDER SELECTION AREA FOR TOOLTIP HOVER
    const greyBars = chartSvg
      .append('g')
      .attr('transform', `translate(${TRANSLATE.left})`)
      .style('fill', `${GREY}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    const tooltipRect = greyBars
      .selectAll('rect')
      .data(filteredData)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'tooltip-rect')
      .attr('height', '100%')
      .attr('width', '30px') // three times width
      .attr('y', '0') // start at bottom
      .attr('x', (data) => xScale(data[this.xKey])); // not data.data because this is not stacked data

    // MOUSE EVENT FOR TOOLTIP
    tooltipRect.on('mouseover', (data) => {
      const hoveredMonth = data[this.xKey];
      const stackedNumbers = []; // accumulates stacked dataset values to calculate total
      this.tooltipStats = []; // clear stats
      this.args.chartLegend.forEach(({ key, label }) => {
        stackedNumbers.push(data[key]);
        // since we're relying on D3 not ember reactivity,
        // pushing directly to this.tooltipStats updates the DOM
        this.tooltipStats.push(`${formatNumber([data[key]])} ${label}`);
      });
      this.tooltipTotal = `${formatNumber([calculateSum(stackedNumbers)])} ${
        data.new_clients ? 'total' : 'new'
      } clients`;
      // filter for the tether point that matches the hoveredMonth
      const hoveredElement = tooltipTether.filter((data) => data.month === hoveredMonth).node();
      this.tooltipTarget = hoveredElement; // grab the node from the list of rects
    });
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }
}
