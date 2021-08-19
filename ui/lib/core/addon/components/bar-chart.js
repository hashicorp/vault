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
const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 };
const BAR_COLORS = ['#BFD4FF', '#8AB1FF'];
class BarChart extends Component {
  // make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
  dataset = [
    { label: 'top-namespace', count: 1512, unique: 300 },
    { label: 'namespace2', count: 1300, unique: 250 },
    { label: 'longnamenamespace', count: 1200, unique: 200 },
    { label: 'anothernamespace', count: 1004, unique: 150 },
    { label: 'namespacesomething', count: 950, unique: 100 },
    { label: 'namespace5', count: 800, unique: 75 },
    { label: 'namespace', count: 400, unique: 300 },
    { label: 'namespace999', count: 650, unique: 40 },
    { label: 'name-space', count: 600, unique: 20 },
    { label: 'path/to/namespace', count: 300, unique: 499 },
  ];

  @action
  renderBarChart(element) {
    let dataset = this.dataset.sort((a, b) => a.count + a.unique - (b.count + b.unique)).reverse();

    let stackFunction = stack().keys(['count', 'unique']);
    let stackedData = stackFunction(dataset);

    let xScale = scaleLinear()
      .domain([0, max(dataset, d => d.count + d.unique)]) // min and max values of dataset
      .range([0, 70]); // range in percent (30% reserved for margins)

    let yScale = scaleBand()
      .domain(dataset.map(d => d.label))
      .range([0, dataset.length * 24]) // each bar element has a thickness (bar + padding) of 24 pixels
      // paddingInner takes a number between 0 and 1
      // it tells the scale the percent of the total width it should reserve for white space between bars
      .paddingInner(0.765);

    let svg = select(element);
    // add a group for each row of data

    let groups = svg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => BAR_COLORS[i]);

    // yAxis legend
    let yAxis = axisLeft(yScale);
    yAxis(groups.append('g'));

    let rects = groups
      .selectAll('rect')
      .data(d => d)
      .enter()
      .append('rect')
      .attr('width', data => `${xScale(data[1] - data[0] - 10)}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', data => `${xScale(data[0])}%`)
      .attr('y', data => yScale(data.data.label));

    // style here? or in .css file
    groups.selectAll('.domain, .tick line').remove();
  }
}

export default setComponentTemplate(layout, BarChart);
