import Component from '@ember/component';
import { computed } from '@ember/object';
/**
 * @module SelectableCard
 * SelectableCard components are card-like components that display a title, total, subtotal, and anything after they yield.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <SelectableCard @cardTitle="Tokens" @total={{totalHttpRequests}} @subText="Total" @gridContainer={{gridContainer}}/>
 * ```
 * @param {string} cardTitle - cardTitle displays the card title
 * @param {number} total - the Total number displays like a title, it's the largest text in the component
 * @param {string} subText - subText describes the total
 * @param {string} gridContainer - Optional parameter used to display CSS grid item class.
 */

export default Component.extend({
  tagName: '', // do not wrap component with div
  cardTitleComputed: computed('total', function() {
    let cardTitle = this.cardTitle || '';
    let total = this.total || '';

    if (cardTitle === 'Tokens') {
      return total !== 1 ? 'Tokens' : 'Token';
    } else if (cardTitle === 'Entities') {
      return total !== 1 ? 'Entities' : 'Entity';
    }

    return cardTitle;
  }),
  cardTitleTesting: computed('cardTitle', function() {
    let cardTitle = this.cardTitle || '';

    if (cardTitle === 'Http Requests') {
      return 'Requests';
    }

    return cardTitle;
  }),
});
