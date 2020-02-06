import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module MetricsSelectableCardContainer
 * MetricsSelectableCardContainer components are used to...
 *
 * @example
 * ```js
 * <MetricsSelectableCardContainer @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default Component.extend({
  totalHttpRequests: computed('counters', function() {
    let httpsRequestsArray = this.counters.httpsRequests || [];
    return httpsRequestsArray.firstObject.total;
  }),
});
