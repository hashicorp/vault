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
 * @param {array} dataset - dataset is an array of objects
 */

export default class LineChart extends Component {
  // TODO CMB make just one tracked variable tooltipText?
  @tracked tooltipTarget = '';
  @tracked tooltipMonth = '';
  @tracked tooltipTotal = '';
  @tracked tooltipNew = '';

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, args) {
    let dataset = args[0];
    let chartSvg = select(element);
    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    // DEFINE AXES SCALES
    let yScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d.total))])
      .range([0, 100]);

    let yAxisScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d.total))])
      .range([SVG_DIMENSIONS.height, 0]);

    let xScale = scalePoint() // use scaleTime()?
      .domain(dataset.map((d) => d.month))
      .range([0, SVG_DIMENSIONS.width])
      .padding(0.2);

    // CUSTOMIZE AND APPEND AXES
    let yAxis = axisLeft(yAxisScale)
      .ticks(7)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width) // makes grid lines length of svg
      .tickFormat(formatNumbers);

    let xAxis = axisBottom(xScale).tickSize(0);

    yAxis(chartSvg.append('g'));
    xAxis(chartSvg.append('g').attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`));

    chartSvg.selectAll('.domain').remove();

    // PATH BETWEEN PLOT POINTS
    let lineGenerator = line()
      .x((d) => xScale(d.month))
      .y((d) => yAxisScale(d.total));

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
      .attr('cy', (d) => `${100 - yScale(d.total)}%`)
      .attr('cx', (d) => xScale(d.month))
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
      .attr('cy', (d) => `${100 - yScale(d.total)}%`)
      .attr('cx', (d) => xScale(d.month))
      .attr('r', 10);

    let hoverCircles = chartSvg.selectAll('.hover-circle');

    // MOUSE EVENT FOR TOOLTIP
    hoverCircles.on('mouseover', (data) => {
      this.tooltipMonth = data.month;
      this.tooltipTotal = `${data.total} total clients`;
      this.tooltipNew = `${data.new_clients.total} new clients`;
      let node = hoverCircles.filter((plot) => plot.month === data.month).node();
      this.tooltipTarget = node;
    });
  }
}
