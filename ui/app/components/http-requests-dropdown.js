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
  options: computed('counters', function() {
    let counters = this.counters || [];
    // const options = {'All', 'Last 12 Months'};
    if (counters.length > 2) {
      const years = counters.map(counter => {
        debugger;
        return counter.start_time.slice(0, 3);
      });
    }
  }),
});
