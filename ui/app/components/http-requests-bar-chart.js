import Component from '@ember/component';
import d3 from 'd3-selection';
import d3Scale from 'd3-scale';
import d3Axis from 'd3-axis';
import d3TimeFormat from 'd3-time-format';
import d3Tip from 'd3-tip';
import d3Transition from 'd3-transition';
import d3Ease from 'd3-ease';
import { assign } from '@ember/polyfills';
import { computed } from '@ember/object';
import { run } from '@ember/runloop';
import { task, waitForEvent } from 'ember-concurrency';

/**
 * @module HttpRequestsBarChart
 * HttpRequestsBarChart components are used to render a bar chart with the total number of HTTP Requests to a Vault server per month.
 *
 * @example
 * ```js
 * <HttpRequestsBarChart @counters={{counters}} />
 * ```
 *
 * @param counters=null {Array} - A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint which looks like:
 * COUNTERS = [
 *    {
 *       "start_time": "2019-05-01T00:00:00Z",
 *       "total": 50
 *     }
 * ]
 */

const HEIGHT = 240;
const HOVER_PADDING = 12;
const BASE_SPEED = 150;
const DURATION = BASE_SPEED * 2;

export default Component.extend({
  classNames: ['http-requests-bar-chart-container'],
  counters: null,
  margin: Object.freeze({ top: 24, right: 16, bottom: 24, left: 16 }),
  padding: 0.04,
  width: 0,
  height() {
    const { margin } = this;
    return HEIGHT - margin.top - margin.bottom;
  },

  parsedCounters: computed('counters', function() {
    // parse the start times so bars and ticks display properly
    const { counters } = this;
    return counters.map(counter => {
      return assign({}, counter, { start_time: d3TimeFormat.isoParse(counter.start_time) });
    });
  }),

  yScale: computed('parsedCounters', 'height', function() {
    const { parsedCounters } = this;
    const height = this.height();
    const counterTotals = parsedCounters.map(c => c.total);

    return d3Scale
      .scaleLinear()
      .domain([0, Math.max(...counterTotals)])
      .range([height, 0]);
  }),

  xScale: computed('parsedCounters', 'width', function() {
    const { parsedCounters, width, margin, padding } = this;

    return d3Scale
      .scaleBand()
      .domain(parsedCounters.map(c => c.start_time))
      .rangeRound([0, width - margin.left - margin.right], 0.05)
      .paddingInner(padding)
      .paddingOuter(padding);
  }),

  didInsertElement() {
    this._super(...arguments);
    const { margin } = this;

    // set the width after the element has been rendered because the chart axes depend on it.
    // this helps us avoid an arbitrary hardcoded width which causes alignment & resizing problems.
    run.schedule('afterRender', this, () => {
      this.set('width', this.element.clientWidth - margin.left - margin.right);
      this.renderBarChart();
    });
  },

  didUpdateAttrs() {
    this.renderBarChart();
  },

  renderBarChart() {
    const { margin, width, xScale, yScale, parsedCounters, elementId } = this;
    const height = this.height();
    const barChartSVG = d3.select('.http-requests-bar-chart');
    const barsContainer = d3.select(`#bars-container-${elementId}`);

    // initialize the tooltip
    const tip = d3Tip()
      .attr('class', 'd3-tooltip')
      .offset([HOVER_PADDING / 2, 0])
      .html(function(d) {
        const formatter = d3TimeFormat.utcFormat('%B %Y');
        return `
          <p class="date">${formatter(d.start_time)}</p>
          <p>${Intl.NumberFormat().format(d.total)} Requests</p>
        `;
      });

    barChartSVG.call(tip);

    // render the chart
    d3.select('.http-requests-bar-chart')
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .attr('viewBox', `0 0 ${width} ${height}`);

    // scale and render the axes
    const yAxis = d3Axis
      .axisRight(yScale)
      .ticks(3, '~s')
      .tickSizeOuter(0);
    const xAxis = d3Axis
      .axisBottom(xScale)
      .tickFormat(d3TimeFormat.utcFormat('%b %Y'))
      .tickSizeOuter(0);

    barChartSVG
      .select('g.x-axis')
      .attr('transform', `translate(0,${height})`)
      .call(xAxis);

    barChartSVG
      .select('g.y-axis')
      .attr('transform', `translate(${width - margin.left - margin.right}, 0)`)
      .call(yAxis);

    // render the gridlines
    const gridlines = d3Axis
      .axisRight(yScale)
      .ticks(3)
      .tickFormat('')
      .tickSize(width - margin.left - margin.right);

    barChartSVG.select('.gridlines').call(gridlines);

    // render the bars
    const bars = barsContainer.selectAll('.bar').data(parsedCounters, c => +c.start_time);

    const barsEnter = bars
      .enter()
      .append('rect')
      .attr('class', 'bar');

    const t = d3Transition
      .transition()
      .duration(DURATION)
      .ease(d3Ease.easeQuad);

    bars
      .merge(barsEnter)
      .attr('x', counter => xScale(counter.start_time))
      // set the initial y value to 0 so the bars animate upwards
      .attr('y', () => yScale(0))
      .attr('width', xScale.bandwidth())
      .transition(t)
      .delay(function(d, i) {
        return i * BASE_SPEED;
      })
      .attr('height', counter => height - yScale(counter.total))
      .attr('y', counter => yScale(counter.total));

    bars.exit().remove();

    // render transparent bars and bind the tooltip to them since we cannot
    // bind the tooltip to the actual bars. this is because the bars are
    // within a clipPath & you cannot bind DOM events to non-display elements.
    const shadowBarsContainer = d3.select('.shadow-bars');

    const shadowBars = shadowBarsContainer.selectAll('.bar').data(parsedCounters, c => +c.start_time);

    const shadowBarsEnter = shadowBars
      .enter()
      .append('rect')
      .attr('class', 'bar')
      .on('mouseenter', tip.show)
      .on('mouseleave', tip.hide);

    shadowBars
      .merge(shadowBarsEnter)
      .attr('width', xScale.bandwidth())
      .attr('height', counter => height - yScale(counter.total) + HOVER_PADDING)
      .attr('x', counter => xScale(counter.start_time))
      .attr('y', counter => yScale(counter.total) - HOVER_PADDING)
      .attr('fill', 'transparent')
      .attr('stroke', 'transparent');

    shadowBars.exit().remove();
  },

  updateDimensions() {
    const newWidth = this.element.clientWidth;
    const { margin } = this;

    this.set('width', newWidth - margin.left - margin.right);
    this.renderBarChart();
  },

  waitForResize: task(function*() {
    while (true) {
      yield waitForEvent(window, 'resize');
      run.scheduleOnce('afterRender', this, 'updateDimensions');
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement')
    .drop(),
});
