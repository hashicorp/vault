/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module SecretListHeaderTab
 * SecretListHeaderTab component passes in properties that are used to check capabilities and either display or not display the component.
 * Use case was first for the Database Secret Engine, but should be used in future iterations as we don't generally want to show things the user does not
 * have access to
 *
 *
 * @example
 * <SecretListHeaderTab @displayName='Database' @id='database-2' @path='roles' @label='Roles' @tab='roles'/>
 *
 * @param {string} [displayName] - set on options-for-backend this sets a conditional to see if capabilities are being checked
 * @param {string} [id] - if fetching capabilities used for making the query url.  It is the name the user has assigned to the instance of the engine.
 * @param {string} [path] - set on options-for-backend this tells us the specifics of the URL the query should hit.
 * @param {string} label - The name displayed on the tab.   Set on the options-for-backend.
 * @param {string} [tab] - The name of the tab.  Set on the options-for-backend.
 * @param {string} [link] - If within an engine provide the name of the link that is defined in the routes file fo the engine, example : 'overview'.
 *
 */
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

export default class SecretListHeaderTab extends Component {
  @service capabilities;

  @tracked dontShowTab;

  constructor() {
    super(...arguments);
    this.fetchCapabilities();
  }

  pathQuery(backend, path) {
    return {
      id: `${backend}/${path}/`,
    };
  }

  async fetchCapabilities() {
    // For now only check capabilities for the Database Secrets Engine
    if (this.args.displayName === 'Database') {
      const { path, id: backend } = this.args;
      const pathKey = path === 'config' ? 'databaseConfig' : 'databaseRoles';
      const { canList, canCreate, canUpdate } = await this.capabilities.for(pathKey, { backend });
      this.dontShowTab = !canList && !canCreate && !canUpdate;
    }
  }
}
