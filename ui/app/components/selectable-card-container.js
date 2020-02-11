import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module SelectableCardContainer
 * SelectableCardContainer components are used to hold SelectableCard components.  They act as a CSS grid container, and change grid configurations based on the boolean of @gridContainer.
 *
 * @example
 * ```js
 * <SelectableCardContainer @counters={{model}} @gridContainer="true" />
 * ```
 * @param {object} counters - Counters is an object that returns total entities, tokens, and an array of objects with the total https request per month.
 * @param {string} [gridContainer] - gridContainer is optional.  If true, it's telling the container it will have a nested CSS grid.
 *
 * const MODEL = {
 *  totalEntities: 0,
 *  httpsRequests: [{ start_time: '2019-04-01T00:00:00Z', total: 5500 }],
 *  totalTokens: 1,
 * };
 */

export default Component.extend({
  classNameBindings: ['isGraphContainer'],
  isGraphContainer: computed('counters', function() {
    return this.counters.httpsRequests.length > 1
      ? 'selectable-card-container-graph'
      : 'selectable-card-container';
  }),
  totalHttpRequests: computed('counters', function() {
    let httpsRequestsArray = this.counters.httpsRequests || [];
    return httpsRequestsArray.firstObject.total;
  }),
  // Limit number of months returned to the most recent 12
  filteredHttpsRequests: computed('counters', function() {
    let httpsRequestsArray = this.counters.httpsRequests || [];
    if (httpsRequestsArray.length > 12) {
      httpsRequestsArray = httpsRequestsArray.slice(0, 12);
    }
    return httpsRequestsArray;
  }),
  percentChange: computed('counters', function() {
    let httpsRequestsArray = this.counters.httpsRequests || [];
    let lastTwoMonthsArray = httpsRequestsArray.slice(0, 2);
    let previousMonthVal = lastTwoMonthsArray.lastObject.total;
    let thisMonthVal = lastTwoMonthsArray.firstObject.total;

    let percentChange = (((previousMonthVal - thisMonthVal) / previousMonthVal) * 100).toFixed(1);
    // a negative value indicates a percentage increase, so we swap the value
    percentChange = -percentChange;
    return percentChange;
  }),
});
