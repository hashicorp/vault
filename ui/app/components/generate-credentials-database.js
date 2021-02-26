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
  @tracked errorMessage = '';
  @tracked errorHttpStatus = '';
  @tracked errorTitle = 'Something went wrong';

  constructor() {
    super(...arguments);
    this.fetchCredentials.perform();
  }

  @task(function*() {
    let { roleName, backendPath } = this.args;
    try {
      let newModel = yield this.store.queryRecord('database/credential', {
        backend: backendPath,
        secret: roleName,
        roleType: 'static',
      });
      this.model = newModel;
      this.roleType = 'static';
      return;
    } catch (error) {
      this.errorHttpStatus = error.httpStatus; // set default http
      this.errorMessage = `We ran into a problem and could not continue: ${error.errors[0]}`;
      if (error.httpStatus === 403) {
        // 403 is forbidden
        this.errorTitle = 'You are not authorized';
        this.errorMessage =
          "Role wasn't found or you do not have permissions. Ask your administrator if you think you should have access.";
      }
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
      if (error.httpStatus === 403) {
        // 403 is forbidden
        this.errorHttpStatus = error.httpStatus; // override default httpStatus which could be 400 which always happens on either dynamic or static depending on which kind of role you're querying
        this.errorTitle = 'You are not authorized';
        this.errorMessage =
          "Role wasn't found or you do not have permissions. Ask your administrator if you think you should have access.";
      }
      if (error.httpStatus == 500) {
        // internal server error happens when empty creation statement on dynamic role creation only
        this.errorHttpStatus = error.httpStatus;
        this.errorTitle = 'Internal Error';
        this.errorMessage = error.errors[0];
      }
    }
    this.roleType = 'noRoleFound';
  })
  fetchCredentials;

  @action redirectPreviousPage() {
    window.history.back();
  }
}
