import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

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
 * @param {number} [total = 0] - the number displayed as the largest text in the component.
 * @param {string} [subText] - subText describes the total.
 * @param {string} [actionText] - action text link.
 * @param {string} [actionTo] - route where link will take you.
 * @param {string} [queryParam] - tab for the route the link will take you.
 * @param {string} [type] - type used in the link type.
 * @param {boolean} [hasSearchInput] - boolean to show form inputs
 * @param {string} [models] - name of model
 * @param {string} [backend] - name of backend
 * @param {string} [placeholder] - name of backend
 */

export default class SelectableCard extends Component {
  @service router;

  get total() {
    return this.args.total || 0;
  }
}
