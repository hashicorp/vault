// ARG TODO may use may not
import Component from '@ember/component';
import { computed } from '@ember/object';
/**
 * @module MetricsSelectableCard
 * MetricsSelectableCard components are used to...
 *
 * @example
 * ```js
 * <MetricsSelectableCard @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default Component.extend({
  cardType: computed('type', function() {
    let type = this.type || '';
    let total = this.total || '';

    if (type === 'Tokens') {
      return total !== 1 ? 'Tokens' : 'Token';
    } else if (type === 'Entities') {
      return total !== 1 ? 'Entities' : 'Entity';
    }
    return 'Http Requests';
  }),
  subText: computed('type', function() {
    let type = this.type || '';

    if (type === 'httpRequests') {
      return 'This month';
    }
    return 'Total';
  }),
});
