import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module HttpRequestsContainer
 * HttpRequestsContainer components are used to...
 *
 * @example
 * ```js
 * <HttpRequestsContainer @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */

const COUNTERS = [
  { start_time: '2018-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

export default Component.extend({
  counters: COUNTERS,
  timeWindow: 'All',
  filteredCounters: null,
  didReceiveAttrs() {
    this._super(...arguments);
    const { timeWindow, counters } = this;
    if (timeWindow === 'All') {
      this.set('filteredCounters', counters);
    }
    let filteredCounters = [];
    console.log(timeWindow);
    filteredCounters = counters.filter(counter => {
      const year = counter.start_time.substr(0, 4);
      return year === timeWindow;
    });
    debugger;
  },
  actions: {
    updateTimeWindow(newValue) {
      this.set('timeWindow', newValue);
    },
  },
});
