import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
// eslint-disable-next-line no-unused-vars
import { select, selectAll, node } from 'd3-selection';
import { axisLeft, axisBottom } from 'd3-axis';
import { scaleLinear, scaleBand } from 'd3-scale';
import { stack } from 'd3-shape';
import {
  GREY,
  LIGHT_AND_DARK_BLUE,
  SVG_DIMENSIONS,
  TRANSLATE,
  formatNumbers,
} from '../../utils/chart-helpers';

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
 */

export default class VerticalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipTotal = '';
  @tracked uniqueEntities = '';
  @tracked nonEntityTokens = '';

  get chartLegend() {
    return this.args.chartLegend;
  }

  @action
  registerListener(element, args) {
    let dataset = args[0];
    let stackFunction = stack().keys(this.chartLegend.map((l) => l.key));
    let stackedData = stackFunction(dataset);
    let chartSvg = select(element);
    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    // DEFINE DATA BAR SCALES
    let yScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d.clients))]) // TODO will need to recalculate when you get the data
      .range([0, 100])
      .nice();

    let xScale = scaleBand()
      .domain(dataset.map((d) => d.month))
      .range([0, SVG_DIMENSIONS.width]) // set width to fix number of pixels
      .paddingInner(0.85);

    let dataBars = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .style('fill', (d, i) => LIGHT_AND_DARK_BLUE[i]);

    dataBars
      .selectAll('rect')
      .data((stackedData) => stackedData)
      .enter()
      .append('rect')
      .attr('width', '7px')
      .attr('class', 'data-bar')
      .attr('height', (stackedData) => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data.month)) // uses destructuring because was data.data.month
      .attr('y', (data) => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

    // MAKE AXES //
    let yAxisScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d.clients))]) // TODO will need to recalculate when you get the data
      .range([`${SVG_DIMENSIONS.height}`, 0])
      .nice();

    let yAxis = axisLeft(yAxisScale)
      .ticks(7)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width)
      .tickFormat(formatNumbers);

    let xAxis = axisBottom(xScale).tickSize(0);

    yAxis(chartSvg.append('g'));
    xAxis(chartSvg.append('g').attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`));

    chartSvg.selectAll('.domain').remove(); // remove domain lines

    // WIDER SELECTION AREA FOR TOOLTIP HOVER
    let greyBars = chartSvg
      .append('g')
      .attr('transform', `translate(${TRANSLATE.left})`)
      .style('fill', `${GREY}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let tooltipRect = greyBars
      .selectAll('rect')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'tooltip-rect')
      .attr('height', '100%')
      .attr('width', '30px') // three times width
      .attr('y', '0') // start at bottom
      .attr('x', (data) => xScale(data.month)); // not data.data because this is not stacked data

    // MOUSE EVENT FOR TOOLTIP
    tooltipRect.on('mouseover', (data) => {
      let hoveredMonth = data.month;
      this.tooltipTotal = `${data.clients} ${data.new_clients ? 'total' : 'new'} clients`;
      this.uniqueEntities = `${data.entity_clients} unique entities`;
      this.nonEntityTokens = `${data.non_entity_clients} non-entity tokens`;
      // let node = chartSvg
      //   .selectAll('rect.tooltip-rect')
      //   .filter(data => data.month === this.hoveredLabel)
      //   .node();
      let node = chartSvg
        .selectAll('rect.data-bar')
        // filter for the top data bar (so y-coord !== 0) with matching month
        .filter((data) => data[0] !== 0 && data.data.month === hoveredMonth)
        .node();
      this.tooltipTarget = node; // grab the node from the list of rects
    });
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }
}
