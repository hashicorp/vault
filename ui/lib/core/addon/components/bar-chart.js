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
import { select } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { max } from 'd3-array';
import { stack } from 'd3-shape';
import { axisLeft } from 'd3-axis';

const BAR_THICKNESS = 6; // bar thickness in pixels;
const BAR_SPACING = 20; // spacing between bars in pixels
const CHART_MARGIN = { top: 0, right: 24, bottom: 26, left: 137 };

class BarChart extends Component {
  // make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
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
    let stackFunction = stack().keys(['count', 'unique']);
    let stackedData = stackFunction(dataset);

    let xScale = scaleLinear()
      .domain([0, max(dataset, d => d.count + d.unique)]) // min and max values of dataset
      .range([0, 100]); // range is 0-100%

    let yScale = scaleBand()
      .domain(dataset.map(d => d.label))
      .range([0, 193]);

    let yAxis = axisLeft(yScale);

    let svg = select(element);
    let colors = ['#BFD4FF', '#8AB1FF'];
    // add a group for each row of data
    let groups = svg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => colors[i]);

    yAxis(groups.append('g'));
    // add a rect for each data value
    let rects = groups
      .selectAll('rect')
      .data(d => d)
      .enter()
      .append('rect')
      .attr('width', value => `${xScale(value[1] - value[0])}%`)
      .attr('height', BAR_THICKNESS)
      .attr('x', value => `${xScale(value[0])}%`)
      .attr('y', (label, index) => index * BAR_SPACING);
  }
}

export default setComponentTemplate(layout, BarChart);
