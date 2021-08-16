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
import { scaleLinear } from 'd3-scale';
class BarChart extends Component {
  dataArray = [
    { namespace: 'top-namespace', count: 1512 },
    { namespace: 'namespace2', count: 1300 },
    { namespace: 'longnamenamespace', count: 1200 },
    { namespace: 'anothernamespace', count: 1004 },
    { namespace: 'namespacesomething', count: 950 },
    { namespace: 'namespace5', count: 800 },
    { namespace: 'namespace', count: 700 },
    { namespace: 'namespace999', count: 650 },
    { namespace: 'name-space', count: 600 },
    { namespace: 'path/to/namespace', count: 300 },
  ];

  calculateMaxValue(dataset) {
    let counts = dataset.map(data => data.count);
    // turns array of counts into an argument list
    return Math.max(...counts);
  }

  @action
  renderBarChart() {
    let xScale = scaleLinear()
      .domain([0, this.calculateMaxValue(this.dataArray)]) // min and max values of dataset
      .range([0, 650]);

    let svg = select('#bar-chart');
    svg
      .selectAll('rect')
      .data(this.dataArray)
      .enter()
      .append('rect')
      .attr('width', namespace => xScale(namespace.count))
      .attr('height', 6)
      .attr('y', (namespace, index) => index * 20)
      .attr('fill', '#BFD4FF');
  }
}

export default setComponentTemplate(layout, BarChart);
