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

import Component from '@glimmer/component';
import layout from '../templates/components/bar-chart';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { select } from 'd3-selection';
import { scaleLinear } from 'd3-scale';
class BarChart extends Component {
  dataset = [
    { label: 'top-namespace', count: 1512 },
    { label: 'namespace2', count: 1300 },
    { label: 'longnamenamespace', count: 1200 },
    { label: 'anothernamespace', count: 1004 },
    { label: 'namespacesomething', count: 950 },
    { label: 'namespace5', count: 800 },
    { label: 'namespace', count: 700 },
    { label: 'namespace999', count: 650 },
    { label: 'name-space', count: 600 },
    { label: 'path/to/namespace', count: 300 },
  ];

  calculateMaxValue(data) {
    let counts = data.map(data => data.count);
    // turns array of counts into an argument list
    return Math.max(...counts);
  }

  @action
  renderBarChart() {
    let xScale = scaleLinear()
      .domain([0, this.calculateMaxValue(this.dataset)]) // min and max values of dataset
      .range([0, 90]); // bar length will expand to 90% of container

    let svg = select('#bar-chart');
    svg
      .selectAll('rect')
      .data(this.dataset)
      .enter()
      .append('rect')
      .attr('width', label => `${xScale(label.count)}%`)
      .attr('height', 6)
      .attr('y', (label, index) => index * 20)
      .attr('fill', '#BFD4FF');
  }
}

export default setComponentTemplate(layout, BarChart);
