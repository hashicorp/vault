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
import { line } from 'd3-shape';
import {
  LIGHT_AND_DARK_BLUE,
  UPGRADE_WARNING,
  SVG_DIMENSIONS,
  formatNumbers,
} from 'vault/utils/chart-helpers';
import { parseAPITimestamp, formatChartDate } from 'core/utils/date-formatters';
import { formatNumber } from 'core/helpers/format-number';

/**
 * @module LineChart
 * LineChart components are used to display data in a line plot with accompanying tooltip
 *
 * @example
 * <LineChart @dataset={{dataset}} @upgradeData={{this.versionHistory}}/>
 *
 * @param {string} xKey - string denoting key for x-axis data (data[xKey]) of dataset
 * @param {string} yKey - string denoting key for y-axis data (data[yKey]) of dataset
 * @param {array} upgradeData - array of objects containing version history from the /version-history endpoint
 * @param {string} [noDataMessage] - custom empty state message that displays when no dataset is passed to the chart
 */

export default class LineChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipMonth = '';
  @tracked tooltipTotal = '';
  @tracked tooltipNew = '';
  @tracked tooltipUpgradeText = '';

  get yKey() {
    return this.args.yKey || 'clients';
  }

  get xKey() {
    return this.args.xKey || 'month';
  }

  get upgradeData() {
    const upgradeData = this.args.upgradeData;
    if (!upgradeData) return null;
    if (!Array.isArray(upgradeData)) {
      console.debug('upgradeData must be an array of objects containing upgrade history'); // eslint-disable-line
      return null;
    } else if (!Object.keys(upgradeData[0]).includes('timestampInstalled')) {
      // eslint-disable-next-line
      console.debug(
        `upgrade must be an object with the following key names: ['version', 'previousVersion', 'timestampInstalled']`
      );
      return null;
    } else {
      return upgradeData?.map((versionData) => {
        return {
          [this.xKey]: parseAPITimestamp(versionData.timestampInstalled, 'M/yy'),
          ...versionData,
        };
      });
    }
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, [chartData]) {
    const dataset = chartData;
    const filteredData = dataset.filter((e) => Object.keys(e).includes(this.yKey)); // months with data will contain a 'clients' key (otherwise only a timestamp)
    const domainMax = max(filteredData.map((d) => d[this.yKey]));
    const chartSvg = select(element);
    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions
    // clear out DOM before appending anything
    chartSvg.selectAll('g').remove().exit().data(filteredData).enter();

    // DEFINE AXES SCALES
    const yScale = scaleLinear().domain([0, domainMax]).range([0, 100]).nice();
    const yAxisScale = scaleLinear().domain([0, domainMax]).range([SVG_DIMENSIONS.height, 0]).nice();

    // use full dataset (instead of filteredData) so x-axis spans months with and without data
    const xScale = scalePoint()
      .domain(dataset.map((d) => d[this.xKey]))
      .range([0, SVG_DIMENSIONS.width])
      .padding(0.2);

    // CUSTOMIZE AND APPEND AXES
    const yAxis = axisLeft(yAxisScale)
      .ticks(4)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width) // makes grid lines length of svg
      .tickFormat(formatNumbers);

    const xAxis = axisBottom(xScale).tickSize(0);

    yAxis(chartSvg.append('g').attr('data-test-line-chart', 'y-axis-labels'));
    xAxis(
      chartSvg
        .append('g')
        .attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`)
        .attr('data-test-line-chart', 'x-axis-labels')
    );

    chartSvg.selectAll('.domain').remove();

    const findUpgradeData = (datum) => {
      return this.upgradeData
        ? this.upgradeData.find((upgrade) => upgrade[this.xKey] === datum[this.xKey])
        : null;
    };

    // VERSION UPGRADE INDICATOR
    chartSvg
      .append('g')
      .selectAll('circle')
      .data(filteredData)
      .enter()
      .append('circle')
      .attr('class', 'upgrade-circle')
      .attr('data-test-line-chart', (d) => `upgrade-${d[this.xKey]}`)
      .attr('fill', UPGRADE_WARNING)
      .style('opacity', (d) => (findUpgradeData(d) ? '1' : '0'))
      .attr('cy', (d) => `${100 - yScale(d[this.yKey])}%`)
      .attr('cx', (d) => xScale(d[this.xKey]))
      .attr('r', 10);

    // PATH BETWEEN PLOT POINTS
    const lineGenerator = line()
      .x((d) => xScale(d[this.xKey]))
      .y((d) => yAxisScale(d[this.yKey]));

    chartSvg
      .append('g')
      .append('path')
      .attr('fill', 'none')
      .attr('stroke', LIGHT_AND_DARK_BLUE[1])
      .attr('stroke-width', 0.5)
      .attr('d', lineGenerator(filteredData));

    // LINE PLOTS (CIRCLES)
    chartSvg
      .append('g')
      .selectAll('circle')
      .data(filteredData)
      .enter()
      .append('circle')
      .attr('data-test-line-chart', 'plot-point')
      .attr('cy', (d) => `${100 - yScale(d[this.yKey])}%`)
      .attr('cx', (d) => xScale(d[this.xKey]))
      .attr('r', 3.5)
      .attr('fill', LIGHT_AND_DARK_BLUE[0])
      .attr('stroke', LIGHT_AND_DARK_BLUE[1])
      .attr('stroke-width', 1.5);

    // LARGER HOVER CIRCLES
    chartSvg
      .append('g')
      .selectAll('circle')
      .data(filteredData)
      .enter()
      .append('circle')
      .attr('class', 'hover-circle')
      .style('cursor', 'pointer')
      .style('opacity', '0')
      .attr('cy', (d) => `${100 - yScale(d[this.yKey])}%`)
      .attr('cx', (d) => xScale(d[this.xKey]))
      .attr('r', 10);

    const hoverCircles = chartSvg.selectAll('.hover-circle');

    // MOUSE EVENT FOR TOOLTIP
    hoverCircles.on('mouseover', (data) => {
      // TODO: how to generalize this?
      this.tooltipMonth = formatChartDate(data[this.xKey]);
      this.tooltipTotal = formatNumber([data[this.yKey]]) + ' total clients';
      this.tooltipNew = (formatNumber([data?.new_clients[this.yKey]]) || '0') + ' new clients';
      this.tooltipUpgradeText = '';
      const upgradeInfo = findUpgradeData(data);
      if (upgradeInfo) {
        const { version, previousVersion } = upgradeInfo;
        this.tooltipUpgradeText = `Vault was upgraded
        ${previousVersion ? 'from ' + previousVersion : ''} to ${version}`;
      }

      const node = hoverCircles.filter((plot) => plot[this.xKey] === data[this.xKey]).node();
      this.tooltipTarget = node;
    });
  }
}
