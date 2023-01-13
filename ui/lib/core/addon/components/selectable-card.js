import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
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
 *
 * If hasSearchInput is enabled:
 * @param {boolean} [hasSearchInput] - boolean to show form inputs.
 * @param {array} [models] - list of model names.
 * @param {string} [backend] - backend name.
 * @param {string} [pagePath] - page path for transitionTo in form button.
 * @param {string} [placeholder] - placeholder for input.
 * @param {string} [buttonText] - button text for form button.
 * @param {boolean} [disallowNewItems] - boolean to not enable search dropdown with new items.
 */

export default class SelectableCard extends Component {
  @service router;
  @tracked value = '';

  get total() {
    return this.args.total || 0;
  }

  @action
  transitionToPage() {
    this.router.transitionTo(this.args.pagePath, this.value);
  }
  @action
  handleInput(value) {
    // if it comes in from the fallback component then the value is a string otherwise it's an array
    if (Array.isArray(value)) {
      this.value = value[0];
    } else {
      this.value = value;
    }
  }
}
