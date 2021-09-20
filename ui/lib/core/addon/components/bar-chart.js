/**
 * @module BarChart
 * BarChart components are used to display data in the form of a stacked bar chart, with accompanying legend and tooltip. Anything passed into the block will display in the top right of the chart header.
 *
 * @example
 * ```js
 * <BarChartComponent @title="Top 10 Namespaces" @description="Each namespace's client count includes clients in child namespaces." @labelKey="namespace_path" @dataset={{this.testData}} @mapLegend={{ array (hash key="non_entity_tokens" label="Active direct tokens") (hash key="distinct_entities" label="Unique Entities") }} @onClick={{this.onClick}} >
 *    <button type="button" class="link">
 *      Export all namespace data
 *    </button>/>
 * </BarChartComponent>
 *
 *  mapLegendSample = [{
 *     key: "api_key_for_label",
 *     label: "Label Displayed on Legend"
 *   }]
 * ```
 *
 * @param {string} title - title of the chart
 * @param {array} mapLegend - array of objects with key names 'key' and 'label' for the map legend
 * @param {object} dataset - dataset for the chart
 * @param {string} [description] - description of the chart
 * @param {string} [labelKey=label] - labelKey is the key name in the dataset passed in that corresponds to the value labeling the y-axis
 * @param {function} [onClick] - takes function from parent and passes it to click event on data bars
 *
 */

import Component from '@glimmer/component';
import layout from '../templates/components/bar-chart';
import { setComponentTemplate } from '@ember/component';
import { assert } from '@ember/debug';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import { scaleLinear, scaleBand } from 'd3-scale';
import { axisLeft } from 'd3-axis';
import { max } from 'd3-array';
import { stack } from 'd3-shape';
// eslint-disable-next-line no-unused-vars
import { select, event, selectAll } from 'd3-selection';
// eslint-disable-next-line no-unused-vars
import { transition } from 'd3-transition';

// SIZING CONSTANTS
const CHART_MARGIN = { top: 10, left: 137 }; // makes space for y-axis legend
const CHAR_LIMIT = 18; // character count limit for y-axis labels to trigger truncating
const LINE_HEIGHT = 24; // each bar w/ padding is 24 pixels thick

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#BFD4FF', '#8AB1FF'];
const BAR_COLOR_HOVER = ['#1563FF', '#0F4FD1'];
const BACKGROUND_BAR_COLOR = '#EBEEF2';
const TOOLTIP_BACKGROUND = '#525761';

class BarChartComponent extends Component {
  get labelKey() {
    return this.args.labelKey || 'label';
  }

  get mapLegend() {
    assert(
      'map legend is required, must be an array of objects with key names of "key" and "label"',
      this.hasLegend()
    );
    return this.args.mapLegend;
  }

  get dataset() {
    return this.args.dataset || null;
  }

  hasLegend() {
    if (!this.args.mapLegend || !Array.isArray(this.args.mapLegend)) {
      return false;
    } else {
      let legendKeys = this.args.mapLegend.map(obj => Object.keys(obj));
      return legendKeys.map(array => array.includes('key', 'label')).every(element => element === true);
    }
  }

