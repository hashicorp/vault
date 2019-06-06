import Component from '@ember/component';
import { computed } from '@ember/object';
import { assign } from '@ember/polyfills';

/**
 * @module HttpRequestsTable
 * `HttpRequestsTable` components render a table with the total number of HTTP Requests to a Vault server per month.
 *
 * @example
 * ```js
 * <HttpRequestsTable @counters={{counters}} />
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
  tagName: '',
  counters: null,
  countersWithChange: computed('counters', function() {
    let counters = this.counters || [];
    let countersWithPercentChange = [];
    let previousMonthVal;

    counters.forEach(month => {
      if (previousMonthVal) {
        let percentChange = (((month.total - previousMonthVal) / month.total) * 100).toFixed(1);
        let glyph;
        if (percentChange > 0) {
          glyph = 'arrow-up';
        } else if (percentChange < 0) {
          glyph = 'arrow-down';
        }
        percentChange = Math.abs(percentChange);
        const newCounter = assign({ percentChange, glyph }, month);
        countersWithPercentChange.push(newCounter);
      } else {
        // we're looking at the first counter in the list, so there is no % change value.
        countersWithPercentChange.push(month);
        previousMonthVal = month.total;
      }
    });
    return countersWithPercentChange;
  }),
});
