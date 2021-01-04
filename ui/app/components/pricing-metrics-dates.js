/**
 * @module PricingMetricsDates
 * PricingMetricsDates components are used on the Pricing Metrics page to handle queries related to pricing metrics.
 * This component assumes that query parameters (as in, from route params) are being passed in with the format MM-YYYY,
 * while the inputs expect a format of MM/YYYY.
 *
 * @example
 * ```js
 * <PricingMetricsDates @resultStart="2020-03-01T00:00:00Z" @resultEnd="2020-08-31T23:59:59Z" @queryStart="03-2020" @queryEnd="08-2020" />
 * ```
 * @param {object} resultStart - resultStart is the start date of the metrics returned. Should be a valid date string that the built-in Date() fn can parse
 * @param {object} resultEnd - resultEnd is the end date of the metrics returned. Should be a valid date string that the built-in Date() fn can parse
 * @param {string} [queryStart] - queryStart is the route param (formatted MM-YYYY) that the result will be measured against for showing discrepancy warning
 * @param {string} [queryEnd] - queryEnd is the route param (formatted MM-YYYY) that the result will be measured against for showing discrepancy warning
 * @param {number} [defaultSpan=12] - setting for default time between start and end input dates
 * @param {number} [retentionMonths=24] - setting for the retention months, which informs valid dates to query by
 */
import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import {
  differenceInSeconds,
  isValid,
  subMonths,
  startOfToday,
  format,
  endOfMonth,
  startOfMonth,
  isBefore,
} from 'date-fns';
import layout from '../templates/components/pricing-metrics-dates';
import { parseDateString } from 'vault/helpers/parse-date-string';

export default Component.extend({
  layout,
  router: service(),

  queryStart: null,
  queryEnd: null,
  resultStart: null,
  resultEnd: null,

  start: null,
  end: null,

  defaultSpan: 12,
  retentionMonths: 24,

  startDate: computed('start', function() {
    if (!this.start) return null;
    let date;
    try {
      date = parseDateString(this.start, '/');
      if (date) return date;
      return null;
    } catch (e) {
      return null;
    }
  }),
  endDate: computed('end', function() {
    if (!this.end) return null;
    let date;
    try {
      date = parseDateString(this.end, '/');
      if (date) return endOfMonth(date);
      return null;
    } catch (e) {
      return null;
    }
  }),

  showResultsWarning: computed('resultStart', 'resultEnd', function() {
    if (!this.queryStart || !this.queryEnd || !this.resultStart || !this.resultEnd) {
      return false;
    }
    const resultStart = new Date(this.resultStart);
    const resultEnd = new Date(this.resultEnd);
    let queryStart, queryEnd;
    try {
      queryStart = parseDateString(this.queryStart, '-');
      queryEnd = parseDateString(this.queryEnd, '-');
    } catch (e) {
      // Log error for debugging purposes
      console.debug(e);
    }

    if (!queryStart || !queryEnd || !isValid(resultStart) || !isValid(resultEnd)) {
      return false;
    }
    if (Math.abs(differenceInSeconds(queryStart, resultStart)) >= 86400) {
      return true;
    }
    if (Math.abs(differenceInSeconds(resultEnd, endOfMonth(queryEnd))) >= 86400) {
      return true;
    }
    return false;
  }),

  error: computed('end', 'start', function() {
    if (!this.startDate) {
      return 'Start date is invalid. Please use format MM/YYYY';
    }
    if (!this.endDate) {
      return 'End date is invalid. Please use format MM/YYYY';
    }
    if (isBefore(this.endDate, this.startDate)) {
      return 'Start date is after end date';
    }
    const lastMonthAvailable = endOfMonth(subMonths(startOfToday(), 1));
    if (isBefore(lastMonthAvailable, this.endDate)) {
      return `Data is not available until the end of the month`;
    }
    const earliestRetained = startOfMonth(subMonths(lastMonthAvailable, this.retentionMonths));
    if (isBefore(this.startDate, earliestRetained)) {
      return `No data retained before ${format(earliestRetained, 'MM/YYYY')} due to your settings`;
    }

    return null;
  }),

  init() {
    this._super(...arguments);
    let initialEnd;
    let initialStart;

    initialEnd = subMonths(startOfToday(), 1);
    if (this.queryEnd) {
      initialEnd = parseDateString(this.queryEnd, '-');
    } else {
      // if query isn't passed in, set it so that showResultsWarning works
      this.queryEnd = format(initialEnd, 'MM-YYYY');
    }
    initialStart = subMonths(initialEnd, this.defaultSpan);
    if (this.queryStart) {
      initialStart = parseDateString(this.queryStart, '-');
    } else {
      // if query isn't passed in, set it so that showResultsWarning works
      this.queryStart = format(initialStart, 'MM-YYYY');
    }

    this.start = format(initialStart, 'MM/YYYY');
    this.end = format(initialEnd, 'MM/YYYY');
  },

  actions: {
    handleQuery() {
      const start = format(this.startDate, 'MM-YYYY');
      const end = format(this.endDate, 'MM-YYYY');
      this.router.transitionTo('vault.cluster.metrics', {
        queryParams: {
          start,
          end,
        },
      });
    },
  },
});