  @action
  renderBarChart(element) {
    let elementId = guidFor(element);
    let totalCount = this.dataset.reduce((prevValue, currValue) => prevValue + currValue.total, 0);
    let handleClick = this.args.onClick;
    let labelKey = this.labelKey;
    let dataset = this.dataset;
    let stackFunction = stack().keys(this.mapLegend.map(l => l.key));
    // creates an array of data for each map legend key
    // each array contains coordinates for each data bar
    let stackedData = stackFunction(dataset);

    // creates and appends tooltip
    let container = select('.bar-chart-container');
    container
      .append('div')
      .attr('class', 'chart-tooltip')
      .attr('style', 'position: absolute; opacity: 0;')
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
      .range([0, dataset.length * LINE_HEIGHT])
      .paddingInner(0.765); // percent of the total width to reserve for padding between bars

    let chartSvg = select(element);
    chartSvg.attr('viewBox', `0 0 710 ${(dataset.length + 1) * LINE_HEIGHT}`);
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
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', chartData => yScale(chartData[labelKey]))
      .style('fill', `${BACKGROUND_BAR_COLOR}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let yLegendBars = chartSvg
      .selectAll('.label-bar')
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'label-action-bar')
      .attr('width', CHART_MARGIN.left)
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', chartData => yScale(chartData[labelKey]))
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let dataBars = chartSvg.selectAll('rect.data-bar');
    let actionBarSelection = chartSvg.selectAll('rect.action-bar');
    let compareAttributes = (elementA, elementB, attr) =>
      select(elementA).attr(`${attr}`) === elementB.getAttribute(`${attr}`);

    // handles click and mouseover/out/move event for data bars
    actionBars
      .on('click', function(chartData) {
        if (handleClick) {
          handleClick(chartData);
        }
      })
      .on('mouseover', function() {
        select(this).style('opacity', 1);
        dataBars
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        select('.chart-tooltip')
          .transition()
          .duration(200)
          .style('opacity', 1);
      })
      .on('mouseout', function() {
        select(this).style('opacity', 0);
        select('.chart-tooltip').style('opacity', 0);
        dataBars
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_DEFAULT[i]}`);
      })
      .on('mousemove', function(chartData) {
        select('.chart-tooltip')
          .style('opacity', 1)
          .style('max-width', '200px')
          .style('left', `${event.pageX - 95}px`)
          .style('top', `${event.pageY - 155}px`)
          .text(
            `${Math.round((chartData.total * 100) / totalCount)}% of total client counts:
            ${chartData.non_entity_tokens} non-entity tokens, ${chartData.distinct_entities} unique entities.
          `
          );
      });

    // handles mouseover/out/move event for y-axis legend
    yLegendBars
      .on('click', function(chartData) {
        if (handleClick) {
          handleClick(chartData);
        }
      })
      .on('mouseover', function(chartData) {
        dataBars
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        actionBarSelection
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '1');
        if (chartData.label.length >= CHAR_LIMIT) {
          select('.chart-tooltip')
            .transition()
            .duration(200)
            .style('opacity', 1);
        }
      })
      .on('mouseout', function() {
        select('.chart-tooltip').style('opacity', 0);
        dataBars
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_DEFAULT[i]}`);
        actionBarSelection
          .filter(function() {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '0');
      })
      .on('mousemove', function(chartData) {
        if (chartData.label.length >= CHAR_LIMIT) {
          select('.chart-tooltip')
            .style('left', `${event.pageX - 100}px`)
            .style('top', `${event.pageY - 50}px`)
            .text(`${chartData.label}`)
            .style('max-width', 'fit-content');
        } else {
          select('.chart-tooltip').style('opacity', 0);
        }
      });

    // TODO: these render twice, need to only render and append once per line
    // creates total count text and coordinates to display to the right of data bars
    let totalCountData = [];
    rects.each(function(d) {
      let textDatum = {
        total: d.data.total,
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
      .text(d => d.total)
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .style('font-size', '.8rem')
      .attr('text-anchor', 'start')
      .attr('y', d => `${d.y}`)
      .attr('x', d => `${d.x + 1}%`);

    // removes axes lines
    groups.selectAll('.domain, .tick line').remove();

    // TODO: make more flexible, make legend a div instead of svg?
    // each map key symbol & label takes up 20% of legend SVG width
    let startingXCoordinate = 100 - this.mapLegend.length * 20; // subtract from 100% to find starting x-coordinate
    let legendSvg = select('.legend');
    this.mapLegend.map((legend, i) => {
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
        .text(`${legend.label}`)
        .style('font-size', '.8rem')
        .attr('alignment-baseline', 'middle');
    });
  }
}

export default setComponentTemplate(layout, BarChartComponent);
