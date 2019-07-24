import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module HttpRequestsDropdown
 * HttpRequestsDropdown components are used to render a dropdown that filters the HttpRequestsBarChart.
 *
 * @example
 * ```js
 * <HttpRequestsDropdown @counters={counters} />
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
  classNames: ['http-requests-dropdown'],
  counters: null,
  timeWindow: 'All',
  options: computed('counters', function() {
    let counters = this.counters || [];
    let options = [];
    if (counters.length) {
      const years = counters
        .map(counter => {
          const year = new Date(counter.start_time);
          return year.getUTCFullYear().toString();
        })
        .uniq();
      years.sort().reverse();
      options = options.concat(years);
    }
    return options;
  }),
  onChange() {},
  actions: {
    onSelectTimeWindow(e) {
      const newValue = e.target.value;
      const { timeWindow } = this;
      if (newValue !== timeWindow) {
        this.onChange(newValue);
      }
    },
  },
});
