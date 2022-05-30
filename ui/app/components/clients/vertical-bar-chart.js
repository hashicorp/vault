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
  GREY,
  LIGHT_AND_DARK_BLUE,
  SVG_DIMENSIONS,
  TRANSLATE,
  formatNumbers,
} from 'vault/utils/chart-helpers';

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
 */

export default class VerticalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipTotal = '';
  @tracked entityClients = '';
  @tracked nonEntityClients = '';

  get chartLegend() {
    return this.args.chartLegend;
  }

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
    const stackFunction = stack().keys(this.chartLegend.map((l) => l.key));
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
      .style('fill', (d, i) => LIGHT_AND_DARK_BLUE[i]);

    dataBars
      .selectAll('rect')
      .data((stackedData) => stackedData)
      .enter()
      .append('rect')
      .attr('width', '7px')
      .attr('class', 'data-bar')
      .attr('height', (stackedData) => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data[this.xKey])) // uses destructuring because was data.data.month
      .attr('y', (data) => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

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

    yAxis(chartSvg.append('g'));
    xAxis(chartSvg.append('g').attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`));

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
      let hoveredMonth = data[this.xKey];
      this.tooltipTotal = `${data[this.yKey]} ${data.new_clients ? 'total' : 'new'} clients`;
      this.entityClients = `${data.entity_clients} entity clients`;
      this.nonEntityClients = `${data.non_entity_clients} non-entity clients`;
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
