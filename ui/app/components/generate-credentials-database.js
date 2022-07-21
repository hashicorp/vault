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
 * @param {string} roleType - either 'static', 'dynamic', or falsey.
 * @param {string} roleName - the id of the credential returning.
 * @param {object} model - database/credential model passed in. If no data, should have errorTitle, errorMessage, and errorHttpStatus
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class GenerateCredentialsDatabase extends Component {
  get errorTitle() {
    return this.args.model.errorTitle || 'Something went wrong';
  }

  @action redirectPreviousPage() {
    window.history.back();
  }
}
