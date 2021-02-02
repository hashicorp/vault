/**
 * @module GenerateCredentialsDatabase
 * GenerateCredentialsDatabase component is used on the credentials route for the Database metrics.
 * The component assumes that you will need to make an ajax request using queryRecord to return a model for the component that has username, password, leaseId and leaseDuration
 *
 * @example
 * ```js
 * <GenerateCredentialsDatabase @backendPath="database" @backendType="database" @roleName="my-role"/>
 * ```
 * @param {string} backendPath - the secret backend name.  This is used in the breadcrumb.
 * @param {object} backendType - the secret type.  Expected to be database.
 * @param {string} roleName - the id of the credential returning.
 */

import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class GenerateCredentialsDatabase extends Component {
  @service store;
  // set on the component
  backendType = null;
  backendPath = null;
  roleName = null;
  @tracked roleType = '';
  @tracked model = null;

  constructor() {
    super(...arguments);
    this.fetchCredentials.perform();
  }

  @task(function*() {
    let { roleName, backendPath } = this.args;
    let errors = [];
    console.log(this.args.backendPath, 'ARGS');
    try {
      let newModel = yield this.store.queryRecord('database/credential', {
        backend: backendPath,
        secret: roleName,
        roleType: 'static',
      });
      // if successful will return result
      this.model = newModel;
      this.roleType = 'static';
      return;
    } catch (error) {
      errors.push(error.errors);
    }
    try {
      let newModel = yield this.store.queryRecord('database/credential', {
        backend: backendPath,
        secret: roleName,
        roleType: 'dynamic',
      });
      this.model = newModel;
      this.roleType = 'dynamic';
      return;
    } catch (error) {
      errors.push(error.errors);
    }
    this.roleType = 'noRoleFound';
  })
  fetchCredentials;

  @action redirectPreviousPage() {
    window.history.back();
  }
}
