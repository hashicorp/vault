/**
 * @module SecretListHeaderTab
 * SecretListHeaderTab component passes in properties that are used to check capabilities and either display or not display the component.
 * Use case was first for the Database Secret Engine, but should be used in future iterations as we don't generally want to show things the user does not
 * have access to
 *
 *
 * @example
 * ```js
 * <SecretListHeaderTab @displayName='Database' @id='database-2' @path='roles' @label='Roles' @tab='roles'/>
 * ```
 * @param {string} [displayName] - set on options-for-backend this sets a conditional to see if capabilities are being checked
 * @param {string} [id] - if fetching capabilities used for making the query url.  It is the name the user has assigned to the instance of the engine.
 * @param {string} [path] - set on options-for-backend this tells us the specifics of the URL the query should hit.
 * @param {string} label - The name displayed on the tab.   Set on the options-for-backend.
 * @param {string} [tab] - The name of the tab.  Set on the options-for-backend.
 *
 */
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

export default class SecretListHeaderTab extends Component {
  @service store;
  @tracked dontShowTab;
  constructor() {
    super(...arguments);
    this.fetchCapabilities.perform();
  }

  pathQuery(backend, path) {
    return {
      id: `${backend}/${path}/`,
    };
  }

  @task(function*() {
    // For now only check capabilities for the Database Secrets Engine
    // ARG TODO try and think of a way to clean up?  Unsure how.
    if (this.args.displayName === 'Database') {
      let peekRecordRoles = yield this.store.peekRecord('capabilities', 'database/roles/');
      let peekRecordStaticRoles = yield this.store.peekRecord('capabilities', 'database/static-roles/');
      let peekRecordConnections = yield this.store.peekRecord('capabilities', 'database/config/');
      // peekRecord if the capabilities store data is there for the connections (config) and roles model
      if (
        (peekRecordRoles && this.args.path === 'roles') ||
        (peekRecordStaticRoles && this.args.path === 'roles')
      ) {
        let roles = !peekRecordRoles.canList && !peekRecordRoles.canCreate && !peekRecordRoles.canUpdate;
        let staticRoles =
          !peekRecordStaticRoles.canList &&
          !peekRecordStaticRoles.canCreate &&
          !peekRecordStaticRoles.canUpdate;
        this.dontShowTab = roles && staticRoles;
        return;
      }
      if (peekRecordConnections && this.args.path === 'config') {
        this.dontShowTab =
          !peekRecordConnections.canList &&
          !peekRecordConnections.canCreate &&
          !peekRecordConnections.canUpdate;
        return;
      }
      // otherwise queryRecord and create an instance on the capabilities.
      let response = yield this.store.queryRecord(
        'capabilities',
        this.pathQuery(this.args.id, this.args.path)
      );
      this.dontShowTab = !response.canList && !response.canCreate && !response.canUpdate;
    }
  })
  fetchCapabilities;
}
