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
import { select, event } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { max } from 'd3-array';
import { stack } from 'd3-shape';
import { axisLeft } from 'd3-axis';
import { transition } from 'd3-transition';
import { format } from 'd3-format';

const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 }; // makes space for y-axis legend
const BAR_COLORS = ['#BFD4FF', '#8AB1FF'];
class BarChart extends Component {
  // make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
  variableA = 'Active direct tokens';
  variableB = 'Unique entities';

  dataset = [
    { label: 'top-namespace', count: 1212, unique: 300, total: 1512 },
    { label: 'namespace2', count: 650, unique: 550, total: 1200 },
    { label: 'longnamenamespace', count: 200, unique: 1000, total: 1200 },
    { label: 'anothernamespace', count: 400, unique: 550, total: 1100 },
    { label: 'namespacesomething', count: 400, unique: 400, total: 950 },
    { label: 'namespace5', count: 800, unique: 300, total: 800 },
    { label: 'namespace', count: 400, unique: 300, total: 700 },
    { label: 'namespace999', count: 350, unique: 250, total: 650 },
    { label: 'name-space', count: 450, unique: 200, total: 600 },
    { label: 'path/to/namespace', count: 200, unique: 100, total: 300 },
  ];

  totalCount = this.dataset.reduce((previousValue, currentValue) => previousValue + currentValue.count, 0);
  totalUnique = this.dataset.reduce((previousValue, currentValue) => previousValue + currentValue.unique, 0);
  totalActive = this.totalCount + this.totalUnique;

  createToolTipText(data) {
    let total = data.reduce((prev, acc) => prev + acc, 0);
    return `${total / this.totalActive}%`;
    // console.log(percent)
  }

  @action
  renderBarChart(element) {
    let dataset = this.dataset.sort((a, b) => a.count + a.unique - (b.count + b.unique)).reverse();
    let totalActive = this.totalActive;
    let stackFunction = stack().keys(['count', 'unique']);
    let stackedData = stackFunction(dataset); // returns an array of coordinates for each rectangle group, first group is for counts (left), second for unique (right)

    let container = select('.bar-chart-container');
    container.attr('viewBox', '0 0 751 405');
    container
      .append('div')
      .attr('class', 'chart-tooltip')
      .attr('style', 'position: absolute; opacity: 0;');

    let xScale = scaleLinear()
      .domain([0, max(dataset, d => d.count + d.unique)]) // min and max values of dataset
      .range([0, 75]); // range in percent (30% reserved for margins)

    let yScale = scaleBand()
      .domain(dataset.map(d => d.label))
      .range([0, dataset.length * 24]) // each bar element has a thickness (bar + padding) of 24 pixels
      // paddingInner takes a number between 0 and 1
      // it tells the scale the percent of the total width it should reserve for white space between bars
      .paddingInner(0.765);

    let chartSvg = select(element);
    chartSvg.attr('height', (dataset.length + 1) * 24);

    // add a group for each row of data
    let groups = chartSvg
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
      .style('cursor', 'pointer')
      .attr('width', data => `${xScale(data[1] - data[0] - 6)}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', data => `${xScale(data[0])}%`)
      .attr('y', ({ data }) => yScale(data.label))
      .attr('rx', 3)
      .attr('ry', 3)
      .on('mouseover', function({ data }) {
        select(this).attr('fill', '#1563FF');
        select('.chart-tooltip')
          .transition()
          .duration(200)
          .style('opacity', 1)
          .text(
            ` 
          ${Math.round((data.total * 100) / totalActive)}% of total client counts: \n
          ${data.unique} unique entities, ${data.count} active tokens.
          `
          )
          .style('color', 'white')
          .style('background', '#525761')
          .style('max-width', '200px')
          .style('font-size', '.929rem')
          .style('padding', '10px')
          .style('border-radius', '4px');
      })
      .on('mouseout', function(d) {
        select('.chart-tooltip').style('opacity', 0);
        select(this).attr('fill', function() {
          d[0] === 0 ? `${BAR_COLORS[0]}` : `${BAR_COLORS[1]}`;
        });
      })
      .on('mousemove', function() {
        select('.chart-tooltip')
          .style('left', `${xScale(event.pageX - 150)}%`)
          .style('top', `${event.pageY - 140}px`);
      });

    let totalNumbers = [];
    stackedData[1].forEach(e => {
      let n = e[1];
      totalNumbers.push(n);
    });

    let textData = [];
    rects.each(function(d, i) {
      let textDatum = {
        text: totalNumbers[i],
        x: parseFloat(select(this).attr('width')) + parseFloat(select(this).attr('x')),
        y: parseFloat(select(this).attr('y')) + parseFloat(select(this).attr('height')),
      };
      textData.push(textDatum);
    });

    let text = groups
      .selectAll('text')
      .data(textData)
      .enter()
      .append('text')
      .text(d => d.text)
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .style('font-size', '.8rem')
      .attr('text-anchor', 'start')
      .attr('y', d => {
        return `${d.y}`;
      })
      .attr('x', d => {
        let translateRight = d.x + 1;
        return `${translateRight}%`;
      });

    // remove axes lines
    groups.selectAll('.domain, .tick line').remove();

    let legendSvg = select('.legend');

    legendSvg
      .append('circle')
      .attr('cx', '60%')
      .attr('cy', '50%')
      .attr('r', 6)
      .style('fill', '#BFD4FF');
    legendSvg
      .append('text')
      .attr('x', '62%')
      .attr('y', '50%')
      .text(`${this.variableA}`)
      .style('font-size', '.8rem')
      .attr('alignment-baseline', 'middle');
    legendSvg
      .append('circle')
      .attr('cx', '83%')
      .attr('cy', '50%')
      .attr('r', 6)
      .style('fill', '#8AB1FF');
    legendSvg
      .append('text')
      .attr('x', '85%')
      .attr('y', '50%')
      .text(`${this.variableB}`)
      .style('font-size', '.8rem')
      .attr('alignment-baseline', 'middle');
  }
}

export default setComponentTemplate(layout, BarChart);
