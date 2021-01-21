import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
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
export default class GetCredentialsCard extends Component {
  @service router;
  @service store;
  role = '';
  dynamicRole = '';
  staticRole = '';
  @tracked buttonDisabled = true;

  @action
  getSelectedValue(selectValue) {
    this.role = selectValue[0];
    if (this.role) {
      this.roleType = this.store.peekRecord('database/role', this.role) ? 'dynamic' : 'static';
    }
    // Ember Octane way of getting away from toggleProperty
    this.buttonDisabled = !this.buttonDisabled;
  }
  @action
  transitionToCredential() {
    if (!this.role) {
      return;
    }
    let role = this.role;
    let roleType = this.roleType;
    this.router.transitionTo('vault.cluster.secrets.backend.credentials', role, {
      queryParams: { roleType: roleType },
    });
  }
}
