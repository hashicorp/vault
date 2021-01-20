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
  roleType = 'dynamic';
  @tracked model = null;

  constructor() {
    super(...arguments);
    this.fetchCredentials.perform();
  }
  @task(function*() {
    let { roleType, roleName, backendType } = this.args;
    let newModel = yield this.store.queryRecord('database/credential', {
      backend: backendType,
      secret: roleName,
      roleType,
    });
    this.model = newModel;
  })
  fetchCredentials;

  @action redirectPreviousPage() {
    window.history.back();
  }
}
