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
import { task } from 'ember-concurrency';
export default class GetCredentialsCard extends Component {
  @service router;
  @service store;
  role = '';
  dynamicRole = '';
  staticRole = '';
  @tracked buttonDisabled = true;
  @tracked roleDoesNotExist = false;

  @action
  async getSelectedValue(selectValue) {
    this.role = selectValue[0];
    let role = this.role;
    if (role) {
      let dynamicRole = this.store.peekRecord('database/role', role);
      let staticRole = this.store.peekRecord('database/static-role', role);
      if (!dynamicRole && !staticRole) {
        // situation when they type in the role and they don't have list permissions
        // make network request to creds
        // if nothing is returned, then return error message
        let fetchCredentialsResult = await this.fetchCredentials.perform(role);
        console.log(fetchCredentialsResult, 'FETCH CREDS RESULT');
        if (fetchCredentialsResult.length === 2) {
          // returned an array of two errors, which means both requests failed and no role exist.
          this.roleDoesNotExist = !this.roleDoesNotExist;
        }
      }
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
  @task(function*(role) {
    let backend = this.args.backend;
    let errors = [];
    try {
      let staticRoleAttempt = yield this.store.queryRecord('database/credential', {
        backend,
        secret: role,
        roleType: 'static',
      });
      // if successful will return result
      return staticRoleAttempt;
    } catch (error) {
      errors.push(error.errors);
    }
    try {
      let dynamicRoleAttempt = yield this.store.queryRecord('database/credential', {
        backend,
        secret: role,
        roleType: 'dynamic',
      });
      return dynamicRoleAttempt;
    } catch (error) {
      errors.push(error.errors);
    }
    return errors;
  })
  fetchCredentials;
}
