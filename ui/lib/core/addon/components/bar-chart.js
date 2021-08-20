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
import { format } from 'd3-format';

const BAR_THICKNESS = 6; // bar thickness in pixels;
const BAR_SPACING = 20; // spacing between bars in pixels
const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 };
const BAR_COLORS = ['#BFD4FF', '#8AB1FF'];
class BarChart extends Component {
  // make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
  variableA = 'Active Direct Tokens';
  variableB = 'Unique Entities';

  dataset = [
    { label: 'top-namespace', count: 1212, unique: 300 },
    { label: 'namespace2', count: 650, unique: 550 },
    { label: 'longnamenamespace', count: 200, unique: 1000 },
    { label: 'anothernamespace', count: 400, unique: 550 },
    { label: 'namespacesomething', count: 400, unique: 400 },
    { label: 'namespace5', count: 800, unique: 300 },
    { label: 'namespace', count: 400, unique: 300 },
    { label: 'namespace999', count: 350, unique: 250 },
    { label: 'name-space', count: 450, unique: 200 },
    { label: 'path/to/namespace', count: 200, unique: 100 },
  ];

  @action
  renderBarChart(element) {
    let dataset = this.dataset.sort((a, b) => a.count + a.unique - (b.count + b.unique)).reverse();
    let stackFunction = stack().keys(['count', 'unique']);
    let stackedData = stackFunction(dataset); // returns an array of coordinates for each rectangle group, first group is for counts (left), second for unique (right)

    let xScale = scaleLinear()
      .domain([0, max(dataset, d => d.count + d.unique)]) // min and max values of dataset
      .range([0, 75]); // range in percent (30% reserved for margins)

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
      .attr('width', data => `${xScale(data[1] - data[0] - 6)}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', data => `${xScale(data[0])}%`)
      .attr('y', ({ data }) => yScale(data.label))
      .attr('rx', 3)
      .attr('ry', 3);

    svg.attr('height', (dataset.length + 1) * 24);

    let totalNumbers = [];
    stackedData[1].forEach(e => {
      let n = e[1];
      totalNumbers.push(n);
    });

    let textData = [];
    rects.each(function(d, i) {
      // if (d[0] !== 0){
      let textDatum = {
        text: totalNumbers[i],
        x: parseFloat(select(this).attr('width')) + parseFloat(select(this).attr('x')),
        y: parseFloat(select(this).attr('y')) + parseFloat(select(this).attr('height')),
      };
      textData.push(textDatum);
      // };
    });

    let text = groups
      .selectAll('text')
      .data(textData)
      .enter()
      .append('text')
      .text(d => d.text)
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .attr('text-anchor', 'start')
      .attr('y', d => {
        return `${d.y}`;
      })
      .attr('x', d => {
        let translateRight = d.x + 1;
        return `${translateRight}%`;
      });

    // style here? or in .css file
    groups.selectAll('.domain, .tick line').remove();

    let legend = select('.legend');
    legend
      .append('circle')
      .attr('cx', '60%')
      .attr('cy', '20%')
      .attr('r', 6)
      .style('fill', '#BFD4FF');
    legend
      .append('text')
      .attr('x', '62%')
      .attr('y', '20%')
      .text(`${this.variableA}`)
      .style('font-size', '15px')
      .attr('alignment-baseline', 'middle');
    legend
      .append('circle')
      .attr('cx', '83%')
      .attr('cy', '20%')
      .attr('r', 6)
      .style('fill', '#8AB1FF');
    legend
      .append('text')
      .attr('x', '85%')
      .attr('y', '20%')
      .text(`${this.variableB}`)
      .style('font-size', '15px')
      .attr('alignment-baseline', 'middle');
  }
}

export default setComponentTemplate(layout, BarChart);
