import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { max } from 'd3-array';
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
  { month: 'March', directEntities: 700, nonEntityTokens: 125, total: 825 },
  { month: 'April', directEntities: 1550, nonEntityTokens: 229, total: 1779 },
  { month: 'May', directEntities: 1560, nonEntityTokens: 124, total: 1684 },
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

const AXES_MARGIN = { yLeft: 0, xLeft: 11, yDown: -40, xDown: 265 }; // makes space for y-axis legend
const TRANSLATE = { none: 0, right: 11, down: -30 };
const CHART_HEIGHT = 300;
export default class TotalClientUsage extends Component {
  @tracked tooltipTarget = '#wtf';
  @tracked hoveredLabel = 'init';
  @tracked trackingTest = 0;
  @action
  registerListener(element) {
    // Define the chart
    let dataset = DATA; // will be data passed in as argument

    let stackFunction = stack().keys(['directEntities', 'nonEntityTokens']);
    let stackedData = stackFunction(dataset);

    let yScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([0, 80]); // don't want 100% because will cut off

    let xScale = scaleBand()
      .domain(dataset.map(d => d.month))
      .range([0, 700]) // set width to fix number of pixels
      .paddingInner(0.85);

    let chartSvg = select(element);

    chartSvg.attr('viewBox', `0 0 725 291`); // set aspect ratio

    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .attr('transform', `translate(${TRANSLATE.right}, ${TRANSLATE.down})`) //
      .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);

    groups
      .selectAll('rect')
      .data(stackedData => stackedData)
      .enter()
      .append('rect')
      .attr('width', '7px')
      .attr('height', stackedData => `${yScale(stackedData[1] - stackedData[0])}%`)
      .attr('x', ({ data }) => xScale(data.month)) // uses destructuring because was data.data.month
      .attr('y', data => `${100 - yScale(data[1])}%`); // subtract higher than 100% to give space for x axis ticks

    // MAKE AXES //
    let yAxisScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))]) // TODO will need to recalculate when you get the data
      .range([CHART_HEIGHT, 0]);

    // Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
    let formatNumbers = number => format('.1s')(number).replace('G', 'B'); // for billions to replace G with B.

    // customize y-axis
    let yAxis = axisLeft(yAxisScale)
      .tickSize(5)
      .ticks(6)
      .tickFormat(formatNumbers);

    yAxis(chartSvg.append('g').attr('transform', `translate(${AXES_MARGIN.yLeft}, ${AXES_MARGIN.yDown})`));

    let xAxisGenerator = axisBottom(xScale);
    let xAxis = chartSvg.append('g').call(xAxisGenerator);

    xAxis.attr('transform', `translate(${AXES_MARGIN.xLeft}, ${AXES_MARGIN.xDown})`);

    chartSvg.selectAll('.domain').remove(); // remove domain lines

    // creating wider area for tooltip hover
    let greyBars = chartSvg.append('g').attr('transform', `translate(${TRANSLATE.none}, ${TRANSLATE.down})`);

    let tooltipRect = greyBars
      .selectAll('.tooltip-rect')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'tooltip-rect')
      .attr('height', '100%')
      .attr('width', '30px') // three times width
      .attr('y', '0') // start at bottom
      .attr('x', data => xScale(data.month)) // not data.data because this is not stacked data
      .style('fill', `${BACKGROUND_BAR_COLOR}`)
      .style('opacity', '1')
      .style('mix-blend-mode', 'multiply');

    // for tooltip
    tooltipRect.on('mouseover', data => {
      this.hoveredLabel = data.month;
      let node = chartSvg
        .selectAll('rect.tooltip-rect')
        .filter(data => data.month === this.hoveredLabel)
        .node();
      this.tooltipTarget = node; // grab the second node from the list of 24 rects
    });

    // axis lines

    // let xScaleNotPercent = scaleBand()
    //   .domain(DATA.map(month => month.month))
    //   .range([0, 1200]) //hard coded width because 100 before with percents 0 to width
    //   .paddingInner(0.85);

    // let yScaleNotPercent = scaleLinear()
    //   .domain([0, max(stackedData[1].map(d => d.data.directEntities + d.data.nonDirectTokens))]) // TODO calculate high of total combined
    //   .range([800, 0]); // height to zero (it's inverted), TODO make 98 as a percent instead of a fixed number

    // let xAxis = axisBottom(xScaleNotPercent).tickSize(1);
    // xAxis(chartSvg.append('g').attr('transform', `translate(${AXES_MARGIN.left},375)`));

    // let yAxis = axisLeft(yScaleNotPercent).tickFormat(yAxisTickFormat); // format number as 8k or 8M
    // yAxis(chartSvg.append('g').attr('transform', `translate(${AXES_MARGIN.left},0)`));
  }

  // ARG TODO rename
  @action closeMe() {
    this.tooltipTarget = null;
  }
}
