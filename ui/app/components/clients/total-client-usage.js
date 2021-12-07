import Component from '@glimmer/component';
import { action } from '@ember/object';
import { select } from 'd3-selection';
import { scaleLinear, scaleBand, scaleTime } from 'd3-scale';
import { stack } from 'd3-shape';

/**
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

const DATA = [
  { month: 'January', directEntities: 500, nonDirectTokens: 22 },
  { month: 'February', directEntities: 150, nonDirectTokens: 22 },
  { month: 'March', directEntities: 155, nonDirectTokens: 25 },
  { month: 'April', directEntities: 155, nonDirectTokens: 229 },
  { month: 'May', directEntities: 156, nonDirectTokens: 24 },
  { month: 'June', directEntities: 157, nonDirectTokens: 42 },
  { month: 'July', directEntities: 158, nonDirectTokens: 12 },
  { month: 'August', directEntities: 161, nonDirectTokens: 1 },
  { month: 'September', directEntities: 190, nonDirectTokens: 222 },
  { month: 'October', directEntities: 250, nonDirectTokens: 66 },
  { month: 'November', directEntities: 300, nonDirectTokens: 32 },
  { month: 'December', directEntities: 600, nonDirectTokens: 202 },
];

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#BFD4FF', '#8AB1FF'];

export default class TotalClientUsage extends Component {
  @action
  registerListner(element) {
    let stackFunction = stack().keys(['directEntities', 'nonDirectTokens']);
    let stackedData = stackFunction(DATA);
    let countArray = DATA.map(month => month.directEntities); // change to combined
    let yScale = scaleLinear()
      // .domain([0, Math.max(...countArray)])
      .domain([0, 802]) // TODO calculate high of total combined
      .range([100, 0]); // 250 is the height of the chart
    let xScale = scaleBand()
      // .domain(DATA.map(month => month.month))
      .domain(DATA.map(month => month.month)) // needs to iterate through each one.
      .range([0, 100]) // width of the chart as percent
      .paddingInner(0.85);
    let chartSvg = select(element);
    chartSvg.attr('width', '100%');
    chartSvg.attr('height', '100%');

    // ARG STOP HERE
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
      .attr('width', `${xScale.bandwidth()}%`)
      .attr('height', data => `${100 - yScale(data[1])}%`)
      .attr('x', data => `${xScale(data.data.month)}%`)
      .attr('y', data => `${yScale(data[0]) + yScale(data[1]) - 100}%`);
  }
}
