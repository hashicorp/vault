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

const CHART_MARGIN = { top: 10, right: 24, bottom: 26, left: 137 }; // makes space for y-axis legend

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
    let stackFunction = stack().keys(this.mapLegend.map(l => l.key));
    let stackedData = stackFunction(this.flattenedData()); // returns an array of coordinates for each group of rectangles (first left, then right)
    let container = select('.bar-chart-container');
    let handleClick = this.args.onClick;
    let labelKey = this.labelKey;
    let dataset = this.flattenedData();
    let elementId = guidFor(element);
    // creates and appends tooltip
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
      // each bar element (bar + padding) has a thickness  of 24 pixels
      .range([0, dataset.length * 24])
      // paddingInner takes a number between 0 and 1
      // it tells the scale the percent of the total width it should reserve for white space between bars
      .paddingInner(0.765);

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

    let truncate = function(selection) {
      selection.text(function(string) {
        return string.length < 18 ? string : string.slice(0, 18 - 3) + '...';
      });
    };

    chartSvg.selectAll('.tick text').call(truncate);

    let rects = groups
      .selectAll('rect')
      .data(d => d)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', data => `${xScale(data[1] - data[0] - 6)}%`)
      .attr('height', 6)
      // .attr('height', yScale.bandwidth()) <- don't want to scale because want bar width set at 6 pixels
      .attr('x', data => `${xScale(data[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3)
      .attr('border', 1);

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
      .attr('y', ({ label }) => yScale(label))
      .style('fill', `${BACKGROUND_BAR_COLOR}`)
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
        dataBars.style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        // FUTURE TODO: Make tooltip text a function
        if (data.label.length >= 18 || event.pageX > 522) {
          select('.chart-tooltip')
            .transition()
            .duration(200)
            .style('opacity', 1);
          if (data.label.length >= 18) {
            select('.chart-tooltip').style('max-width', 'fit-content');
          }
        }
      })
      .on('mouseout', function() {
        select(this).style('opacity', 0);
        let dataBars = chartSvg.selectAll('rect.data-bar').filter(function() {
          return select(this).attr('y') === `${event.target.getAttribute('y')}`;
        });
        dataBars.style('fill', (b, i) => `${BAR_COLOR_DEFAULT[i]}`);
        select('.chart-tooltip').style('opacity', 0);
      })
      .on('mousemove', function(data) {
        if (event.pageX < 522) {
          // don't hard code, but use y axis width to determine
          if (data.label.length >= 18) {
            select('.chart-tooltip')
              .style('left', `${event.pageX - 100}px`)
              .style('top', `${event.pageY - 50}px`)
              .text(`${data.label}`)
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
                ${Math.round((data.total * 100) / totalCount)}% of total client counts:
                ${data.distinct_entities} unique entities, ${data.non_entity_tokens} active tokens.`);
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
