/**
 * @module BarChart
 * BarChart components are used to...
 *
 * @example
 * ```js
 * <BarChart @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}} @onClick/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/bar-chart';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { select, event, selection, selectAll } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { max } from 'd3-array';
import { stack } from 'd3-shape';
import { axisLeft } from 'd3-axis';

const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 }; // makes space for y-axis legend
const BAR_COLORS_UNSELECTED = ['#BFD4FF', '#8AB1FF'];
const BAR_COLORS_SELECTED = ['#1563FF', '#0F4FD1'];
class BarChart extends Component {
  // TODO: make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
  variableA = 'Active direct tokens';
  variableB = 'Unique entities';
  mapKey = [{ key: 'count', label: 'Active direct tokens' }, { key: 'unique', label: 'Unique entities' }];
  // can pass in @labelKey but default will be 'label'
  get labelKey() {
    return this.args.labelKey || 'label';
  }

  realData = [
    {
      namespace_id: 'root',
      namespace_path: '',
      counts: {
        distinct_entities: 268,
        non_entity_tokens: 985,
        clients: 1253,
      },
    },
    {
      namespace_id: 'O0i4m',
      namespace_path: 'admin/',
      counts: {
        distinct_entities: 648,
        non_entity_tokens: 220,
        clients: 868,
      },
    },
    {
      namespace_id: '1oihz',
      namespace_path: 'ns1/',
      counts: {
        distinct_entities: 547,
        non_entity_tokens: 337,
        clients: 884,
      },
    },
    {
      namespace_id: '1oihz',
      namespace_path: 'ns1/',
      counts: {
        distinct_entities: 8078,
        non_entity_tokens: 5349,
        clients: 13427,
      },
    },
  ];

  dataset = [
    { label: 'top-namespace', count: 1212, unique: 300, total: 1512 },
    { label: 'namespace2', count: 650, unique: 550, total: 1200 },
    { label: 'longnamenamespace', count: 200, unique: 1000, total: 1200 },
    { label: 'namespacesomething', count: 400, unique: 400, total: 950 },
    { label: 'anothernamespace', count: 400, unique: 550, total: 1100 },
    { label: 'namespace5', count: 800, unique: 300, total: 800 },
    { label: 'namespace', count: 400, unique: 300, total: 700 },
    { label: 'namespace999', count: 350, unique: 250, total: 650 },
    { label: 'name-space', count: 450, unique: 200, total: 600 },
    { label: 'path/to/namespace', count: 200, unique: 100, total: 300 },
  ];

  totalCount = this.dataset.reduce((previousValue, currentValue) => previousValue + currentValue.count, 0);
  totalUnique = this.dataset.reduce((previousValue, currentValue) => previousValue + currentValue.unique, 0);
  totalActive = this.totalCount + this.totalUnique;

  @action
  renderBarChart(element) {
    let dataset = this.dataset.sort((a, b) => a.count + a.unique - (b.count + b.unique)).reverse();
    let totalActive = this.totalActive;
    let stackFunction = stack().keys(this.mapKey.map(l => l.key));
    let stackedData = stackFunction(dataset); // returns an array of coordinates for each group of rectangles, first group is for counts (left), second for unique (right)
    let container = select('.bar-chart-container');
    let handleClick = this.args.onClick;
    let keys = this.mapKey.map(e => e.key);
    let labelKey = this.labelKey;
    // creates and appends tooltip
    container
      .append('div')
      .attr('class', 'chart-tooltip')
      .attr('style', 'position: absolute; opacity: 0;')
      .style('color', 'white')
      .style('background', '#525761')
      .style('max-width', '200px')
      .style('font-size', '.929rem')
      .style('padding', '10px')
      .style('border-radius', '4px');

    let xScale = scaleLinear()
      .domain([0, max(dataset, d => keys.reduce((prevValue, currValue) => prevValue + d[currValue], 0))])
      .range([0, 75]); // 25% reserved for margins

    let yScale = scaleBand()
      .domain(dataset.map(d => d[labelKey]))
      // each bar element (bar + padding) has a thickness  of 24 pixels
      .range([0, dataset.length * 24])
      // paddingInner takes a number between 0 and 1
      // it tells the scale the percent of the total width it should reserve for white space between bars
      .paddingInner(0.765);

    let chartSvg = select(element);
    chartSvg.attr('viewBox', `0 0 710 ${(dataset.length + 1) * 24}`);

    // creates group for each array of stackedData
    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => BAR_COLORS_UNSELECTED[i]);

    let yAxis = axisLeft(yScale);
    yAxis(groups.append('g'));

    let rects = groups
      .selectAll('rect')
      .data(d => d)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', data => `${xScale(data[1] - data[0] - 6)}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', data => `${xScale(data[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('border', 1);

    let actionBars = chartSvg
      .selectAll('.foreground-bar')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'foreground-bar')
      .attr('width', '100%')
      .attr('height', '24px')
      .attr('x', '0')
      .attr('y', ({ label }) => yScale(label))
      .style('fill', '#EBEEF2')
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    actionBars
      .on('click', function(barData) {
        if (handleClick) {
          handleClick(barData);
        }
      })
      .on('mouseover', function(data) {
        select(this).style('opacity', 1);
        let dataBars = chartSvg.selectAll('rect.data-bar').filter(function() {
          return select(this).attr('y') === `${event.target.getAttribute('y')}`;
        });
        dataBars.style('fill', (b, i) => `${BAR_COLORS_SELECTED[i]}`);
        // FUTURE TODO: Make tooltip text a function
        select('.chart-tooltip')
          .transition()
          .duration(200)
          .style('opacity', 1).text(` 
      ${Math.round((data.total * 100) / totalActive)}% of total client counts: \n
      ${data.unique} unique entities, ${data.count} active tokens.
      `);
      })
      .on('mouseout', function() {
        select(this).style('opacity', 0);
        let dataBars = chartSvg.selectAll('rect.data-bar').filter(function() {
          return select(this).attr('y') === `${event.target.getAttribute('y')}`;
        });
        dataBars.style('fill', (b, i) => `${BAR_COLORS_UNSELECTED[i]}`);
        select('.chart-tooltip').style('opacity', 0);
      })
      .on('mousemove', function() {
        select('.chart-tooltip')
          .style('left', `${event.pageX - 30}px`)
          .style('top', `${event.pageY - 143}px`);
      });

    // TODO: fix this inflexible business
    let totalNumbers = [];
    stackedData[1].forEach(e => {
      let n = e[1];
      totalNumbers.push(n);
    });

    let totalCountData = [];
    rects.each(function(d, i) {
      let textDatum = {
        text: totalNumbers[i],
        x: parseFloat(select(this).attr('width')) + parseFloat(select(this).attr('x')),
        y: parseFloat(select(this).attr('y')) + parseFloat(select(this).attr('height')),
      };
      totalCountData.push(textDatum);
    });

    groups
      .selectAll('text')
      .data(totalCountData)
      .enter()
      .append('text')
      .text(d => d.text)
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .style('font-size', '.8rem')
      .attr('text-anchor', 'start')
      .attr('y', d => `${d.y}`)
      .attr('x', d => `${d.x + 1}%`);

    // removes axes lines
    groups.selectAll('.domain, .tick line').remove();

    // TODO: y needs to change when move onto another line
    // 20% of map key is reserved for each symbol + label, calculates starting x coord
    let startingXCoordinate = 100 - this.mapKey.length * 20;
    let legendSvg = select('.legend');
    this.mapKey.map((v, i) => {
      let xCoordinate = startingXCoordinate + i * 20;
      legendSvg
        .append('circle')
        .attr('cx', `${xCoordinate}%`)
        .attr('cy', '50%')
        .attr('r', 6)
        .style('fill', `${BAR_COLORS_UNSELECTED[i]}`);
      legendSvg
        .append('text')
        .attr('x', `${xCoordinate + 2}%`)
        .attr('y', '50%')
        .text(`${v.label}`)
        .style('font-size', '.8rem')
        .attr('alignment-baseline', 'middle');
    });
  }
}

export default setComponentTemplate(layout, BarChart);
