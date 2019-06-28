import Component from '@ember/component';
import { computed } from '@ember/object';
import isWithinRange from 'date-fns/is_within_range';
import addMonths from 'date-fns/add_months';

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

export default Component.extend({
  classNames: ['http-requests-container'],
  counters: null,
  timeWindow: 'All',
  filteredCounters: computed('timeWindow', function() {
    const { counters, timeWindow } = this;
    if (timeWindow === 'All') {
      return counters;
    }

    let filteredCounters = [];
    if (timeWindow === 'Last 12 Months') {
      const today = new Date();
      const TwelveMonthsAgo = addMonths(today, -12);
      filteredCounters = counters.filter(counter => {
        return isWithinRange(counter.start_time, TwelveMonthsAgo, today);
      });

      return filteredCounters;
    }

    filteredCounters = counters.filter(counter => {
      const year = counter.start_time.substr(0, 4);
      return year === timeWindow;
    });
    return filteredCounters;
  }),
  actions: {
    updateTimeWindow(newValue) {
      this.set('timeWindow', newValue);
    },
  },
});
