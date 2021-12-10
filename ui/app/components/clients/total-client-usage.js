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
  { month: 'January', directEntities: 5000, nonDirectTokens: 22 },
  { month: 'February', directEntities: 1500, nonDirectTokens: 22 },
  { month: 'March', directEntities: 1550, nonDirectTokens: 25 },
  { month: 'April', directEntities: 1550, nonDirectTokens: 229 },
  { month: 'May', directEntities: 1560, nonDirectTokens: 24 },
  { month: 'June', directEntities: 1570, nonDirectTokens: 42 },
  { month: 'July', directEntities: 1580, nonDirectTokens: 12 },
  { month: 'August', directEntities: 1610, nonDirectTokens: 1 },
  { month: 'September', directEntities: 1900, nonDirectTokens: 222 },
  { month: 'October', directEntities: 2500, nonDirectTokens: 66 },
  { month: 'November', directEntities: 3000, nonDirectTokens: 32 },
  { month: 'December', directEntities: 6000, nonDirectTokens: 202 },
];

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#1563FF', '#8AB1FF'];
const BACKGROUND_BAR_COLOR = '#EBEEF2';

const CHART_MARGIN = { top: 10, right: 0, bottom: 0, left: 5 }; // makes space for y-axis legend
const LINE_HEIGHT = 30; // each bar w/ padding is 24 pixels thick

export default class TotalClientUsage extends Component {
  @tracked tooltipTarget = '#wtf';
  @tracked hoveredLabel = 'init';
  @tracked trackingTest = 0;
  @action
  registerListener(element) {
    // Define the chart
    let chartSvg = select(element);
    chartSvg.attr('width', '100%');
    chartSvg.attr('height', '100%');
    chartSvg.attr('viewBox', `0 0 1000 ${(DATA.length + 1) * LINE_HEIGHT}`);

    let stackFunction = stack().keys(['directEntities', 'nonDirectTokens']);
    let stackedData = stackFunction(DATA);

    let yScale = scaleLinear()
      .domain([0, max(stackedData[1].map(d => d.data.directEntities + d.data.nonDirectTokens))]) // TODO will need to recalculate when you get the data
      .range([98, 0]);

    let xScale = scaleBand()
      .domain(DATA.map(month => month.month))
      .range([0, 100]) // unsure about this 100 range?
      .paddingInner(0.85);

    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .attr('transform', `translate(${CHART_MARGIN.left}})`)
      .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);

    groups
      .selectAll('rect')
      .data(stackedData => stackedData)
      .enter()
      .append('rect')
      .attr('width', `${xScale.bandwidth()}%`)
      .attr('height', data => `${100 - yScale(data[1])}%`)
      .attr('x', data => `${xScale(data.data.month)}%`)
      .attr('y', data => `${yScale(data[0]) + yScale(data[1]) - 102}%`) // subtract higher than 100% to give space for x axis ticks
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left})`);

    // creating wider area for tooltip hover
    let tooltipRect = chartSvg
      .selectAll('.tooltip-rect')
      .data(DATA)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'tooltip-rect')
      .attr('height', '100%')
      .attr('width', '50px') // three times width
      .attr('y', '0') // start at bottom
      .attr('x', data => `${xScale(data.month) - 1.1}%`) // not data.data because this is not stacked data
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

    let xScaleNotPercent = scaleBand()
      .domain(DATA.map(month => month.month))
      .range([0, 1200]) //hard coded width because 100 before with percents 0 to width
      .paddingInner(0.85);

    let yScaleNotPercent = scaleLinear()
      .domain([0, max(stackedData[1].map(d => d.data.directEntities + d.data.nonDirectTokens))]) // TODO calculate high of total combined
      .range([800, 0]); // height to zero (it's inverted), TODO make 98 as a percent instead of a fixed number

    let xAxis = axisBottom(xScaleNotPercent).tickSize(1);
    xAxis(chartSvg.append('g').attr('transform', `translate(${CHART_MARGIN.left},375)`));

    // Reference for tickFormat https://www.youtube.com/watch?v=c3MCROTNN8g
    let yAxisTickFormat = number => format('.1s')(number).replace('G', 'B'); // for billions to replace G with B.

    let yAxis = axisLeft(yScaleNotPercent).tickFormat(yAxisTickFormat); // format number as 8k or 8M
    yAxis(chartSvg.append('g').attr('transform', `translate(${CHART_MARGIN.left},0)`));
  }

  // ARG TODO rename
  @action closeMe() {
    this.tooltipTarget = null;
  }
}
