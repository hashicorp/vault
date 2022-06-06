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

/**
 * @module LineChart
 * LineChart components are used to display data in a line plot with accompanying tooltip
 *
 * @example
 * ```js
 * <LineChart @dataset={dataset} />
 * ```
 * @param {string} xKey - string denoting key for x-axis data (data[xKey]) of dataset
 * @param {string} yKey - string denoting key for y-axis data (data[yKey]) of dataset
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

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, [chartData]) {
    const dataset = chartData;
    const upgradeData = [];
    if (this.args.upgradeData) {
      this.args.upgradeData.forEach((versionData) =>
        upgradeData.push({ month: parseAPITimestamp(versionData.timestampInstalled, 'M/yy'), ...versionData })
      );
    }
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
      return upgradeData.find((upgrade) => upgrade[this.xKey] === datum[this.xKey]);
    };

    // VERSION UPGRADE INDICATOR
    chartSvg
      .append('g')
      .selectAll('circle')
      .data(filteredData)
      .enter()
      .append('circle')
      .attr('class', 'upgrade-circle')
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
      this.tooltipTotal = data[this.yKey] + ' total clients';
      this.tooltipNew = (data?.new_clients[this.yKey] || '0') + ' new clients';
      this.tooltipUpgradeText = '';
      let upgradeInfo = findUpgradeData(data);
      if (upgradeInfo) {
        let { id, previousVersion } = upgradeInfo;
        this.tooltipUpgradeText = `Vault was upgraded 
        ${previousVersion ? 'from ' + previousVersion : ''} to ${id}`;
      }

      let node = hoverCircles.filter((plot) => plot[this.xKey] === data[this.xKey]).node();
      this.tooltipTarget = node;
    });
  }
}
