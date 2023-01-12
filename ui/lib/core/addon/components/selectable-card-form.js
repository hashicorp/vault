/**
 * @module SelectableCardForm
 * SelectableCardForm components are card-like components that display a title, and SearchSelect component that sends you to a route for the selected item.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <SelectableCardForm @title="Get Credentials" @searchLabel="Role to use" @models={{array 'database/roles'}} @type="role" @backend={{model.backend}}/>
 * ```
 * @param {string} title - The title displays the card title
 * @param {string} searchLabel - The text above the searchSelect component
 * @param {array} models - An array of model types to fetch from the API. Passed through to SearchSelect component
 * @param {string} [subText] - Text below title
 * @param {string} [placeholder] - Input placeholder text (default for SearchSelect is 'Search', none for InputSearch)
 * @param {string} backend - Passed to SearchSelect query method to fetch dropdown options
 * @param {string} testLabel - The name for the data-test attribute
 * @param {string} pagePath - The name of the path where the page will transition to
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
export default class SelectableCardForm extends Component {
  @service router;
  @tracked value = '';

  get buttonDisabled() {
    return !this.value;
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
