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
 * @param cardTitle='' {String} - cardTitle displays the card title
 * @param total=0 {Number} - the Total number displays like a title, it's the largest text in the component
 * @param subText='' {String} - subText describes the total
 * @param gridContainer=false {Boolean} - Optional parameter used to display CSS grid item class.
 */

export default Component.extend({
  cardTitle: '',
  total: 0,
  subText: '',
  gridContainer: false,
  tagName: '', // do not wrap component with div
  formattedCardTitle: computed('total', function() {
    const { cardTitle, total } = this;

    if (cardTitle === 'Tokens') {
      return total !== 1 ? 'Tokens' : 'Token';
    } else if (cardTitle === 'Entities') {
      return total !== 1 ? 'Entities' : 'Entity';
    }

    return cardTitle;
  }),
});
