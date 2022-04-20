import Component from '@glimmer/component';
import { action } from '@ember/object';
import { stack } from 'd3-shape';
// eslint-disable-next-line no-unused-vars
import { select, event, selectAll } from 'd3-selection';
import { scaleLinear, scaleBand } from 'd3-scale';
import { axisLeft } from 'd3-axis';
import { max, maxIndex } from 'd3-array';
import { BAR_COLOR_HOVER, GREY, LIGHT_AND_DARK_BLUE, formatTooltipNumber } from 'vault/utils/chart-helpers';
import { tracked } from '@glimmer/tracking';

/**
 * @module HorizontalBarChart
 * HorizontalBarChart components are used to display data in the form of a horizontal, stacked bar chart with accompanying tooltip.
 *
 * @example
 * ```js
 * <HorizontalBarChart @dataset={{@dataset}} @chartLegend={{@chartLegend}}/>
 * ```
 * @param {array} dataset - dataset for the chart, must be an array of flattened objects
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 */

// SIZING CONSTANTS
const CHART_MARGIN = { top: 10, left: 95 }; // makes space for y-axis legend
const TRANSLATE = { down: 14, left: 99 };
const CHAR_LIMIT = 15; // character count limit for y-axis labels to trigger truncating
const LINE_HEIGHT = 24; // each bar w/ padding is 24 pixels thick

export default class HorizontalBarChart extends Component {
  @tracked tooltipTarget = '';
  @tracked tooltipText = '';
  @tracked isLabel = null;

  get labelKey() {
    return this.args.labelKey || 'label';
  }

  get chartLegend() {
    return this.args.chartLegend;
  }

