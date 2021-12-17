import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
// eslint-disable-next-line no-unused-vars
import { select, selectAll, node } from 'd3-selection';
import { axisLeft, axisBottom } from 'd3-axis';
import { scaleLinear, scaleBand } from 'd3-scale';
import { format } from 'd3-format';
import { stack } from 'd3-shape';

/**
 * ARG TODO fill out
 * @module TotalClientUsage
 * TotalClientUsage components are used to...
 *
 * @example
 * ```js
 * <TotalClientUsage @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

// ARG TODO pull in data
const DATA = [
  { month: 'January', directEntities: 1000, nonEntityTokens: 322, total: 1322 },
  { month: 'February', directEntities: 1500, nonEntityTokens: 122, total: 1622 },
  { month: 'March', directEntities: 4300, nonEntityTokens: 700, total: 5000 },
  { month: 'April', directEntities: 1550, nonEntityTokens: 229, total: 1779 },
  { month: 'May', directEntities: 5560, nonEntityTokens: 124, total: 5684 },
  { month: 'June', directEntities: 1570, nonEntityTokens: 142, total: 1712 },
  { month: 'July', directEntities: 300, nonEntityTokens: 112, total: 412 },
  { month: 'August', directEntities: 1610, nonEntityTokens: 130, total: 1740 },
  { month: 'September', directEntities: 1900, nonEntityTokens: 222, total: 2122 },
  { month: 'October', directEntities: 500, nonEntityTokens: 166, total: 666 },
  { month: 'November', directEntities: 480, nonEntityTokens: 132, total: 612 },
  { month: 'December', directEntities: 980, nonEntityTokens: 202, total: 1182 },
];

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#8AB1FF', '#1563FF'];
const BACKGROUND_BAR_COLOR = '#EBEEF2';

// TRANSLATIONS:
const TRANSLATE = { left: -11 };
const SVG_DIMENSIONS = { height: 190, width: 500 };

export default class TotalClientUsage extends Component {
  @tracked tooltipTarget = '';
  @tracked hoveredLabel = '';
  @tracked trackingTest = 0;

  @action
  registerListener(element) {
    // Define the chart
    let dataset = DATA; // will be data passed in as argument

    let stackFunction = stack().keys(['directEntities', 'nonEntityTokens']);
    let stackedData = stackFunction(dataset);

    let yScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([0, 100]);

    let xScale = scaleBand()
      .domain(dataset.map(d => d.month))
      .range([0, SVG_DIMENSIONS.width]) // set width to fix number of pixels
      .paddingInner(0.85);

    let chartSvg = select(element);

    chartSvg.attr('viewBox', `-50 20 600 ${SVG_DIMENSIONS.height}`); // set svg dimensions

    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);

    groups
      .selectAll('rect')
      .data(stackedData => stackedData)
      .enter()
      .append('rect')
      .attr('width', '7px')
      .attr('class', 'data-bar')
      .attr('height', stackedData => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data.month)) // uses destructuring because was data.data.month
      .attr('y', data => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

    // MAKE AXES //
    let yAxisScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([`${SVG_DIMENSIONS.height}`, 0]);

    // Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
    let formatNumbers = number => format('.2s')(number).replace('G', 'B'); // for billions to replace G with B.

    // customize y-axis
    let yAxis = axisLeft(yAxisScale)
      .ticks(7)
      .tickPadding(10)
      .tickSizeInner(-SVG_DIMENSIONS.width) // makes grid lines correct length
      .tickFormat(formatNumbers);

    yAxis(chartSvg.append('g'));

    // customize x-axis
    let xAxisGenerator = axisBottom(xScale).tickSize(0);
    let xAxis = chartSvg.append('g').call(xAxisGenerator);

    xAxis.attr('transform', `translate(0, ${SVG_DIMENSIONS.height + 10})`);

    chartSvg.selectAll('.domain').remove(); // remove domain lines

    // creating wider area for tooltip hover
    let greyBars = chartSvg
      .append('g')
      .attr('transform', `translate(${TRANSLATE.left})`)
      .style('fill', `${BACKGROUND_BAR_COLOR}`)
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
      .attr('x', data => xScale(data.month)); // not data.data because this is not stacked data

    // for tooltip
    tooltipRect.on('mouseover', data => {
      this.hoveredLabel = data.month;
      // let node = chartSvg
      //   .selectAll('rect.tooltip-rect')
      //   .filter(data => data.month === this.hoveredLabel)
      //   .node();
      let node = chartSvg
        .selectAll('rect.data-bar')
        // filter for the top data bar (so y-coord !== 0) with matching month
        .filter(data => data[0] !== 0 && data.data.month === this.hoveredLabel)
        .node();
      this.tooltipTarget = node; // grab the node from the list of rects
    });
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }
}
