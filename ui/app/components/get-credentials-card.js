/**
 * @module GetCredentialsCard
 * GetCredentialsCard components are card-like components that display a title, and SearchSelect component that sends you to a credentials route for the selected item.
 * They are designed to be used in containers that act as flexbox or css grid containers.
 *
 * @example
 * ```js
 * <GetCredentialsCard @title="Get Credentials" @searchLabel="Role to use" @models={{array 'database/roles'}} />
 * ```
 * @param title=null {String} - The title displays the card title
 * @param searchLabel=null {String} - The text above the searchSelect component
 * @param models=null {Array} - An array of model types to fetch from the API.  Passed through to SearchSelect component
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
export default class GetCredentialsCard extends Component {
  @service router;
  @service store;
  role = '';
  manuallyEnteredRole = '';
  @tracked buttonDisabled = true;

  @action
  async getSelectedValue(selectValue) {
    this.role = selectValue[0];
    let role = this.role;
    if (role) {
      this.buttonDisabled = false;
    }
  }
  @action
  async transitionToCredential() {
    if (this.role || this.manuallyEnteredRole) {
      let role = this.role || this.manuallyEnteredRole;
      this.router.transitionTo('vault.cluster.secrets.backend.credentials', role);
    }
  }
  @action
  getManuallyEnteredValue(path, value) {
    this.manuallyEnteredRole = value;
    if (value) {
      this.buttonDisabled = false;
    }
  }
}
