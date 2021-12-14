import Component from '@glimmer/component';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { stack } from 'd3-shape';
// eslint-disable-next-line no-unused-vars
import { select, event, selectAll } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { axisLeft } from 'd3-axis';
import { max, maxIndex } from 'd3-array';

/**
 * @module HorizontalBarChart
 * HorizontalBarChart components are used to display data in the form of a horizontal, stacked bar chart with accompanying legend and tooltip.
 *
 * @example
 * ```js
 * <HorizontalBarChart @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} dataset - dataset for the chart
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' for the chart legend
 * @param {string} [labelKey=label] - labelKey is the key name in the dataset passed in that corresponds to the value labeling the y-axis (i.e. 'namespace_path')
 * @param {string} [param1=defaultValue] - param1 is...
 */

// TODO: delete original bar chart component
// TODO: Move constants to helper

// SIZING CONSTANTS
const CHART_MARGIN = { top: 10, left: 137 }; // makes space for y-axis legend
const CHAR_LIMIT = 18; // character count limit for y-axis labels to trigger truncating
const LINE_HEIGHT = 24; // each bar w/ padding is 24 pixels thick

// COLOR THEME:
const BAR_COLOR_DEFAULT = ['#BFD4FF', '#1563FF'];
const BAR_COLOR_HOVER = ['#1563FF', '#0F4FD1'];
const BACKGROUND_BAR_COLOR = '#EBEEF2';

const SAMPLE_DATA = [
  {
    label: 'namespace80/',
    non_entity_tokens: 1696,
    distinct_entities: 1652,
    total: 3348,
  },
  {
    label: 'namespace12/',
    non_entity_tokens: 1568,
    distinct_entities: 1663,
    total: 3231,
  },
  {
    label: 'namespace44/',
    non_entity_tokens: 1511,
    distinct_entities: 1708,
    total: 3219,
  },
  {
    label: 'namespace36/',
    non_entity_tokens: 1574,
    distinct_entities: 1553,
    total: 3127,
  },
  {
    label: 'namespace2/',
    non_entity_tokens: 1784,
    distinct_entities: 1333,
    total: 3117,
  },
  {
    label: 'namespace82/',
    non_entity_tokens: 1245,
    distinct_entities: 1702,
    total: 2947,
  },
  {
    label: 'namespace28/',
    non_entity_tokens: 1579,
    distinct_entities: 1364,
    total: 2943,
  },
  {
    label: 'namespace60/',
    non_entity_tokens: 1962,
    distinct_entities: 929,
    total: 2891,
  },
  {
    label: 'namespace5/',
    non_entity_tokens: 1448,
    distinct_entities: 1418,
    total: 2866,
  },
  {
    label: 'namespace67/',
    non_entity_tokens: 1758,
    distinct_entities: 1065,
    total: 2823,
  },
];
export default class HorizontalBarChart extends Component {
  get labelKey() {
    return this.args.labelKey || 'label';
  }

  get chartLegend() {
    return this.args.chartLegend;
  }

  get topNamespace() {
    return this.args.dataset[maxIndex(this.args.dataset, d => d.total)];
  }

  @action
  renderChart(element, args) {
    // chart legend tells stackFunction how to stack/organize data
    // creates an array of data for each key name
    // each array contains coordinates for each data bar
    let stackFunction = stack().keys(this.chartLegend.map(l => l.key));
    let dataset = args[0];
    let stackedData = stackFunction(dataset);
    let labelKey = this.labelKey;

    let xScale = scaleLinear()
      .domain([0, max(dataset.map(d => d.total))])
      .range([0, 75]); // 25% reserved for margins

    let yScale = scaleBand()
      .domain(dataset.map(d => d[labelKey]))
      .range([0, dataset.length * LINE_HEIGHT])
      .paddingInner(0.765); // percent of the total width to reserve for padding between bars

    let chartSvg = select(element);
    chartSvg.attr('viewBox', `0 0 710 ${(dataset.length + 1) * LINE_HEIGHT}`);

    let groups = chartSvg
      .selectAll('g')
      .remove()
      .exit()
      .data(stackedData)
      .enter()
      .append('g')
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => BAR_COLOR_DEFAULT[i]);

    let yAxis = axisLeft(yScale).tickSize(0);
    yAxis(chartSvg.append('g').attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`));

    let truncate = selection =>
      selection.text(string =>
        string.length < CHAR_LIMIT ? string : string.slice(0, CHAR_LIMIT - 3) + '...'
      );

    chartSvg.selectAll('.tick text').call(truncate);
    chartSvg.select('.domain').remove();

    groups
      .selectAll('rect')
      // iterate through the stacked data and chart respectively
      .data(stackedData => stackedData)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', chartData => `${xScale(chartData[1] - chartData[0]) - 0.25}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', chartData => `${xScale(chartData[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3);
  }
}
