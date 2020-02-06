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
  cardTitleComputed: computed('type', function() {
    let cardTitle = this.cardTitle || '';
    let total = this.total || '';

    if (cardTitle === 'Tokens') {
      return total !== 1 ? 'Tokens' : 'Token';
    } else if (cardTitle === 'Entities') {
      return total !== 1 ? 'Entities' : 'Entity';
    }

    return cardTitle;
  }),
});
