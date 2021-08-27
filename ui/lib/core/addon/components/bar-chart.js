/**
 * @module BarChart
 * BarChart components are used to...
 *
 * @example
 * ```js
 * <BarChart @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}} @onClick= @labelKey=/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} labelKey - labelKey is the key name the dataset uses to label each bar on the y-axis (i.e. "namespace_path" if the object passed in was: { namespace_path: "top-namespace", count: 4000 })
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 *
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
import { transition } from 'd3-transition';
import { guidFor } from '@ember/object/internals';

// SIZING CONSTANTS
const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 }; // makes space for y-axis legend
const CHAR_LIMIT = 18; // character count limit (for label truncating)

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#BFD4FF', '#8AB1FF'];
const BAR_COLOR_HOVER = ['#1563FF', '#0F4FD1'];
const BACKGROUND_BAR_COLOR = '#EBEEF2';
const TOOLTIP_BACKGROUND = '#525761';

class BarChart extends Component {
  // TODO: make xValue and yValue consts? i.e. yValue = dataset.map(d => d.label)
  // mapLegend = [{ key: 'count', label: 'Active direct tokens' }, { key: 'unique', label: 'Unique entities' }];
  mapLegend = [
    { key: 'non_entity_tokens', label: 'Active direct tokens' },
    { key: 'distinct_entities', label: 'Unique entities' },
  ];

  realData = [
    {
      namespace_id: 'root',
      namespace_path: 'root',
      counts: {
        distinct_entities: 268,
        non_entity_tokens: 985,
        clients: 1253,
      },
    },
    {
      namespace_id: 'O0i4m',
      namespace_path: 'top-namespace',
      counts: {
        distinct_entities: 648,
        non_entity_tokens: 220,
        clients: 868,
      },
    },
    {
      namespace_id: '1oihz',
      namespace_path: 'anotherNamespace',
      counts: {
        distinct_entities: 547,
        non_entity_tokens: 337,
        clients: 884,
      },
    },
    {
      namespace_id: '1oihz',
      namespace_path: 'someOtherNamespaceawgagawegawgawgawgaweg',
      counts: {
        distinct_entities: 807,
        non_entity_tokens: 234,
        clients: 1041,
      },
    },
  ];

  // the key name the dataset uses to label each bar on the y-axis
  get labelKey() {
    return this.args.labelKey || 'label';
  }

  get chartData() {
    return this.args.chartData || ['count', 'unique', 'total'];
  }

  // TODO: turn this into a get function that responds to arguments passed to component
  flattenedData() {
    return this.realData.map(d => {
      return {
        label: d['namespace_path'],
        non_entity_tokens: d['counts']['non_entity_tokens'],
        distinct_entities: d['counts']['distinct_entities'],
        total: d['counts']['clients'],
      };
    });
  }

  // TODO: separate into function for specifically creating tooltip text
  totalCount = this.realData.reduce((prevValue, currValue) => prevValue + currValue.counts.clients, 0);

  @action
  renderBarChart(element) {
    let totalCount = this.totalCount;
    let handleClick = this.args.onClick;
    let labelKey = this.labelKey;
    let dataset = this.flattenedData();
    let elementId = guidFor(element);
    let stackFunction = stack().keys(this.mapLegend.map(l => l.key));
    // creates an array of data for each map legend key
    // each array contains coordinates for each data bar
    let stackedData = stackFunction(this.flattenedData());

    // creates and appends tooltip
    let container = select('.bar-chart-container');
    container
      .append('div')
      .attr('class', 'chart-tooltip')
      .attr('style', 'position: fixed; opacity: 0;')
      .style('color', 'white')
      .style('background', `${TOOLTIP_BACKGROUND}`)
      .style('font-size', '.929rem')
      .style('padding', '10px')
      .style('border-radius', '4px');

    let xScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))])
      .range([0, 75]); // 25% reserved for margins

    let yScale = scaleBand()
      .domain(dataset.map(d => d[labelKey]))
      .range([0, dataset.length * 24]) // each bar w/ padding is 24 pixels wide
      .paddingInner(0.765); // percent of the total width to reserve for white space between bars

    let chartSvg = select(element);
    chartSvg.attr('viewBox', `0 0 710 ${(dataset.length + 1) * 24}`);
    chartSvg.attr('id', elementId);

    // creates group for each array of stackedData
    let groups = chartSvg
      .selectAll('g')
      .data(stackedData)
      .enter()
      .append('g')
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);

    let yAxis = axisLeft(yScale);
    yAxis(groups.append('g'));

    let truncate = selection =>
      selection.text(string =>
        string.length < CHAR_LIMIT ? string : string.slice(0, CHAR_LIMIT - 3) + '...'
      );

    chartSvg.selectAll('.tick text').call(truncate);

    let rects = groups
      .selectAll('rect')
      .data(d => d)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', chartData => `${xScale(chartData[1] - chartData[0] - 5)}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', chartData => `${xScale(chartData[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3);

    let actionBars = chartSvg
      .selectAll('.action-bar')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'action-bar')
      .attr('width', '100%')
      .attr('height', '24px')
      .attr('x', '0')
      .attr('y', chartData => yScale(chartData[labelKey]))
      .style('fill', `${BACKGROUND_BAR_COLOR}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let dataBars = chartSvg.selectAll('rect.data-bar');
    actionBars
      .on('click', function(chartData) {
        if (handleClick) {
          handleClick(chartData);
        }
      })
      .on('mouseover', function(chartData) {
        select(this).style('opacity', 1);
        dataBars
          .filter(function() {
            return select(this).attr('y') === `${event.target.getAttribute('y')}`;
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        // FUTURE TODO: Make tooltip text a function
        if (chartData.label.length >= CHAR_LIMIT || event.pageX > 522) {
          select('.chart-tooltip')
            .transition()
            .duration(200)
            .style('opacity', 1);
          if (chartData.label.length >= CHAR_LIMIT) {
            select('.chart-tooltip').style('max-width', 'fit-content');
          }
        }
      })
      .on('mouseout', function() {
        select(this).style('opacity', 0);
        select('.chart-tooltip').style('opacity', 0);
        dataBars
          .filter(function() {
            return select(this).attr('y') === `${event.target.getAttribute('y')}`;
          })
          .style('fill', (b, i) => `${BAR_COLOR_DEFAULT[i]}`);
      })
      .on('mousemove', function(chartData) {
        if (event.pageX < 522) {
          // don't hard code, but use y axis width to determine
          if (chartData.label.length >= CHAR_LIMIT) {
            select('.chart-tooltip')
              .style('left', `${event.pageX - 100}px`)
              .style('top', `${event.pageY - 50}px`)
              .text(`${chartData.label}`)
              .style('max-width', 'fit-content');
          } else {
            select('.chart-tooltip').style('opacity', 0);
          }
        } else {
          select('.chart-tooltip')
            .style('opacity', 1)
            .style('max-width', '200px')
            .style('left', `${event.pageX - 90}px`)
            .style('top', `${event.pageY - 90}px`).text(` 
                ${Math.round((chartData.total * 100) / totalCount)}% of total client counts:
                ${chartData.distinct_entities} unique entities, ${
            chartData.non_entity_tokens
          } active tokens.`);
        }
      });

    // creates total count text and coordinates to display to the right of data bars
    let totalCountData = [];
    rects.each(function(d) {
      let textDatum = {
        text: d.data.total,
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

    // TODO: y value needs to change when move onto another line
    // 20% of map key is reserved for each symbol + label, calculates starting x coord
    let startingXCoordinate = 100 - this.mapLegend.length * 20;
    let legendSvg = select('.legend');
    this.mapLegend.map((v, i) => {
      let xCoordinate = startingXCoordinate + i * 20;
      legendSvg
        .append('circle')
        .attr('cx', `${xCoordinate}%`)
        .attr('cy', '50%')
        .attr('r', 6)
        .style('fill', `${BAR_COLOR_DEFAULT[i]}`);
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
