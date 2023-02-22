/**
 * @module GetCredentialsCard
 * GetCredentialsCard components are card-like components that display a title, and SearchSelect component that sends you to a credentials route for the selected item.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <GetCredentialsCard @title="Get Credentials" @searchLabel="Role to use" @models={{array 'database/roles'}} @type="role" @backend={{model.backend}}/>
 * ```
 * @param {string} title - The title displays the card title
 * @param {string} searchLabel - The text above the searchSelect component
 * @param {array} models - An array of model types to fetch from the API. Passed through to SearchSelect component
 * @param {string} type - 'role' or 'secret' - determines where the transitionTo goes
 * @param {boolean} [renderInputSearch=false] - If true, renders InputSearch instead of SearchSelect
 * @param {string} [subText] - Text below title
 * @param {string} [placeholder] - Input placeholder text (default for SearchSelect is 'Search', none for InputSearch)
 * @param {string} backend - Passed to SearchSelect query method to fetch dropdown options
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
export default class GetCredentialsCard extends Component {
  @service router;
  @service store;
  @tracked role = '';
  @tracked secret = '';

  get buttonDisabled() {
    return !this.role && !this.secret;
  }

  @action
  transitionToCredential() {
    const role = this.role;
    const secret = this.secret;
    if (role) {
      this.router.transitionTo('vault.cluster.secrets.backend.credentials', role);
    }
    if (secret) {
      this.router.transitionTo('vault.cluster.secrets.backend.show', secret);
    }
  }

  @action
  handleInput(value) {
    if (this.args.type === 'role') {
      // if it comes in from the fallback component then the value is a string otherwise it's an array
      // which currently only happens if type is role.
      if (Array.isArray(value)) {
        this.role = value[0];
      } else {
        this.role = value;
      }
    }
    if (this.args.type === 'secret') {
      this.secret = value;
    }
  }
}
