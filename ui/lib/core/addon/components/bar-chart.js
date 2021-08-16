/**
 * @module BarChart
 * BarChart components are used to...
 *
 * @example
 * ```js
 * <BarChart @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
// import Component from '@glimmer/component';
import layout from '../templates/components/bar-chart';
import { setComponentTemplate } from '@ember/component';
import { select } from 'd3-selection';
import { action } from '@ember/object';
import { scaleLinear, scaleBand } from 'd3-scale';
class BarChart extends Component {
  dataArray = [
    { namespace: 'top-namespace', count: 87 },
    { namespace: 'namespace2', count: 43 },
    { namespace: 'longnamenamespace', count: 23 },
    { namespace: 'path/to/namespace', count: 14 },
  ];

  calculateMaxValue(dataset) {
    let counts = dataset.map(data => data.count);
    return Math.max(...counts);
  }

  @action
  renderBarChart() {
    let xScale = scaleLinear()
      .domain([0, this.calculateMaxValue(this.dataArray)]) // min and max values of dataset
      .range([0, 150]);

    let svg = select('#bar-chart');
    svg
      .selectAll('rect')
      .data(this.dataArray)
      .enter()
      .append('rect')
      .attr('width', namespace => xScale(namespace.count))
      .attr('height', 6)
      .attr('y', (namespace, index) => index * 20);
  }
}

export default setComponentTemplate(layout, BarChart);
