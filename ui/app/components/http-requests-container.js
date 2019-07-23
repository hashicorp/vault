import Component from '@ember/component';
import { computed } from '@ember/object';
import isWithinRange from 'date-fns/is_within_range';
import addMonths from 'date-fns/add_months';

/**
 * @module HttpRequestsContainer
 * The HttpRequestsContainer component is the parent component of the HttpRequestsDropdown, HttpRequestsBarChart, and HttpRequestsTable components. It is used to handle filtering the bar chart and table according to selected time window from the dropdown.
 *
 * @example
 * ```js
 * <HttpRequestsContainer @counters={counters}/>
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

export default Component.extend({
  classNames: ['http-requests-container'],
  counters: null,
  timeWindow: 'All',
  filteredCounters: computed('counters', 'timeWindow', function() {
    const { counters, timeWindow } = this;
    if (timeWindow === 'All') {
      return counters;
    }

    let filteredCounters = [];
    if (timeWindow === 'Last 12 Months') {
      const today = new Date();
      const twelveMonthsAgo = addMonths(today, -12);
      filteredCounters = counters.filter(counter => {
        return isWithinRange(counter.start_time, twelveMonthsAgo, today);
      });

      return filteredCounters;
    }

    filteredCounters = counters.filter(counter => {
      const year = new Date(counter.start_time).getUTCFullYear();
      return year.toString() === timeWindow;
    });
    return filteredCounters;
  }),
  actions: {
    updateTimeWindow(newValue) {
      this.set('timeWindow', newValue);
    },
  },
});
