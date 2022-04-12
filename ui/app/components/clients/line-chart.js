import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
// eslint-disable-next-line no-unused-vars
import { select, selectAll, node } from 'd3-selection';
import { axisLeft, axisBottom } from 'd3-axis';
import { scaleLinear, scalePoint } from 'd3-scale';
import { line } from 'd3-shape';
import { LIGHT_AND_DARK_BLUE, SVG_DIMENSIONS, formatNumbers } from '../../utils/chart-helpers';

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

  get yKey() {
    return this.args.yKey || 'total';
  }

  get xKey() {
    return this.args.xKey || 'month';
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, args) {
    const dataset = args[0];
    const chartSvg = select(element);
    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    // DEFINE AXES SCALES
    const yScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d[this.yKey]))])
      .range([0, 100])
      .nice();

    const yAxisScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d[this.yKey]))])
      .range([SVG_DIMENSIONS.height, 0])
      .nice();

    const xScale = scalePoint() // use scaleTime()?
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

    yAxis(chartSvg.append('g'));
    xAxis(chartSvg.append('g').attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`));

    chartSvg.selectAll('.domain').remove();

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
      .attr('d', lineGenerator(dataset));

    // LINE PLOTS (CIRCLES)
    chartSvg
      .append('g')
      .selectAll('circle')
      .data(dataset)
      .enter()
      .append('circle')
      .attr('class', 'data-plot')
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
      .data(dataset)
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
      // TODO: how to genericize this?
      this.tooltipMonth = data[this.xKey];
      this.tooltipTotal = `${data[this.yKey]} total clients`;
      this.tooltipNew = `${data?.new_clients[this.yKey]} new clients`;
      let node = hoverCircles.filter((plot) => plot[this.xKey] === data[this.xKey]).node();
      this.tooltipTarget = node;
    });
  }
}
