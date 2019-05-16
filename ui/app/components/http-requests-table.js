import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module HttpRequestsTable
 * HttpRequestsTable components render a table with the total number of HTTP Requests to a Vault server per month.
 *
 * @example
 * ```js
 * const COUNTERS = [
 *    {
 *       "start_time": "2019-05-01T00:00:00Z",
 *       "total": 50
 *     }
 * ]
 *
 * <HttpRequestsTable @counters={{COUNTERS}} />
 * ```
 *
 * @param counters {Array} - A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint.
 */

export default Component.extend({
  classNames: ['http-requests-table'],
  counters: null,
  showChangeColumn: computed('counters', function() {
    const { counters } = this;
    return counters && counters.length > 1;
  }),
  countersWithChange: computed('counters', function() {
    const { counters } = this;

    if (counters) {
      let countersWithPercentChange = [];
      let previousMonthVal;

      counters.forEach(month => {
        if (previousMonthVal) {
          const change = (((month.total - previousMonthVal) / month.total) * 100).toFixed(1);
          const newCounter = Object.assign({ change }, month);
          countersWithPercentChange.push(newCounter);
        } else {
          // we're looking at the first counter in the list, so there is no % change value.
          countersWithPercentChange.push(month);
          previousMonthVal = month.total;
        }
      });
      return countersWithPercentChange;
    }
  }),
});