  get topNamespace() {
    return this.args.dataset[maxIndex(this.args.dataset, (d) => d.clients)];
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  renderChart(element, args) {
    // chart legend tells stackFunction how to stack/organize data
    // creates an array of data for each key name
    // each array contains coordinates for each data bar
    let stackFunction = stack().keys(this.chartLegend.map((l) => l.key));
    let dataset = args[0];
    let stackedData = stackFunction(dataset);
    let labelKey = this.labelKey;

    let xScale = scaleLinear()
      .domain([0, max(dataset.map((d) => d.clients))])
      .range([0, 75]); // 25% reserved for margins

    let yScale = scaleBand()
      .domain(dataset.map((d) => d[labelKey]))
      .range([0, dataset.length * LINE_HEIGHT])
      .paddingInner(0.765); // percent of the total width to reserve for padding between bars

    let chartSvg = select(element);
    chartSvg.attr('width', '100%').attr('viewBox', `0 0 564 ${(dataset.length + 1) * LINE_HEIGHT}`);

    let dataBarGroup = chartSvg
      .selectAll('g')
      .remove()
      .exit()
      .data(stackedData)
      .enter()
      .append('g')
      .attr('data-test-group', (d) => `${d.key}`)
      // shifts chart to accommodate y-axis legend
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`)
      .style('fill', (d, i) => LIGHT_AND_DARK_BLUE[i]);

    let yAxis = axisLeft(yScale).tickSize(0);

    let yLabelsGroup = chartSvg
      .append('g')
      .attr('data-test-group', 'y-labels')
      .attr('transform', `translate(${CHART_MARGIN.left}, ${CHART_MARGIN.top})`);
    yAxis(yLabelsGroup);

    chartSvg.select('.domain').remove();

    let truncate = (selection) =>
      selection.text((string) =>
        string.length < CHAR_LIMIT ? string : string.slice(0, CHAR_LIMIT - 3) + '...'
      );

    chartSvg.selectAll('.tick text').call(truncate);

    dataBarGroup
      .selectAll('rect')
      .remove()
      .exit()
      // iterate through the stacked data and chart respectively
      .data((stackedData) => stackedData)
      .enter()
      .append('rect')
      .attr('class', 'data-bar')
      .style('cursor', 'pointer')
      .attr('width', (chartData) => `${xScale(Math.abs(chartData[1] - chartData[0]))}%`)
      .attr('height', yScale.bandwidth())
      .attr('x', (chartData) => `${xScale(chartData[0])}%`)
      .attr('y', ({ data }) => yScale(data[labelKey]))
      .attr('rx', 3)
      .attr('ry', 3);

    let actionBarGroup = chartSvg.append('g').attr('data-test-group', 'action-bars');

    let actionBars = actionBarGroup
      .selectAll('.action-bar')
      .remove()
      .exit()
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'action-bar')
      .attr('width', '100%')
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', (chartData) => yScale(chartData[labelKey]))
      .style('fill', `${GREY}`)
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let labelActionBarGroup = chartSvg.append('g').attr('data-test-group', 'label-action-bars');

    let labelActionBar = labelActionBarGroup
      .selectAll('.label-action-bar')
      .remove()
      .exit()
      .data(dataset)
      .enter()
      .append('rect')
      .style('cursor', 'pointer')
      .attr('class', 'label-action-bar')
      .attr('width', CHART_MARGIN.left)
      .attr('height', `${LINE_HEIGHT}px`)
      .attr('x', '0')
      .attr('y', (chartData) => yScale(chartData[labelKey]))
      .style('opacity', '0')
      .style('mix-blend-mode', 'multiply');

    let dataBars = chartSvg.selectAll('rect.data-bar');
    let actionBarSelection = chartSvg.selectAll('rect.action-bar');

    let compareAttributes = (elementA, elementB, attr) =>
      select(elementA).attr(`${attr}`) === select(elementB).attr(`${attr}`);

    // MOUSE EVENTS FOR DATA BARS
    actionBars
      .on('mouseover', (data) => {
        let hoveredElement = actionBars.filter((bar) => bar.label === data.label).node();
        this.tooltipTarget = hoveredElement;
        this.isLabel = false;
        this.tooltipText = `${Math.round((data.clients * 100) / this.args.totalUsageCounts.clients)}% 
        of total client counts:
        ${formatTooltipNumber(data.entity_clients)} entity clients, 
        ${formatTooltipNumber(data.non_entity_clients)} non-entity clients.`;

        select(hoveredElement).style('opacity', 1);

        dataBars
          .filter(function () {
            return compareAttributes(this, hoveredElement, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
      })
      .on('mouseout', function () {
        select(this).style('opacity', 0);
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${LIGHT_AND_DARK_BLUE[i]}`);
      });

    // MOUSE EVENTS FOR Y-AXIS LABELS
    labelActionBar
      .on('mouseover', (data) => {
        if (data.label.length >= CHAR_LIMIT) {
          let hoveredElement = labelActionBar.filter((bar) => bar.label === data.label).node();
          this.tooltipTarget = hoveredElement;
          this.isLabel = true;
          this.tooltipText = data.label;
        } else {
          this.tooltipTarget = null;
        }
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${BAR_COLOR_HOVER[i]}`);
        actionBarSelection
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '1');
      })
      .on('mouseout', function () {
        this.tooltipTarget = null;
        dataBars
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('fill', (b, i) => `${LIGHT_AND_DARK_BLUE[i]}`);
        actionBarSelection
          .filter(function () {
            return compareAttributes(this, event.target, 'y');
          })
          .style('opacity', '0');
      });

    // client count total values to the right
    let totalValueGroup = chartSvg
      .append('g')
      .attr('data-test-group', 'total-values')
      .attr('transform', `translate(${TRANSLATE.left}, ${TRANSLATE.down})`);

    totalValueGroup
      .selectAll('text')
      .data(dataset)
      .enter()
      .append('text')
      .text((d) => d.clients)
      .attr('fill', '#000')
      .attr('class', 'total-value')
      .style('font-size', '.8rem')
      .attr('text-anchor', 'start')
      .attr('alignment-baseline', 'middle')
      .attr('x', (chartData) => `${xScale(chartData.clients)}%`)
      .attr('y', (chartData) => yScale(chartData.label));
  }
}
