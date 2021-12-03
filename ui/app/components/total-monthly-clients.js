import Component from '@glimmer/component';
import { action } from '@ember/object';
import { select } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { stack } from 'd3-shape';

/**
 * @module TotalMonthlyClients
 * TotalMonthlyClients components are used to...
 *
 * @example
 * ```js
 * <TotalMonthlyClients @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

const DATA = [
  { month: 'January', directEntities: 150, nonDirectTokens: 22 },
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

export default class TotalMonthlyClients extends Component {
  @action
  registerListner(element) {
    let stackFunction = stack().keys(['directEntities', 'nonDirectTokens']);
    let stackedData = stackFunction(DATA);
    stackedData.forEach(item => console.log(item, 'here'));
    let countArray = DATA.map(month => month.directEntities); // change to combined
    let yScale = scaleLinear()
      // .domain([0, Math.max(...countArray)])
      .domain([0, 1000]) // TODO calculate high of total combined
      .range([0, 250]); // 250 is the height of the chart
    let xScale = scaleBand()
      // .domain(DATA.map(month => month.month))
      .domain([DATA[0].month, DATA[11].month]) // jan to dec
      .range([0, 500])
      .paddingInner(0.85);

    let chartSvg = select(element);
    chartSvg.attr('width', '500px');
    chartSvg.attr('height', '250px');

    // ARG STOP HERE

    chartSvg
      .selectAll('rect')
      .data(stackedData)
      .enter()
      .append('rect')
      // .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);
      .attr('width', xScale.bandwidth())
      .attr('height', 200)
      // .attr('height', month => yScale(month.directEntities))
      .attr('x', month => xScale(month.month)) // offset each bar.
      .attr('y', month => 250 - yScale(month.directEntities)); // flips it right side up.
  }
}
