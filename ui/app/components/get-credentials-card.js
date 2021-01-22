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
  manuallyEnteredRole = '';
  dynamicRole = '';
  staticRole = '';
  @tracked buttonDisabled = true;
  @tracked roleDoesNotExist = false;

  @action
  async getSelectedValue(selectValue) {
    this.role = selectValue[0];
    let role = this.role;
    if (role) {
      this.buttonDisabled = false;
      let dynamicRole = this.store.peekRecord('database/role', role);
      let staticRole = this.store.peekRecord('database/static-role', role);
      if (!dynamicRole && !staticRole) {
        // situation when they type in the role and they don't have list permissions
        // if nothing is returned, then return error message
        let fetchCredentialsResult = await this.fetchCredentials.perform(role);
        console.log(fetchCredentialsResult, 'FETCH CREDS RESULT');
        if (fetchCredentialsResult.length === 2) {
          // returned an array of two errors, which means both requests failed and no role exist.
          this.roleDoesNotExist = !this.roleDoesNotExist;
        }
      }
    }
  }
  @action
  async transitionToCredential() {
    if (this.role || this.manuallyEnteredRole) {
      let role = this.role || this.manuallyEnteredRole;
      let roleType = this.roleType;
      let results = await this.fetchCredentials.perform(role);
      // ARG TODO handle better.
      if (results.length === 2) {
        // error, return noRoleFound as type and let generate-credentials-database handle
        roleType = 'noRoleFound';
      } else {
        roleType = results;
      }

      this.router.transitionTo('vault.cluster.secrets.backend.credentials', role, {
        queryParams: { roleType: roleType },
      });
    }
  }
  @action
  getManuallyEnteredValue(path, value) {
    this.manuallyEnteredRole = value;
    if (value) {
      this.buttonDisabled = false;
    }
  }
  @task(function*(role) {
    let backend = this.args.backend;
    let errors = [];
    try {
      yield this.store.queryRecord('database/credential', {
        backend,
        secret: role,
        roleType: 'static',
      });
      // if successful will return result
      return 'static';
    } catch (error) {
      errors.push(error.errors);
    }
    try {
      yield this.store.queryRecord('database/credential', {
        backend,
        secret: role,
        roleType: 'dynamic',
      });
      return 'dynamic';
    } catch (error) {
      errors.push(error.errors);
    }
    return errors;
  })
  fetchCredentials;
}
