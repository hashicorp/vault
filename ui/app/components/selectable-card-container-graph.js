import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module MetricsSelectableCardContainerGraph
 * MetricsSelectableCardContainerGraph components are used to...
 *
 * @example
 * ```js
 * <MetricsSelectableCardContainerGraph @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default Component.extend({
  classNames: ['selectable-card-container-graph'],
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
