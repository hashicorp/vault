/**
 * @module SelectableCard
 * SelectableCard components are card-like components that display a title, total, subtotal, and anything after the yield.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <SelectableCard @cardTitle="Tokens" @total={{totalHttpRequests}} @subText="Total"/>
 * ```
 * @param {string} [cardTitle] - cardTitle displays the card title.
 * @param {number} [total = 0] - the Total number displays like a title, it's the largest text in the component.
 * @param {string} [subText] - subText describes the total.
 * @param {boolean} [actionCard = false] - false default selectable card container used in metrics, true a card that focus on actions as seen in database secret engine overview.
 * @param {string} [actionText] - that action that happens in an actionCard.
 * @param {string} [actionTo] - route where link will take you.
 * @param {string} [queryParam] - tab for the route the link will take you.
 * @param {string} [type] - type used in the link type.
 */

import Component from '@glimmer/component';
export default class SelectableCard extends Component {
  get actionCard() {
    return this.args.actionCard || false;
  }
  get gridContainer() {
    return this.args.gridContainer || false;
  }
  get total() {
    return this.args.total || 0;
  }
}
