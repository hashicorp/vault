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
import { max } from 'd3-array';
import { stack as d3Layout } from 'd3-shape';

const BAR_HEIGHT = 6; // bar height in pixels;
const BAR_SPACING = 20; // bar spacing in pixel
class BarChart extends Component {
  dataset = [
    { label: 'top-namespace', count: 1512, unique: 300 },
    { label: 'namespace2', count: 1300, unique: 250 },
    { label: 'longnamenamespace', count: 1200, unique: 200 },
    { label: 'anothernamespace', count: 1004, unique: 150 },
    { label: 'namespacesomething', count: 950, unique: 100 },
    { label: 'namespace5', count: 800, unique: 75 },
    { label: 'namespace', count: 700, unique: 50 },
    { label: 'namespace999', count: 650, unique: 40 },
    { label: 'name-space', count: 600, unique: 20 },
    { label: 'path/to/namespace', count: 300, unique: 10 },
  ];

  @action
  renderBarChart(element) {
    let dataset = this.dataset;
    let stack = d3Layout().keys(['count', 'unique']);

    let stackedData = stack(dataset);
    console.log(stackedData, 'stackedData');

    let xScale = scaleLinear()
      .domain([0, max(this.dataset, d => d.count + d.unique)]) // min and max values of dataset
      .range([0, 100]); // range is from 0-100%

    let svg = select(element);

    let colors = ['#BFD4FF', '#8AB1FF'];
    // Add a group for each row of data
    let groups = svg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .style('fill', function(d, i) {
        return colors[i];
      });
    // Add a rect for each data value
    let rects = groups
      .selectAll('rect')
      .data(function(d) {
        return d;
      })
      .enter()
      .append('rect')
      .attr('width', dataValue => `${xScale(dataValue[1] - dataValue[0])}%`)
      .attr('x', dataValue => `${xScale(dataValue[0])}%`)
      .attr('height', BAR_HEIGHT)
      .attr('y', (label, index) => index * BAR_SPACING);

    // svg
    //   .selectAll('rect')
    //   .data(this.dataset)
    //   .enter()
    //   .append('rect')
    //   .attr('width', label => `${xScale(label.count)}%`)
    //   .attr('height', BAR_HEIGHT)
    //   .attr('y', (label, index) => index * BAR_SPACING)
    //   .attr('fill', '#BFD4FF');
  }
}

export default setComponentTemplate(layout, BarChart);
